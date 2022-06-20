package gaccauth

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfig(b *googleAccountAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: pathConfigPattern,
		Callbacks: ActionCallback{
			logical.UpdateOperation: b.pathConfigWrite,
			logical.ReadOperation:   b.pathConfigRead,
		},
		Fields: Schema{
			clientIDConfigProp: {
				Type:        framework.TypeString,
				Description: "Google OAuth client id",
			},
			clientSecretConfigProp: {
				Type:        framework.TypeString,
				Description: "Google OAuth client secret",
			},
			clientOAuthRedirectURLProp: {
				Type:        framework.TypeString,
				Description: "Google OAuth redirect URL",
			},
			clientFetchGroupsConfigProp: {
				Type:        framework.TypeBool,
				Description: "Google fetch groups",
			},
			clientServiceAccountKeyConfigProp: {
				Type:        framework.TypeString,
				Description: "Google service account key content",
			},
			clientDelegationUserConfigProp: {
				Type:        framework.TypeString,
				Description: "Google delegation email address",
			},
		},
	}
}

func (b *googleAccountAuthBackend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON(pathConfigEntry, googleAuth{
		ClientID:       data.Get(clientIDConfigProp).(string),
		ClientSecret:   data.Get(clientSecretConfigProp).(string),
		RedirectURL:    data.Get(clientOAuthRedirectURLProp).(string),
		FetchGroups:    data.Get(clientFetchGroupsConfigProp).(bool),
		ServiceAccount: data.Get(clientServiceAccountKeyConfigProp).(string),
		DelegationUser: data.Get(clientDelegationUserConfigProp).(string),
	})

	if err != nil {
		return nil, err
	}

	return nil, req.Storage.Put(ctx, entry)
}

func (b *googleAccountAuthBackend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)

	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, nil
	}

	configMap := GenericMap{
		clientIDConfigProp:                config.ClientID,
		clientSecretConfigProp:            config.ClientSecret,
		clientOAuthRedirectURLProp:        config.RedirectURL,
		clientFetchGroupsConfigProp:       config.FetchGroups,
		clientServiceAccountKeyConfigProp: config.ServiceAccount,
		clientDelegationUserConfigProp:    config.DelegationUser,
	}

	return &logical.Response{Data: configMap}, nil
}

// Config returns the configuration for this backend.
func (b *googleAccountAuthBackend) config(ctx context.Context, s logical.Storage) (*googleAuth, error) {
	entry, err := s.Get(ctx, pathConfigEntry)

	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	var result googleAuth
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, fmt.Errorf("error reading configuration: %s", err)
	}

	return &result, nil
}

const (
	clientDelegationUserConfigProp    = "delegation_user"
	clientFetchGroupsConfigProp       = "fetch_groups"
	clientIDConfigProp                = "client_id"
	clientOAuthRedirectURLProp        = "redirect_url"
	clientSecretConfigProp            = "client_secret"
	clientServiceAccountKeyConfigProp = "service_acc_key"
	pathConfigEntry                   = "config"
	pathConfigPattern                 = "config"
)
