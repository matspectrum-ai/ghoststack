package networking

import (
	"context"
	"sync"
)

type Gateway struct {
	mu      sync.RWMutex
	running bool
	config  string
}

func NewGateway(config string) *Gateway {
	return &Gateway{config: config}
}

func (g *Gateway) Start(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.running {
		return ErrGatewayAlreadyStarted
	}

	g.running = true
	return nil
}

func (g *Gateway) Stop(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.running {
		return ErrGatewayNotStarted
	}

	g.running = false
	return nil
}

func (g *Gateway) Status() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if g.running {
		return "running"
	}

	return "stopped"
}
