package networking

import (
	"context"
	"fmt"
	"sync"
)

type TunDevice interface {
	Create(ctx context.Context, name string, mtu int) error
	Up(ctx context.Context) error
	Down(ctx context.Context) error
	Addresses(ctx context.Context) ([]string, error)
}


var (
	ErrTUNNotCreated = fmt.Errorf("tun not created")
)

type tunDevice struct {
	mu     sync.RWMutex
	name   string
	mtu    int
	up     bool
}

func newTUN() TunDevice {
	return &tunDevice{}
}

func (t *tunDevice) Create(ctx context.Context, name string, mtu int) error {
	if name == "" {
		return fmt.Errorf("name must not be empty")
	}
	if mtu <= 0 {
		mtu = 1500
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	t.name = name
	t.mtu = mtu
	return nil
}

func (t *tunDevice) Up(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.name == "" {
		return ErrTUNNotCreated
	}

	t.up = true
	return nil
}

func (t *tunDevice) Down(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.name == "" {
		return ErrTUNNotCreated
	}

	t.up = false
	return nil
}

func (t *tunDevice) Addresses(ctx context.Context) ([]string, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.name == "" || !t.up {
		return nil, ErrTUNNotCreated
	}

	return nil, nil
}
