package gaccauth

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const backendHelp = `
This credential provider allows you to authenticate with Google accounts.
Documentation can be found at <https://github.com/caian-org/vault-plugin-auth-google-acc>
`

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
		Help:        backendHelp,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				pathLoginPattern,
				pathCodeUrlPattern,
			},
		},
		Paths: framework.PathAppend(
			pathRoles(b),
			[]*framework.Path{
				pathConfig(b),
				pathLogin(b),
				pathCodeUrl(b),
			},
		),
	}

	return b
}
