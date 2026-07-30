package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/ghoststack/ghoststack/internal/runtime"
)

type Server struct {
	daemon  *runtime.Daemon
	hub     *WSHub
	metrics *MetricsCollector
}

func NewServer(daemon *runtime.Daemon) *Server {
	hub := NewWSHub()
	return &Server{
		daemon:  daemon,
		hub:     hub,
		metrics: NewMetricsCollector(hub),
	}
}

func (s *Server) Start(ctx context.Context, addr string) (*http.Server, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/status", s.handleStatus)
	mux.HandleFunc("/api/monitoring", s.handleMonitoring)
	mux.HandleFunc("/api/logs", s.handleLogs)
	mux.HandleFunc("/api/ws", s.hub.HandleWS)
	mux.HandleFunc("/api/events", providerSSEHandler(s.hub))
	mux.HandleFunc("/api/providers", s.handleProviders)
	mux.HandleFunc("/api/config", s.handleConfig)

	go s.metrics.Collect(ctx)

	if dashboardDir := os.Getenv("GHOSTSTACK_DASHBOARD_DIR"); dashboardDir != "" {
		mux.Handle("/", http.FileServer(http.Dir(dashboardDir)))
	}

	server := &http.Server{Addr: addr, Handler: mux}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
	}()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		}
	}()

	return server, nil
}

func (s *Server) Hub() *WSHub {
	return s.hub
}

func (s *Server) Metrics() *MetricsCollector {
	return s.metrics
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, map[string]any{
		"state":  s.daemon.State(),
		"uptime": s.daemon.Uptime().String(),
		"config": s.daemon.Config(),
	})
}

func (s *Server) handleMonitoring(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	snap := s.metrics.Snapshot()
	writeJSON(w, map[string]any{
		"cpu":     snap.CPU,
		"memory":  snap.Memory,
		"rxBytes": snap.RXBytes,
		"txBytes": snap.TXBytes,
		"uptime":  snap.Uptime,
	})
}

func (s *Server) handleLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	logs := make([]map[string]string, 0)
	for _, e := range s.daemon.Events() {
		logs = append(logs, map[string]string{
			"id":        e.ID,
			"timestamp": e.Timestamp.Format(time.RFC3339),
			"type":      e.Type,
			"source":    e.Source,
		})
	}

	writeJSON(w, logs)
}

func (s *Server) handleProviders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, map[string]any{
		"available": []string{"wireguard", "tor", "sing-box", "unbound", "socks5"},
	})
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, s.daemon.Config())
	case http.MethodPost:
		var cfg map[string]any
		if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, map[string]string{"status": "config_updated"})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
