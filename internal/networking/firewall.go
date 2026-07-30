package networking

import "context"

type Firewall interface {
	Allow(ctx context.Context, rule string) error
	Drop(ctx context.Context, rule string) error
	List(ctx context.Context) ([]string, error)
}
