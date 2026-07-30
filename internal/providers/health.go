package providers

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type HealthCheckResult struct {
	Provider  string
	Healthy   bool
	Message   string
	CheckedAt time.Time
}

type HealthChecker struct {
	mu        sync.RWMutex
	providers map[string]HealthCheckFunc
	results   map[string]HealthCheckResult
}

type HealthCheckFunc func(ctx context.Context) (bool, string)

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		providers: make(map[string]HealthCheckFunc),
		results:   make(map[string]HealthCheckResult),
	}
}

func (h *HealthChecker) Register(name string, check HealthCheckFunc) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.providers[name] = check
}

func (h *HealthChecker) Check(ctx context.Context, name string) (HealthCheckResult, error) {
	h.mu.RLock()
	check, ok := h.providers[name]
	h.mu.RUnlock()

	if !ok {
		return HealthCheckResult{}, fmt.Errorf("provider not registered: %s", name)
	}

	healthy, message := check(ctx)
	result := HealthCheckResult{
		Provider:  name,
		Healthy:   healthy,
		Message:   message,
		CheckedAt: time.Now(),
	}

	h.mu.Lock()
	h.results[name] = result
	h.mu.Unlock()

	return result, nil
}

func (h *HealthChecker) CheckAll(ctx context.Context) ([]HealthCheckResult, error) {
	h.mu.RLock()
	providers := make([]string, 0, len(h.providers))
	for name := range h.providers {
		providers = append(providers, name)
	}
	h.mu.RUnlock()

	var results []HealthCheckResult
	for _, name := range providers {
		result, err := h.Check(ctx, name)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (h *HealthChecker) LastResult(name string) (HealthCheckResult, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	result, ok := h.results[name]
	return result, ok
}

func (h *HealthChecker) Fallback(ctx context.Context, primary, fallback string) (string, error) {
	result, err := h.Check(ctx, primary)
	if err != nil {
		return fallback, err
	}

	if result.Healthy {
		return primary, nil
	}

	if _, ok := h.providers[fallback]; !ok {
		return "", fmt.Errorf("fallback provider not registered: %s", fallback)
	}

	return fallback, nil
}
