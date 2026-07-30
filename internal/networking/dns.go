package networking

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrResolverNotStarted = errors.New("resolver not started")

type Resolver interface {
	SetServers(ctx context.Context, servers []string) error
	FlushCache(ctx context.Context) error
	Restore(ctx context.Context) error
}

type resolver struct {
	mu         sync.RWMutex
	servers    []string
	interface_ string
	backupPath string
}

func newResolver() Resolver {
	return &resolver{interface_: "ghost0"}
}

func (r *resolver) SetServers(ctx context.Context, servers []string) error {
	if len(servers) == 0 {
		return fmt.Errorf("servers must not be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.servers = append([]string(nil), servers...)
	return r.applyServers(ctx, servers)
}

func (r *resolver) applyServers(ctx context.Context, servers []string) error {
	err := setResolvectl(ctx, r.interface_, servers)
	if err == nil {
		return nil
	}

	return setResolvConf(ctx, servers)
}

func (r *resolver) FlushCache(ctx context.Context) error {
	if err := flushResolvectl(ctx); err != nil {
		return fmt.Errorf("flush dns cache: %w", err)
	}
	return nil
}

func (r *resolver) Restore(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.servers = nil
	return restoreResolvConf(ctx)
}
