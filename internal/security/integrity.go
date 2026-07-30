package security

import "context"

type IntegrityChecker interface {
	Verify(ctx context.Context, target string) ([]string, error)
}
