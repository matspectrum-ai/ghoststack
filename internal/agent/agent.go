package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/ghoststack/ghoststack/internal/storage"
)

type Runtime struct {
	controllerURL string
	apiKey        string
	agentID       int
	name          string
	client        *http.Client
	interval      time.Duration
	maxBackoff    time.Duration
	status        string
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Metadata string `json:"metadata"`
}

type RegisterResponse struct {
	AgentID int    `json:"agent_id"`
	Status  string `json:"status"`
}

type HeartbeatPayload struct {
	Status   string `json:"status"`
	Version  string `json:"version"`
	Metadata string `json:"metadata"`
}

type CommandResult struct {
	Status string `json:"status"`
	Result string `json:"result"`
}

func NewRuntime(controllerURL, apiKey, name, version string) *Runtime {
	return &Runtime{
		controllerURL: controllerURL,
		apiKey:        apiKey,
		name:          name,
		client:        &http.Client{Timeout: 30 * time.Second},
		interval:      30 * time.Second,
		maxBackoff:    60 * time.Second,
	}
}

func (r *Runtime) Run(ctx context.Context) error {
	if err := r.register(ctx); err != nil {
		return fmt.Errorf("register: %w", err)
	}

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	backoff := 1 * time.Second

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := r.heartbeat(ctx); err != nil {
				time.Sleep(backoff)
				backoff = time.Duration(math.Min(
					float64(backoff*2),
					float64(r.maxBackoff),
				))
				continue
			}
			backoff = r.interval

			cmds, err := r.pollCommands(ctx)
			if err != nil {
				continue
			}

			for _, cmd := range cmds {
				result := r.executeCommand(cmd)
				r.reportResult(ctx, cmd.ID, result)
			}
		}
	}
}

func (r *Runtime) register(ctx context.Context) error {
	body := RegisterRequest{
		Name:     r.name,
		Version:  "0.5.0",
		Metadata: "{}",
	}

	data, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST",
		r.controllerURL+"/api/v2/agents/register",
		bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.apiKey)

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("register failed: %s", string(body))
	}

	var regResp RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&regResp); err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	r.agentID = regResp.AgentID
	r.status = regResp.Status
	return nil
}

func (r *Runtime) heartbeat(ctx context.Context) error {
	payload := HeartbeatPayload{
		Status:   "running",
		Version:  "0.5.0",
		Metadata: "{}",
	}

	data, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("%s/api/v2/agents/%d/heartbeat", r.controllerURL, r.agentID),
		bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.apiKey)

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("heartbeat failed: %d", resp.StatusCode)
	}
	return nil
}

func (r *Runtime) pollCommands(ctx context.Context) ([]storage.Command, error) {
	req, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/api/v2/agents/%d/commands/pending", r.controllerURL, r.agentID),
		nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+r.apiKey)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("poll commands failed: %d", resp.StatusCode)
	}

	var cmds []storage.Command
	if err := json.NewDecoder(resp.Body).Decode(&cmds); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	jitter := time.Duration(rand.Int63n(5000)) * time.Millisecond
	time.Sleep(jitter)
	return cmds, nil
}

func (r *Runtime) executeCommand(cmd storage.Command) CommandResult {
	switch cmd.Action {
	case "status":
		return CommandResult{Status: "completed", Result: `{"state":"running"}`}
	case "ping":
		return CommandResult{Status: "completed", Result: `{"pong":true}`}
	default:
		return CommandResult{Status: "failed", Result: fmt.Sprintf("unknown action: %s", cmd.Action)}
	}
}

func (r *Runtime) reportResult(ctx context.Context, cmdID int, result CommandResult) {
	data, _ := json.Marshal(result)
	req, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("%s/api/v2/agents/%d/commands/%d/result", r.controllerURL, r.agentID, cmdID),
		bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.apiKey)

	resp, err := r.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
