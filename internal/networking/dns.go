package networking

import "context"

type Resolver interface {
	SetServers(ctx context.Context, servers []string) error
	FlushCache(ctx context.Context) error
}
