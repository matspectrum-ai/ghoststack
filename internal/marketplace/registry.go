package marketplace

import (
	"context"
	"fmt"
	"sync"
)

type PluginPackage struct {
	ID        string
	Name      string
	Version   string
	Signature string
}

type Registry struct {
	mu       sync.RWMutex
	packages map[string]PluginPackage
}

func NewRegistry() *Registry {
	return &Registry{packages: make(map[string]PluginPackage)}
}

func (r *Registry) Publish(ctx context.Context, pkg PluginPackage) error {
	if pkg.ID == "" || pkg.Version == "" {
		return fmt.Errorf("package id and version must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.packages[pkg.ID] = pkg
	return nil
}

func (r *Registry) Install(ctx context.Context, id, version string) (PluginPackage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	pkg, ok := r.packages[id]
	if !ok {
		return PluginPackage{}, fmt.Errorf("package not found: %s", id)
	}
	return pkg, nil
}
