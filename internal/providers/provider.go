package providers

import "context"

type WireGuardProvider interface {
	Start(ctx context.Context, configPath string) error
	Stop(ctx context.Context) error
}

type TorProvider interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type SingBoxProvider interface {
	Start(ctx context.Context, configPath string) error
	Stop(ctx context.Context) error
}

type UnboundProvider interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Socks5ProxyProvider interface {
	Start(ctx context.Context, listen string) error
	Stop(ctx context.Context) error
}

type Provider interface {
	Name() string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
