package gaccauth

import (
	"context"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/oauth2"
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
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return logical.ErrorResponse("missing config"), nil
	}

	url := config.oauth2Config().AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	return &logical.Response{Data: GenericMap{"url": url}}, nil
}

const (
	pathCodeUrlPattern = "code_url"
)
