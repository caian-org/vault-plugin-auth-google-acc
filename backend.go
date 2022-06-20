package gaccauth

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type ActionCallback map[logical.Operation]framework.OperationFunc

type Schema map[string]*framework.FieldSchema

type googleAccountAuthBackend struct {
	*framework.Backend
}

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := newBackend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	return b, nil
}

func newBackend() *googleAccountAuthBackend {
	b := &googleAccountAuthBackend{}

	b.Backend = &framework.Backend{
		BackendType: logical.TypeCredential,
		AuthRenew:   b.authRenew,

		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				pathLoginPattern,
				pathCodeUrlPattern,
			},
		},

		Paths: framework.PathAppend(
			[]*framework.Path{
				pathConfig(b),
				pathLogin(b),
				pathCodeUrl(b),
			},
			pathRoles(b),
		),

		Help: `
            The Google credential provider allows you to authenticate with Google.
            Documentation can be found at https://github.com/erozario/vault-auth-google.
        `,
	}

	return b
}
