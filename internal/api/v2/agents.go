package v2

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ghoststack/ghoststack/internal/storage"
)

type registerRequest struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Metadata string `json:"metadata"`
}

type registerResponse struct {
	AgentID int    `json:"agent_id"`
	Status  string `json:"status"`
}

type heartbeatRequest struct {
	Status   string `json:"status"`
	Version  string `json:"version"`
	Metadata string `json:"metadata"`
}

type commandResultRequest struct {
	Status string `json:"status"`
	Result string `json:"result"`
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		writeError(w, "name required", http.StatusBadRequest)
		return
	}

	now := time.Now().Unix()
	id, err := s.store.SaveAgent(r.Context(), storage.Agent{
		Name:          req.Name,
		ControllerURL: r.Host,
		APIKeyHash:    "",
		Status:        "online",
		LastHeartbeat: now,
		Version:       req.Version,
		Metadata:      req.Metadata,
	})
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, registerResponse{AgentID: id, Status: "online"})
}

func (s *Server) handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := parseAgentID(r.URL.Path)
	if err != nil {
		writeError(w, "invalid agent id", http.StatusBadRequest)
		return
	}

	var req heartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid json", http.StatusBadRequest)
		return
	}

	status := req.Status
	if status == "" {
		status = "online"
	}

	if err := s.store.UpdateAgentHeartbeat(r.Context(), id, status); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"status": "ok"})
}

func (s *Server) handlePendingCommands(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := parseAgentID(r.URL.Path)
	if err != nil {
		writeError(w, "invalid agent id", http.StatusBadRequest)
		return
	}

	cmds, err := s.store.PendingCommands(r.Context(), id)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if cmds == nil {
		cmds = []storage.Command{}
	}

	writeJSON(w, cmds)
}

func (s *Server) handleCommandResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	agentID, cmdID, err := parseAgentAndCmdID(r.URL.Path)
	if err != nil {
		writeError(w, "invalid path", http.StatusBadRequest)
		return
	}

	var req commandResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid json", http.StatusBadRequest)
		return
	}

	now := time.Now().Unix()
	if err := s.store.UpdateCommand(r.Context(), storage.Command{
		ID:         cmdID,
		AgentID:    agentID,
		Status:     req.Status,
		ExecutedAt: now,
		Result:     req.Result,
	}); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"status": "ok"})
}

func (s *Server) handleListAgents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	agents, err := s.store.ListAgents(r.Context())
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if agents == nil {
		agents = []storage.Agent{}
	}

	writeJSON(w, agents)
}

func (s *Server) handleGetAgent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := parseAgentID(r.URL.Path)
	if err != nil {
		writeError(w, "invalid agent id", http.StatusBadRequest)
		return
	}

	agent, err := s.store.LoadAgent(r.Context(), id)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if agent == nil {
		writeError(w, "agent not found", http.StatusNotFound)
		return
	}

	writeJSON(w, agent)
}

func (s *Server) handleAgentByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasSuffix(path, "/commands/pending") {
		s.handlePendingCommands(w, r)
		return
	}
	if strings.Contains(path, "/commands/") && strings.HasSuffix(path, "/result") {
		s.handleCommandResult(w, r)
		return
	}
	if strings.HasSuffix(path, "/heartbeat") {
		s.handleHeartbeat(w, r)
		return
	}
	s.handleGetAgent(w, r)
}

func parseAgentID(path string) (int, error) {
	parts := strings.Split(strings.TrimRight(path, "/"), "/")
	for i, p := range parts {
		if p == "agents" && i+1 < len(parts) {
			id, err := strconv.Atoi(parts[i+1])
			if err != nil {
				return 0, err
			}
			if parts[i+1] != "register" {
				return id, nil
			}
		}
	}
	return 0, strconv.ErrSyntax
}

func parseAgentAndCmdID(path string) (int, int, error) {
	parts := strings.Split(strings.TrimRight(path, "/"), "/")
	for i, p := range parts {
		if p == "agents" && i+2 < len(parts) {
			agentID, err := strconv.Atoi(parts[i+1])
			if err != nil {
				return 0, 0, err
			}
			if i+3 < len(parts) && parts[i+2] == "commands" {
				cmdID, err := strconv.Atoi(parts[i+3])
				if err != nil {
					return 0, 0, err
				}
				return agentID, cmdID, nil
			}
		}
	}
	return 0, 0, strconv.ErrSyntax
}
