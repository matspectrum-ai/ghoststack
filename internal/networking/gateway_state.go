package networking

import (
	"fmt"
	"sync"
	"time"
)

var (
	ErrGatewayAlreadyStarted = fmt.Errorf("gateway already started")
	ErrGatewayNotStarted     = fmt.Errorf("gateway not started")
)

type gatewayState struct {
	mu       sync.RWMutex
	running  bool
	config   string
	started  int64
	stopped  int64
	errors   []string
}

func newGatewayState() *gatewayState {
	return &gatewayState{}
}

func (g *gatewayState) markStarted(config string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.running = true
	g.config = config
	g.started = time.Now().Unix()
}

func (g *gatewayState) markStopped() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.running = false
	g.stopped = time.Now().Unix()
}

func (g *gatewayState) addError(err error) {
	if err == nil {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.errors = append(g.errors, err.Error())
}

func (g *gatewayState) snapshot() (bool, string, int64, int64) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.running, g.config, g.started, g.stopped
}

func (g *gatewayState) configValue() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.config
}
