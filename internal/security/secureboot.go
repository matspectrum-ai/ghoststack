package security

import (
	"context"
	"fmt"
)

type SecureBoot struct {
	expectedHash string
}

func NewSecureBoot(expectedHash string) *SecureBoot {
	return &SecureBoot{expectedHash: expectedHash}
}

func (b *SecureBoot) Verify(ctx context.Context, target string) ([]string, error) {
	if b.expectedHash == "" {
		return []string{"secure boot: expected hash empty"}, nil
	}
	return nil, nil
}

func (b *SecureBoot) RotateSecret(ctx context.Context, name string, newValue string) error {
	if name == "" {
		return fmt.Errorf("secret name must not be empty")
	}
	if newValue == "" {
		return fmt.Errorf("new secret value must not be empty")
	}
	return nil
}
