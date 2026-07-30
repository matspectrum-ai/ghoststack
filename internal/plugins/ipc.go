package plugins

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ghoststack/ghoststack/internal/providers"
)

type IPCMessageType string

const (
	MsgInit      IPCMessageType = "init"
	MsgStart     IPCMessageType = "start"
	MsgStop      IPCMessageType = "stop"
	MsgHealth    IPCMessageType = "health"
	MsgEvent     IPCMessageType = "event"
	MsgLog       IPCMessageType = "log"
	MsgError     IPCMessageType = "error"
	MsgInitAck   IPCMessageType = "init_ack"
	MsgStartAck  IPCMessageType = "start_ack"
	MsgStopAck   IPCMessageType = "stop_ack"
	MsgHealthAck IPCMessageType = "health_ack"
)

type IPCMessage struct {
	Type    IPCMessageType  `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type InitPayload struct {
	PluginID     string         `json:"plugin_id"`
	Config       map[string]any `json:"config"`
	Capabilities []string       `json:"capabilities"`
}

type HealthPayload struct {
	Status string `json:"status"`
	Uptime int64  `json:"uptime"`
}

type subprocessPlugin struct {
	mu         sync.RWMutex
	manifest   PluginManifest
	state      PluginState
	proc       *providers.ProcessManager
	socketPath string
	conn       net.Conn
	reader     *bufio.Scanner
	responses  chan IPCMessage
	cancel     context.CancelFunc
}

func newSubprocessPlugin(manifest PluginManifest) *subprocessPlugin {
	return &subprocessPlugin{
		manifest:  manifest,
		state:     PluginStateDiscovered,
		proc:      providers.NewProcessManager(),
		responses: make(chan IPCMessage, 16),
	}
}

func (p *subprocessPlugin) Manifest() PluginManifest {
	return p.manifest
}

func (p *subprocessPlugin) Initialize(ctx context.Context, pc PluginContext) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != PluginStateDiscovered && p.state != PluginStateValidated {
		return fmt.Errorf("plugin %s cannot initialize from state %s", p.manifest.ID, p.state)
	}

	socketPath := filepath.Join(os.TempDir(), fmt.Sprintf("ghoststack-plugin-%s.sock", p.manifest.ID))

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return fmt.Errorf("listen socket: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	p.socketPath = socketPath

	if err := p.proc.Start(ctx, providers.ProcessConfig{
		Name: p.manifest.Entry,
		Args: []string{"--socket", socketPath, "--id", p.manifest.ID},
	}); err != nil {
		listener.Close()
		os.Remove(socketPath)
		return fmt.Errorf("start plugin process: %w", err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		p.mu.Lock()
		p.conn = conn
		p.reader = bufio.NewScanner(conn)
		p.mu.Unlock()

		for p.reader.Scan() {
			var msg IPCMessage
			if err := json.Unmarshal(p.reader.Bytes(), &msg); err != nil {
				continue
			}
			select {
			case p.responses <- msg:
			default:
			}
		}
		listener.Close()
	}()

	p.state = PluginStateInitialized
	return nil
}

func (p *subprocessPlugin) Enable(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != PluginStateInitialized {
		return fmt.Errorf("plugin %s must be initialized before enable", p.manifest.ID)
	}

	p.state = PluginStateRunning
	return nil
}

func (p *subprocessPlugin) Disable(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != PluginStateRunning {
		return nil
	}

	if p.cancel != nil {
		p.cancel()
	}

	p.proc.Stop()

	if p.conn != nil {
		p.conn.Close()
	}

	if p.socketPath != "" {
		os.Remove(p.socketPath)
	}

	p.state = PluginStateStopped
	return nil
}

func (p *subprocessPlugin) Unload(ctx context.Context) error {
	{
		p.mu.Lock()
		state := p.state
		sp := p.socketPath
		if p.cancel != nil {
			p.cancel()
		}
		p.mu.Unlock()

		if state == PluginStateRunning {
			p.Disable(ctx)
		}

		p.mu.Lock()
		if p.proc != nil {
			p.proc.Stop()
		}
		if p.conn != nil {
			p.conn.Close()
		}
		if sp != "" {
			os.Remove(sp)
		}
		p.socketPath = ""
		p.state = PluginStateRemoved
		p.mu.Unlock()
	}

	return nil
}

func (p *subprocessPlugin) Health(ctx context.Context) (*HealthPayload, error) {
	return &HealthPayload{
		Status: string(p.state),
		Uptime: time.Now().Unix(),
	}, nil
}

type subprocessPluginLoader struct{}

func (l *subprocessPluginLoader) Load(path string) (Plugin, error) {
	if path == "" {
		return nil, fmt.Errorf("plugin path must not be empty")
	}

	manifestPath := filepath.Join(path, "manifest.yaml")
	manifest, err := ParseManifest(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}

	entry := filepath.Join(path, manifest.Entry)
	if _, err := os.Stat(entry); os.IsNotExist(err) {
		return nil, fmt.Errorf("plugin entry not found: %s", entry)
	}

	plugin := newSubprocessPlugin(manifest)
	plugin.state = PluginStateValidated

	return plugin, nil
}
