package providers

import "context"

type Provider interface {
	Name() string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type WireGuardProvider interface {
	Name() string
	Start(ctx context.Context, configPath string) error
	Stop(ctx context.Context) error
}

type TorProvider interface {
	Name() string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type SingBoxProvider interface {
	Name() string
	Start(ctx context.Context, configPath string) error
	Stop(ctx context.Context) error
}

type UnboundProvider interface {
	Name() string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Socks5ProxyProvider interface {
	Name() string
	Start(ctx context.Context, listen string) error
	Stop(ctx context.Context) error
}
