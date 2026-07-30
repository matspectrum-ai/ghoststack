package networking

import (
	"context"
	"fmt"
	"sync"
)

type Resolver interface {
	SetServers(ctx context.Context, servers []string) error
	FlushCache(ctx context.Context) error
}

var (
	ErrResolverNotStarted = fmt.Errorf("resolver not started")
)

type resolver struct {
	mu       sync.RWMutex
	servers  []string
	started  bool
	flushed  int
}

func newResolver() Resolver {
	return &resolver{}
}

func (r *resolver) SetServers(ctx context.Context, servers []string) error {
	if len(servers) == 0 {
		return fmt.Errorf("servers must not be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.servers = append([]string(nil), servers...)
	return nil
}

func (r *resolver) FlushCache(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.started {
		return ErrResolverNotStarted
	}

	r.flushed++
	return nil
}
