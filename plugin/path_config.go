package gaccauth

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	pathConfigDelegationUserProp    = "delegation_user"
	pathConfigFetchGroupsProp       = "fetch_groups"
	pathConfigClientIDProp          = "client_id"
	pathConfigRedirectURLProp       = "redirect_url"
	pathConfigClientSecretProp      = "client_secret"
	pathConfigServiceAccountKeyProp = "service_acc_key"
	pathConfigEntry                 = "config"
	pathConfigPattern               = "config"
)

func pathConfig(b *googleAccountAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: pathConfigPattern,
		Fields: Schema{
			pathConfigClientIDProp: {
				Type:        framework.TypeString,
				Description: "Google OAuth client id",
			},
			pathConfigClientSecretProp: {
				Type:        framework.TypeString,
				Description: "Google OAuth client secret",
			},
			pathConfigRedirectURLProp: {
				Type:        framework.TypeString,
				Description: "Google OAuth redirect URL",
			},
			pathConfigFetchGroupsProp: {
				Type:        framework.TypeBool,
				Description: "Whether Google groups should be fetched or not",
			},
			pathConfigServiceAccountKeyProp: {
				Type:        framework.TypeString,
				Description: "Google service account key content",
			},
			pathConfigDelegationUserProp: {
				Type:        framework.TypeString,
				Description: "Google delegation email address",
			},
		},
		Callbacks: ActionCallback{
			logical.UpdateOperation: b.pathConfigWrite,
			logical.ReadOperation:   b.pathConfigRead,
		},
	}
}

func (b *googleAccountAuthBackend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	gauthc := googleOAuth{
		ServiceAccount: data.Get(pathConfigServiceAccountKeyProp).(string),
		DelegationUser: data.Get(pathConfigDelegationUserProp).(string),
	}

	if clientID, err := getRequiredStringData(data, pathConfigClientIDProp); err == nil {
		gauthc.ClientID = *clientID
	} else {
		return nil, err
	}

	if clientSecret, err := getRequiredStringData(data, pathConfigClientSecretProp); err == nil {
		gauthc.ClientSecret = *clientSecret
	} else {
		return nil, err
	}

	if redirectURL, err := getRequiredStringData(data, pathConfigRedirectURLProp); err == nil {
		url := *redirectURL
		if !isValidUrl(url) {
			return nil, fmt.Errorf("property '%s' must be a valid URL; got '%s'", pathConfigRedirectURLProp, url)
		}

		gauthc.RedirectURL = url
	} else {
		return nil, err
	}

	if fetchGroups, ok := data.GetOk(pathConfigFetchGroupsProp); ok {
		gauthc.FetchGroups = fetchGroups.(bool)
	} else {
		gauthc.FetchGroups = false
	}

	entry, err := logical.StorageEntryJSON(pathConfigEntry, gauthc)
	if err != nil {
		return nil, err
	}

	return nil, req.Storage.Put(ctx, entry)
}

func (b *googleAccountAuthBackend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	googleOAuth, err := b.getGoogleOAuthConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if googleOAuth == nil {
		return nil, nil
	}

	response := &logical.Response{
		Data: GenericMap{
			/* client secret and service account key are not returned for security reasons */
			pathConfigClientIDProp:       googleOAuth.ClientID,
			pathConfigRedirectURLProp:    googleOAuth.RedirectURL,
			pathConfigFetchGroupsProp:    googleOAuth.FetchGroups,
			pathConfigDelegationUserProp: googleOAuth.DelegationUser,
		},
	}

	return response, nil
}

// Config returns the configuration for this backend.
func (b *googleAccountAuthBackend) getGoogleOAuthConfig(ctx context.Context, s logical.Storage) (*googleOAuth, error) {
	entry, err := s.Get(ctx, pathConfigEntry)

	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	var result googleOAuth
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, fmt.Errorf("error reading configuration: %s", err)
	}

	return &result, nil
}
