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

const (
	pathLoginPattern            = "login"
	pathLoginGoogleAuthCodeProp = "code"
	pathLoginRoleNameProp       = "role"
)

func pathLogin(b *googleAccountAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: pathLoginPattern,
		Fields: Schema{
			pathLoginGoogleAuthCodeProp: {
				Type:        framework.TypeString,
				Description: "Google authentication code",
			},
			pathLoginRoleNameProp: {
				Type:        framework.TypeString,
				Description: "Name of the role against which the login is being attempted",
			},
		},
		Callbacks: ActionCallback{
			logical.UpdateOperation:         b.pathLoginAuthFlow,
			logical.AliasLookaheadOperation: b.pathLoginAuthFlow,
		},
	}
}

func (b *googleAccountAuthBackend) pathLoginAuthFlow(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	code := data.Get(pathLoginGoogleAuthCodeProp).(string)
	roleName := data.Get(pathLoginRoleNameProp).(string)
	role, err := b.getDecodedRole(ctx, req.Storage, roleName)

	if err != nil {
		return nil, err
	}

	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("role '%s' not found", roleName)), nil
	}

	googleOAuth, err := b.getGoogleOAuthConfig(ctx, req.Storage)

	if err != nil {
		return nil, err
	}

	if googleOAuth == nil {
		return logical.ErrorResponse("missing config"), nil
	}

	googleConfig := googleOAuth.build()
	token, err := googleConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	user, groups, err := b.authenticate(googleOAuth, token)
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

	response := &logical.Response{
		Auth: &logical.Auth{
			DisplayName: user.Email,
			Policies:    policies,
			InternalData: GenericMap{
				"token": encodedToken,
				"role":  roleName,
			},
			Metadata: map[string]string{
				"username": user.Email,
			},
			LeaseOptions: logical.LeaseOptions{
				TTL:       role.TTL,
				Renewable: true,
			},
		},
	}

	return response, nil
}

///////////////////////////////////////////////////////////////////////////////

func (b *googleAccountAuthBackend) authRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	encodedToken, ok := req.Auth.InternalData["token"].(string)
	if !ok {
		return nil, errors.New("no refresh token from previous login")
	}

	roleName, ok := req.Auth.InternalData["role"].(string)
	if !ok {
		return nil, errors.New("no role name from previous login")
	}

	role, err := b.getDecodedRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("role '%s' not found", roleName)
	}

	googleOAuth, err := b.getGoogleOAuthConfig(ctx, req.Storage)

	if err != nil {
		return nil, err
	}

	if googleOAuth == nil {
		return logical.ErrorResponse("missing Google OAuth config"), nil
	}

	token, err := decodeToken(encodedToken)
	if err != nil {
		return nil, err
	}

	user, groups, err := b.authenticate(googleOAuth, token)
	if err != nil {
		return nil, err
	}

	policies, err := b.authorize(req.Storage, role, user, groups)
	if err != nil {
		return nil, err
	}

	if !sliceEquals(policies, req.Auth.Policies) {
		return logical.ErrorResponse(fmt.Sprintf("policies do not match. new policies: %s. old policies: %s.", policies, req.Auth.Policies)), nil
	}

	return framework.LeaseExtend(role.TTL, role.MaxTTL, b.System())(ctx, req, d)
}

func (b *googleAccountAuthBackend) authenticate(googleOAuth *googleOAuth, token *oauth2.Token) (*goauth.Userinfo, []string, error) {
	client := googleOAuth.build().Client(context.Background(), token)

	userService, err := goauth.New(client)
	if err != nil {
		return nil, nil, err
	}

	user, err := goauth.NewUserinfoV2MeService(userService).Get().Do()
	if err != nil {
		return nil, nil, err
	}

	groups := []string{}

	if !googleOAuth.FetchGroups {
		return user, groups, nil
	}

	saScopes := "https://www.googleapis.com/auth/admin.directory.group.readonly"
	saCredential, err := google.JWTConfigFromJSON([]byte(googleOAuth.ServiceAccount), saScopes)
	if err != nil {
		return nil, nil, err
	}

	saCredential.Subject = googleOAuth.DelegationUser
	saClient, err := directory.New(saCredential.Client(context.Background()))
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

	return user, groups, nil
}

func (b *googleAccountAuthBackend) authorize(storage logical.Storage, role *googleAuthRole, user *goauth.Userinfo, groups []string) ([]string, error) {
	isGroupMember := sliceContains(groups, role.BoundGroups)
	isUserMember := sliceContains([]string{user.Email}, role.BoundEmails)

	if isUserMember || isGroupMember {
		return role.Policies, nil
	}

	return nil, fmt.Errorf("user is not allowed to use this role")
}
