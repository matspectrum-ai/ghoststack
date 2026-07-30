package networking

import "context"

type TunDevice interface {
	Create(ctx context.Context, name string, mtu int) error
	Up(ctx context.Context) error
	Down(ctx context.Context) error
	Addresses(ctx context.Context) ([]string, error)
}
