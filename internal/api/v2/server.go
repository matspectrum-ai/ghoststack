package v2

import (
	"encoding/json"
	"net/http"

	"github.com/ghoststack/ghoststack/internal/storage"
)

type Server struct {
	mux   *http.ServeMux
	store storage.StorageProvider
}

func NewServer(store storage.StorageProvider) *Server {
	s := &Server{
		mux:   http.NewServeMux(),
		store: store,
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/v2/agents/register", s.handleRegister)
	s.mux.HandleFunc("/api/v2/agents/", s.handleAgentByID)
	s.mux.HandleFunc("/api/v2/agents", s.handleListAgents)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, msg string, code int) {
	http.Error(w, msg, code)
}
