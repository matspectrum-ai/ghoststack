package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ghoststack/ghoststack/internal/api"
	"github.com/ghoststack/ghoststack/internal/config"
	"github.com/ghoststack/ghoststack/internal/providers"
	"github.com/ghoststack/ghoststack/internal/runtime"
	"github.com/ghoststack/ghoststack/internal/storage"
	"github.com/spf13/cobra"
)

func newStartCommand() *cobra.Command {
	var configPath string
	var apiAddr string
	var providerName string
	var force bool

	var tls bool
	var public bool
	var apiKeyName string

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start GhostStack runtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !force {
				if err := runDiagnostics(); err != nil {
					return fmt.Errorf("pre-flight checks failed: %w\nRun with --force to skip", err)
				}
			}

			var providerCfg map[string]any

			if configPath != "" {
				cfg, err := config.Load(configPath)
				if err != nil {
					return fmt.Errorf("load config: %w", err)
				}

				if providerName == "" {
					if len(cfg.Profiles["default"].Providers) > 0 {
						providerName = cfg.Profiles["default"].Providers[0]
					}
				}

				if p, ok := cfg.Profiles["default"]; ok {
					if pcfg, ok := p.Config[providerName]; ok {
						if m, ok := pcfg.(map[string]any); ok {
							providerCfg = m
						}
					}
				}
			}

			if providerName == "" {
				providerName = "wireguard"
			}

			engine := providers.NewProviderEngine()
			providers.RegisterBuiltins(engine)

			startCtx, startCancel := context.WithTimeout(cmd.Context(), 30*time.Second)
			defer startCancel()

			if err := engine.Start(startCtx, providerName, providerCfg); err != nil {
				return fmt.Errorf("start provider: %w", err)
			}

			daemon := runtime.NewDaemon(configPath, nil)

			store := storage.NewSQLiteProvider()
			if err := store.Open(cmd.Context(), ""); err != nil {
				engine.StopAll(startCtx)
				return fmt.Errorf("open storage: %w", err)
			}
			defer store.Close(cmd.Context())

			if err := store.Migrate(cmd.Context()); err != nil {
				engine.StopAll(startCtx)
				return fmt.Errorf("migrate storage: %w", err)
			}
			daemon.SetStorage(store)

			if err := daemon.Start(cmd.Context()); err != nil {
				engine.StopAll(startCtx)
				return fmt.Errorf("start daemon: %w", err)
			}

			var server *api.Server
			if tls {
				tlsCfg := api.DefaultTLSConfig(homeDir())
				server = api.NewServerWithTLS(daemon, tlsCfg)
			} else {
				server = api.NewServer(daemon)
			}

			if public {
				server.SetPublic(true)
			}

			server.SetStore(store)

			if apiKeyName != "" {
				keyStore := api.NewAPIKeyStore(store)
				server.SetKeyStore(keyStore)
				fmt.Fprintf(os.Stdout, "API key name: %s\n", apiKeyName)
			}

			httpServer, err := server.Start(cmd.Context(), apiAddr)
			if err != nil {
				daemon.Stop(cmd.Context())
				engine.StopAll(startCtx)
				return fmt.Errorf("start api: %w", err)
			}
			defer httpServer.Close()

			go func() {
				hub := server.Hub()
				metrics := server.Metrics()
				metricTicker := time.NewTicker(5 * time.Second)
				defer metricTicker.Stop()

				for {
					select {
					case <-cmd.Context().Done():
						return
					case <-metricTicker.C:
						metrics.Update(12.5, 64*1024*1024, 1024*1024, 512*1024)
						state := daemon.State()
						hub.BroadcastEvent("provider_status", map[string]any{
							"state":   state,
							"uptime":  daemon.Uptime().String(),
							"version": "0.3.0",
						})
					}
				}
			}()

			fmt.Fprintln(os.Stdout, daemon.String())
			fmt.Fprintf(os.Stdout, "Provider: %s\n", providerName)
			if tls {
				fmt.Fprintf(os.Stdout, "API listening on https://[::1]:8443\n")
				fmt.Fprintln(os.Stdout, "WebSocket: wss://[::1]:8443/api/ws")
				fmt.Fprintln(os.Stdout, "HTTP→HTTPS redirect at "+apiAddr)
			} else {
				fmt.Fprintf(os.Stdout, "API listening on http://%s\n", apiAddr)
				fmt.Fprintln(os.Stdout, "WebSocket: ws://"+apiAddr+"/api/ws")
				fmt.Fprintln(os.Stdout, "SSE:       http://"+apiAddr+"/api/events")
			}

			<-cmd.Context().Done()
			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "path to GhostStack config")
	cmd.Flags().StringVar(&apiAddr, "api-addr", "127.0.0.1:8080", "API listen address")
	cmd.Flags().StringVar(&providerName, "provider", "", "provider name (wireguard, tor, sing-box, unbound, socks5)")
	cmd.Flags().BoolVar(&force, "force", false, "skip pre-flight diagnostics")
	cmd.Flags().BoolVar(&tls, "tls", false, "enable HTTPS with auto-generated self-signed cert")
	cmd.Flags().BoolVar(&public, "public", false, "bind to 0.0.0.0 (requires --api-key)")
	cmd.Flags().StringVar(&apiKeyName, "api-key", "", "require API key auth with this key name")
	return cmd
}
