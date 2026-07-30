package plugins

import (
	"context"
	"fmt"
	"sync"
)

type pluginEntry struct {
	state    PluginState
	plugin   Plugin
	manifest PluginManifest
	loader   PluginLoader
}

type pluginManager struct {
	mu           sync.RWMutex
	plugins      map[string]*pluginEntry
	validator    PluginValidator
	discoverDirs []string
}

func NewPluginManager(validator PluginValidator, discoverDirs []string) PluginManager {
	if validator == nil {
		validator = &defaultValidator{}
	}
	if discoverDirs == nil {
		discoverDirs = officialDirectories
	}
	return &pluginManager{
		plugins:      make(map[string]*pluginEntry),
		validator:    validator,
		discoverDirs: discoverDirs,
	}
}

func (m *pluginManager) Discover(ctx context.Context, dirs []string) ([]PluginManifest, error) {
	targets := dirs
	if len(targets) == 0 {
		targets = m.discoverDirs
	}

	manifestPaths, err := DiscoverPlugins(targets)
	if err != nil {
		return nil, fmt.Errorf("discover plugins: %w", err)
	}

	var manifests []PluginManifest
	for _, path := range manifestPaths {
		manifest, err := ParseManifest(path)
		if err != nil {
			continue
		}
		manifests = append(manifests, manifest)
	}

	return manifests, nil
}

func (m *pluginManager) Validate(ctx context.Context, manifest PluginManifest) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.validator.Validate(manifest)
}

func (m *pluginManager) Load(ctx context.Context, pluginPath string) (Plugin, error) {
	if pluginPath == "" {
		return nil, fmt.Errorf("plugin path must not be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.plugins[pluginPath]; exists {
		return nil, fmt.Errorf("plugin already loaded: %s", pluginPath)
	}

	manifest, err := ParseManifest(manifestPathFor(pluginPath))
	if err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}

	if err := m.validator.Validate(manifest); err != nil {
		return nil, fmt.Errorf("validate manifest: %w", err)
	}

	plugin, err := (&subprocessPluginLoader{}).Load(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("load plugin: %w", err)
	}

	m.plugins[manifest.ID] = &pluginEntry{
		state:    PluginStateLoaded,
		plugin:   plugin,
		manifest: manifest,
		loader:   &defaultPluginLoader{},
	}

	return plugin, nil
}

func (m *pluginManager) Initialize(ctx context.Context, plugin Plugin) error {
	if plugin == nil {
		return fmt.Errorf("plugin must not be nil")
	}

	manifest := plugin.Manifest()

	m.mu.Lock()
	entry, exists := m.plugins[manifest.ID]
	m.mu.Unlock()

	if !exists {
		return fmt.Errorf("plugin not loaded: %s", manifest.ID)
	}

	pc := PluginContext{}
	if err := plugin.Initialize(ctx, pc); err != nil {
		m.mu.Lock()
		entry.state = PluginStateFailed
		m.mu.Unlock()
		return fmt.Errorf("initialize plugin: %w", err)
	}

	m.mu.Lock()
	entry.state = PluginStateInitialized
	m.mu.Unlock()

	return nil
}

func (m *pluginManager) Enable(ctx context.Context, pluginID string) error {
	if pluginID == "" {
		return fmt.Errorf("plugin id must not be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, entry := range m.plugins {
		if entry.manifest.ID == pluginID {
			if err := entry.plugin.Enable(ctx); err != nil {
				entry.state = PluginStateFailed
				return fmt.Errorf("enable plugin: %w", err)
			}
			entry.state = PluginStateEnabled
			return nil
		}
	}

	return fmt.Errorf("plugin not found: %s", pluginID)
}

func (m *pluginManager) Disable(ctx context.Context, pluginID string) error {
	if pluginID == "" {
		return fmt.Errorf("plugin id must not be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, entry := range m.plugins {
		if entry.manifest.ID == pluginID {
			if err := entry.plugin.Disable(ctx); err != nil {
				return fmt.Errorf("disable plugin: %w", err)
			}
			entry.state = PluginStateStopped
			return nil
		}
	}

	return fmt.Errorf("plugin not found: %s", pluginID)
}

func (m *pluginManager) Unload(ctx context.Context, pluginID string) error {
	if pluginID == "" {
		return fmt.Errorf("plugin id must not be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, entry := range m.plugins {
		if entry.manifest.ID == pluginID {
			if err := entry.plugin.Unload(ctx); err != nil {
				return fmt.Errorf("unload plugin: %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("plugin not found: %s", pluginID)
}

func (m *pluginManager) List(ctx context.Context) ([]PluginManifest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	manifests := make([]PluginManifest, 0, len(m.plugins))
	for _, entry := range m.plugins {
		manifests = append(manifests, entry.manifest)
	}
	return manifests, nil
}

func (m *pluginManager) Get(ctx context.Context, pluginID string) (Plugin, error) {
	if pluginID == "" {
		return nil, fmt.Errorf("plugin id must not be empty")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, entry := range m.plugins {
		if entry.manifest.ID == pluginID {
			return entry.plugin, nil
		}
	}

	return nil, fmt.Errorf("plugin not found: %s", pluginID)
}

type defaultValidator struct{}

func (v *defaultValidator) Validate(manifest PluginManifest) error {
	if manifest.ID == "" {
		return fmt.Errorf("manifest.id is required")
	}
	if manifest.Version == "" {
		return fmt.Errorf("manifest.version is required")
	}
	if manifest.SDKVersion == "" {
		return fmt.Errorf("manifest.sdk is required")
	}
	if manifest.Entry == "" {
		return fmt.Errorf("manifest.entry is required")
	}
	return nil
}

type defaultPluginLoader struct{}

func (l *defaultPluginLoader) Load(path string) (Plugin, error) {
	loader := &subprocessPluginLoader{}
	return loader.Load(path)
}

func manifestPathFor(path string) string {
	return path + "/manifest.yaml"
}
