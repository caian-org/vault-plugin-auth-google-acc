package gaccauth

import (
	"context"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/oauth2"
)

const (
	pathCodeUrlPattern = "code_url"
)

func pathCodeUrl(b *googleAccountAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: pathCodeUrlPattern,
		Fields:  Schema{},
		Callbacks: ActionCallback{
			logical.ReadOperation: b.pathCodeUrlRead,
		},
	}
}

func (b *googleAccountAuthBackend) pathCodeUrlRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	googleOAuth, err := b.getGoogleOAuthConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if googleOAuth == nil {
		return logical.ErrorResponse("missing Google OAuth config"), nil
	}

	response := &logical.Response{
		Data: GenericMap{
			"url": googleOAuth.build().AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce),
		},
	}

	return response, nil
}
