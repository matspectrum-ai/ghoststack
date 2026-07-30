package networking

import "context"

type RouteTable interface {
	Add(ctx context.Context, cidr string, gateway string) error
	Remove(ctx context.Context, cidr string) error
	List(ctx context.Context) ([]string, error)
}
