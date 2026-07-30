package networking

import (
	"context"
	"fmt"
	"sync"
)

type RouteTable interface {
	Add(ctx context.Context, cidr string, gateway string) error
	Remove(ctx context.Context, cidr string) error
	List(ctx context.Context) ([]string, error)
}

var (
	ErrRouteNotFound     = fmt.Errorf("route not found")
	ErrRouteAlreadyExists = fmt.Errorf("route already exists")
)

type routeTable struct {
	mu      sync.RWMutex
	routes  map[string]string
}

func newRouteTable() *routeTable {
	return &routeTable{routes: make(map[string]string)}
}

func (r *routeTable) Add(ctx context.Context, cidr string, gateway string) error {
	if cidr == "" || gateway == "" {
		return fmt.Errorf("cidr and gateway must not be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.routes[cidr]; exists {
		return fmt.Errorf("%w: %s", ErrRouteAlreadyExists, cidr)
	}

	r.routes[cidr] = gateway
	return nil
}

func (r *routeTable) Remove(ctx context.Context, cidr string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.routes[cidr]; !exists {
		return fmt.Errorf("%w: %s", ErrRouteNotFound, cidr)
	}

	delete(r.routes, cidr)
	return nil
}

func (r *routeTable) List(ctx context.Context) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]string, 0, len(r.routes))
	for cidr, gw := range r.routes {
		out = append(out, fmt.Sprintf("%s via %s", cidr, gw))
	}
	return out, nil
}
