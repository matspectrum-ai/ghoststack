package health

import (
	"context"
	"errors"
	"sync"
	"time"
)

type HealthStatus string

const (
	HealthUp       HealthStatus = "up"
	HealthDown     HealthStatus = "down"
	HealthDegraded HealthStatus = "degraded"
)

type CheckFunc func(ctx context.Context) (HealthStatus, string)

type Check struct {
	Name      string
	Status    HealthStatus
	Message   string
	CheckedAt time.Time
}

type Checker struct {
	mu     sync.RWMutex
	checks map[string]CheckFunc
}

func NewChecker() *Checker {
	return &Checker{checks: make(map[string]CheckFunc)}
}

func (c *Checker) Register(name string, fn CheckFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.checks[name] = fn
}

func (c *Checker) Run(ctx context.Context) ([]Check, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var err error
	out := make([]Check, 0, len(c.checks))
	for name, fn := range c.checks {
		status, message := fn(ctx)
		out = append(out, Check{
			Name:      name,
			Status:    status,
			Message:   message,
			CheckedAt: time.Now(),
		})
		if status == HealthDown && err == nil {
			err = errors.New("health check failed: " + name)
		}
	}
	return out, err
}
