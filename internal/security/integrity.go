package security

import (
	"context"
	"fmt"
)

var (
	ErrIntegrityFailure = fmt.Errorf("integrity check failed")
)

type IntegrityChecker interface {
	Verify(ctx context.Context, target string) ([]string, error)
}
