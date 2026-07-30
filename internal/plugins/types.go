package plugins

import "context"

type PluginState string

const (
	PluginStateDiscovered PluginState = "discovered"
	PluginStateValidated  PluginState = "validated"
	PluginStateLoaded     PluginState = "loaded"
	PluginStateInitialized PluginState = "initialized"
	PluginStateEnabled    PluginState = "enabled"
	PluginStateRunning    PluginState = "running"
	PluginStateStopped    PluginState = "stopped"
	PluginStateFailed     PluginState = "failed"
	PluginStateRemoved    PluginState = "removed"
)

type PluginManifest struct {
	ID          string
	Name        string
	Version     string
	Author      string
	SDKVersion  string
	Entry       string
	Description string
	License     string
	Capabilities []string
	Permissions  []string
	Dependencies  map[string]string
	Resources     PluginResources
}

type PluginResources struct {
	Memory string
	CPU    string
	Network bool
}

type PluginContext struct {
	Logger    Logger
	Config    ConfigAPI
	Events    EventAPI
	Storage   StorageAPI
	Secrets   SecretAPI
	Health    HealthAPI
}

type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
}

type ConfigAPI interface {
	GetPluginConfig(pluginID string) (map[string]interface{}, error)
}

type EventAPI interface {
	Publish(ctx context.Context, event interface{}) error
	Subscribe(ctx context.Context, eventType string, handler func(interface{})) error
}

type StorageAPI interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
}

type SecretAPI interface {
	GetSecret(name string) (string, error)
}

type HealthAPI interface {
	ReportHealth(status string) error
}

type Plugin interface {
	Manifest() PluginManifest
	Initialize(ctx context.Context, pc PluginContext) error
	Enable(ctx context.Context) error
	Disable(ctx context.Context) error
	Unload(ctx context.Context) error
}

type PluginLoader interface {
	Load(path string) (Plugin, error)
}

type PluginValidator interface {
	Validate(manifest PluginManifest) error
}

type PluginManager interface {
	Discover(ctx context.Context, dirs []string) ([]PluginManifest, error)
	Validate(ctx context.Context, manifest PluginManifest) error
	Load(ctx context.Context, path string) (Plugin, error)
	Initialize(ctx context.Context, plugin Plugin) error
	Enable(ctx context.Context, pluginID string) error
	Disable(ctx context.Context, pluginID string) error
	Unload(ctx context.Context, pluginID string) error
	List(ctx context.Context) ([]PluginManifest, error)
	Get(ctx context.Context, pluginID string) (Plugin, error)
}
