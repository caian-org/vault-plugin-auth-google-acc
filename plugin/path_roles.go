package gaccauth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	pathRolesNameProp        = "name"
	pathRolesPoliciesProp    = "policies"
	pathRolesBoundEmailsProp = "bound_emails"
	pathRolesBoundGroupsProp = "bound_groups"
	pathRolesMaxTTLProp      = "max_ttl"
	pathRolesTTLProp         = "ttl"
	errEmptyRoleName         = "role name is required"
)

const pathRolesHelpSyn = `
A role is required to login under the Google auth backend. A role binds Vault policies and has required attributes that
an authenticating entity must fulfill to login against this role. After authenticating the instance, Vault uses the
bound policies to determine which resources the authorization token for the instance can access.
`

type googleAuthRole struct {
	Policies    []string      `json:"policies" structs:"policies" mapstructure:"policies"`
	BoundGroups []string      `json:"bound_groups" structs:"bound_groups" mapstructure:"bound_groups"`
	BoundEmails []string      `json:"bound_emails" structs:"bound_emails" mapstructure:"bound_emails"`
	TTL         time.Duration `json:"ttl" structs:"ttl" mapstructure:"ttl"`
	MaxTTL      time.Duration `json:"max_ttl" structs:"max_ttl" mapstructure:"max_ttl"`
}

func pathRoles(b *googleAccountAuthBackend) []*framework.Path {
	role := &framework.Path{
		Pattern:         fmt.Sprintf("role/%s", framework.GenericNameRegex("name")),
		HelpSynopsis:    "Create a Google role with associated policies and required attributes.",
		HelpDescription: pathRolesHelpSyn,
		ExistenceCheck:  b.pathRoleExistenceCheck,
		Fields: Schema{
			pathRolesNameProp: {
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},
			pathRolesPoliciesProp: {
				Type:        framework.TypeCommaStringSlice,
				Description: "Policies to be set on tokens issued using this role.",
			},
			pathRolesBoundGroupsProp: {
				Type:        framework.TypeCommaStringSlice,
				Description: "Comma separate list of groups, at least one of which the user must be in to grant this role.",
			},
			pathRolesBoundEmailsProp: {
				Type:        framework.TypeCommaStringSlice,
				Description: "Comma separate list of usernames, which the user must be in to grant this role.",
			},
			pathRolesTTLProp: {
				Type:        framework.TypeDurationSecond,
				Description: "Duration in seconds after which the issued token should expire.",
			},
			pathRolesMaxTTLProp: {
				Type:        framework.TypeDurationSecond,
				Description: "The maximum allowed lifetime of tokens issued using this role.",
			},
		},
		Callbacks: ActionCallback{
			logical.CreateOperation: b.pathRoleUpsert,
			logical.ReadOperation:   b.pathRoleRead,
			logical.UpdateOperation: b.pathRoleUpsert,
			logical.DeleteOperation: b.pathRoleDelete,
		},
	}

	roles := &framework.Path{
		Pattern:         "role/?",
		HelpSynopsis:    "Lists all the roles that are registered with Vault.",
		HelpDescription: "Lists all roles under the Google backend by name.",
		Callbacks:       ActionCallback{logical.ListOperation: b.pathRoleList},
	}

	return []*framework.Path{role, roles}
}

func (b *googleAccountAuthBackend) pathRoleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.getDecodedRole(ctx, req.Storage, data.Get(pathRolesNameProp).(string))
	if err != nil {
		return false, err
	}

	return entry != nil, nil
}

func (b *googleAccountAuthBackend) pathRoleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get(pathRolesNameProp).(string)
	if name == "" {
		return logical.ErrorResponse(errEmptyRoleName), nil
	}

	if err := req.Storage.Delete(ctx, fmt.Sprintf("role/%s", name)); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *googleAccountAuthBackend) pathRoleUpsert(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := strings.ToLower(data.Get(pathRolesNameProp).(string))
	if name == "" {
		return logical.ErrorResponse(errEmptyRoleName), nil
	}

	r, err := b.getDecodedRole(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}

	// new role
	if r == nil {
		r = &googleAuthRole{}
	}

	if err := r.parseAndValidateInput(b.System(), req.Operation, data); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	entry, err := logical.StorageEntryJSON(fmt.Sprintf("role/%s", name), r)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *googleAccountAuthBackend) pathRoleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roles, err := req.Storage.List(ctx, "role/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(roles), nil
}

func (b *googleAccountAuthBackend) pathRoleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get(pathRolesNameProp).(string)
	if name == "" {
		return logical.ErrorResponse(errEmptyRoleName), nil
	}

	role, err := b.getDecodedRole(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, nil
	}

	response := &logical.Response{
		Data: GenericMap{
			pathRolesNameProp:        name,
			pathRolesPoliciesProp:    role.Policies,
			pathRolesBoundGroupsProp: role.BoundGroups,
			pathRolesBoundEmailsProp: role.BoundEmails,
			pathRolesTTLProp:         fmt.Sprint(role.TTL / time.Second),
			pathRolesMaxTTLProp:      fmt.Sprint(role.MaxTTL / time.Second),
		},
	}

	return response, nil
}

///////////////////////////////////////////////////////////////////////////////

func (b *googleAccountAuthBackend) getDecodedRole(ctx context.Context, s logical.Storage, name string) (*googleAuthRole, error) {
	entry, err := s.Get(ctx, fmt.Sprintf("role/%s", name))
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	role := &googleAuthRole{}
	if err := entry.DecodeJSON(role); err != nil {
		return nil, err
	}

	return role, nil
}

///////////////////////////////////////////////////////////////////////////////

func (r *googleAuthRole) parseAndValidateInput(sys logical.SystemView, op logical.Operation, data *framework.FieldData) error {
	boundEmails := getFilteredStringSliceData(data, pathRolesBoundEmailsProp)
	if boundEmails == nil {
		r.BoundEmails = []string{}
	} else {
		r.BoundEmails = *boundEmails
	}

	boundGroups := getFilteredStringSliceData(data, pathRolesBoundGroupsProp)
	if boundGroups == nil {
		r.BoundGroups = []string{}
	} else {
		r.BoundGroups = *boundGroups
	}

	if len(r.BoundEmails)+len(r.BoundGroups) == 0 {
		return fmt.Errorf("at least one email address or group must be set")
	}

	invalidEmailAddrs := []string{}
	for _, emailAddr := range append(r.BoundEmails, r.BoundGroups...) {
		if !isValidEmail(emailAddr) {
			invalidEmailAddrs = append(invalidEmailAddrs, emailAddr)
		}
	}

	if len(invalidEmailAddrs) > 0 {
		return fmt.Errorf("one or more provided email addresses are invalid: %s", strings.Join(invalidEmailAddrs, ", "))
	}

	//////////////////////

	if policies, ok := data.GetOk(pathRolesPoliciesProp); ok {
		r.Policies = policyutil.ParsePolicies(policies)
	} else {
		return fmt.Errorf("unable to retrieve policies")
	}

	if len(r.Policies) == 0 {
		return fmt.Errorf("at least one policy must be defined")
	}

	if strings.ToLower(r.Policies[0]) == "root" {
		return fmt.Errorf("cannot use root policy")
	}

	//////////////////////

	if ttl, err := getPositiveIntData(data, pathRolesTTLProp); err == nil {
		if ttl == nil {
			// fallbacks to 1 hour when unset
			r.TTL = time.Duration(1) * time.Hour
		} else {
			r.TTL = time.Duration(*ttl) * time.Second
		}
	} else {
		return err
	}

	if mttl, err := getPositiveIntData(data, pathRolesMaxTTLProp); err == nil {
		if mttl == nil {
			// fallbacks to 1 day when unset
			r.MaxTTL = time.Duration(24) * time.Hour
		} else {
			r.MaxTTL = time.Duration(*mttl) * time.Second
		}
	} else {
		return err
	}

	if r.TTL.Hours() > r.MaxTTL.Hours() {
		return fmt.Errorf("ttl (%s) cannot be greater than max_ttl (%s)", r.TTL.String(), r.MaxTTL.String())
	}

	return nil
}
