package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ghoststack/ghoststack/internal/runtime"
)

type Server struct {
	daemon *runtime.Daemon
}

func NewServer(daemon *runtime.Daemon) *Server {
	return &Server{daemon: daemon}
}

func (s *Server) Start(ctx context.Context, addr string) (*http.Server, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/status", s.handleStatus)
	mux.HandleFunc("/api/monitoring", s.handleMonitoring)
	mux.HandleFunc("/api/logs", s.handleLogs)

	server := &http.Server{Addr: addr, Handler: mux}

	go func() {
		<-ctx.Done()
		server.Close()
	}()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// handle listen error
		}
	}()

	return server, nil
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
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

	json.NewEncoder(w).Encode(map[string]any{
		"cpu":     0,
		"memory":  0,
		"network": map[string]int{"in": 0, "out": 0},
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

	json.NewEncoder(w).Encode(logs)
}
