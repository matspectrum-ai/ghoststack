package providers

import "context"

type ProviderState string

const (
	ProviderStopped   ProviderState = "stopped"
	ProviderRunning   ProviderState = "running"
	ProviderFailed    ProviderState = "failed"
)

type Provider interface {
	Name() string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	State() ProviderState
}

type ProviderFactory func(config map[string]any) (Provider, error)

type WireGuardProvider interface {
	Provider
	Status(ctx context.Context) (map[string]any, error)
}

type TorProvider interface {
	Provider
}

type SingBoxProvider interface {
	Provider
}

type UnboundProvider interface {
	Provider
}

type Socks5ProxyProvider interface {
	Provider
}
