package gaccauth

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	directory "google.golang.org/api/admin/directory/v1"
	goauth "google.golang.org/api/oauth2/v2"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathLogin(b *googleAccountAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: pathLoginPattern,
		Fields: Schema{
			pathLoginGoogleAuthCodeParam: {
				Type:        framework.TypeString,
				Description: "Google authentication code. Required.",
			},
			pathLoginRoleNameParam: {
				Type:        framework.TypeString,
				Description: "Name of the role against which the login is being attempted. Required.",
			},
		},
		Callbacks: ActionCallback{
			logical.UpdateOperation:         b.pathLoginAuthFlow,
			logical.AliasLookaheadOperation: b.pathLoginAuthFlow,
		},
	}
}

func (b *googleAccountAuthBackend) pathLoginAuthFlow(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	code := data.Get(pathLoginGoogleAuthCodeParam).(string)

	roleName := data.Get(pathLoginRoleNameParam).(string)
	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("role '%s' not found", roleName)), nil
	}

	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return logical.ErrorResponse("missing config"), nil
	}

	googleConfig := config.oauth2Config()
	token, err := googleConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	user, groups, err := b.authenticate(config, token)
	if err != nil {
		return nil, err
	}

	policies, err := b.authorize(req.Storage, role, user, groups)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	encodedToken, err := encodeToken(token)
	if err != nil {
		return nil, err
	}

	auth := &logical.Auth{
		DisplayName: user.Email,
		Policies:    policies,
		InternalData: GenericMap{
			"token": encodedToken,
			"role":  roleName,
		},
		Metadata: map[string]string{
			"username": user.Email,
			"domain":   user.Hd,
		},
		LeaseOptions: logical.LeaseOptions{
			TTL:       role.TTL,
			Renewable: true,
		},
	}

	return &logical.Response{Auth: auth}, nil
}

func (b *googleAccountAuthBackend) authRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	encodedToken, ok := req.Auth.InternalData["token"].(string)
	if !ok {
		return nil, errors.New("no refresh token from previous login")
	}

	roleName, ok := req.Auth.InternalData["role"].(string)
	if !ok {
		return nil, errors.New("no role name from previous login")
	}

	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("role '%s' not found", roleName)
	}

	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return logical.ErrorResponse("missing config"), nil
	}

	token, err := decodeToken(encodedToken)
	if err != nil {
		return nil, err
	}

	user, groups, err := b.authenticate(config, token)
	if err != nil {
		return nil, err
	}

	policies, err := b.authorize(req.Storage, role, user, groups)
	if err != nil {
		return nil, err
	}

	if !strSliceEquals(policies, req.Auth.Policies) {
		return logical.ErrorResponse(fmt.Sprintf("policies do not match. new policies: %s. old policies: %s.", policies, req.Auth.Policies)), nil
	}

	return framework.LeaseExtend(role.TTL, role.MaxTTL, b.System())(ctx, req, d)
}

func (b *googleAccountAuthBackend) authenticate(googleAuth *googleAuth, token *oauth2.Token) (*goauth.Userinfo, []string, error) {
	client := googleAuth.oauth2Config().Client(context.Background(), token)

	userService, err := goauth.New(client)
	if err != nil {
		return nil, nil, err
	}

	user, err := goauth.NewUserinfoV2MeService(userService).Get().Do()
	if err != nil {
		return nil, nil, err
	}

	groups := []string{}
	if googleAuth.FetchGroups {
		scope := "https://www.googleapis.com/auth/admin.directory.group.readonly"

		serviceAccountCredential, err := google.JWTConfigFromJSON([]byte(googleAuth.ServiceAccount), scope)
		if err != nil {
			return nil, nil, err
		}

		serviceAccountCredential.Subject = googleAuth.DelegationUser
		saClient, err := directory.New(serviceAccountCredential.Client(context.Background()))
		if err != nil {
			return nil, nil, err
		}

		response, err := saClient.Groups.List().UserKey(user.Email).Do()
		if err != nil {
			return nil, nil, err
		}

		for _, g := range response.Groups {
			groups = append(groups, g.Email)
		}
	}

	return user, groups, nil
}

func (b *googleAccountAuthBackend) authorize(storage logical.Storage, role *role, user *goauth.Userinfo, groups []string) ([]string, error) {
	if user.Hd != role.BoundDomain && role.BoundDomain != "" {
		return nil, fmt.Errorf("user %s is not part of required domain %s, found %s", user.Email, role.BoundDomain, user.Hd)
	}

	// Is this user in one of the bound groups for this role?
	isGroupMember := strSliceHasIntersection(groups, role.BoundGroups)
	isUserMember := strSliceHasIntersection([]string{user.Email}, role.BoundEmails)

	if !isGroupMember && !isUserMember {
		return nil, fmt.Errorf("user is not allowed to use this role")
	}

	return role.Policies, nil
}

const (
	pathLoginPattern             = "login"
	pathLoginGoogleAuthCodeParam = "code"
	pathLoginRoleNameParam       = "role"
)
