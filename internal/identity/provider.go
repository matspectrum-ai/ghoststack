package identity

import "context"

type IdentityProvider interface {
	Authenticate(ctx context.Context, token string) (string, error)
}

type noopIdentityProvider struct{}

func NewNoopIdentityProvider() IdentityProvider {
	return &noopIdentityProvider{}
}

func (p *noopIdentityProvider) Authenticate(ctx context.Context, token string) (string, error) {
	return token, nil
}
