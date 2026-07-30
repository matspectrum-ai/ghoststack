package security

import (
	"context"
	"fmt"
	"sync"
)

type integrityChecker struct {
	mu        sync.RWMutex
	lastCheck string
	failures  []string
}

func newIntegrityChecker() IntegrityChecker {
	return &integrityChecker{}
}

func (i *integrityChecker) Verify(ctx context.Context, target string) ([]string, error) {
	if target == "" {
		return nil, fmt.Errorf("target must not be empty")
	}

	i.mu.Lock()
	defer i.mu.Unlock()
	i.lastCheck = target

	if len(i.failures) > 0 {
		return append([]string(nil), i.failures...), ErrIntegrityFailure
	}

	return nil, nil
}

func (i *integrityChecker) recordFailure(reason string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.failures = append(i.failures, reason)
}
