package storage

import "context"

type RuntimeState struct {
	ID        int    `json:"id"`
	Status    string `json:"status"`
	Mode      string `json:"mode"`
	StartedAt int64  `json:"started_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type ProviderState struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	State     string `json:"state"`
	Config    string `json:"config"`
	UpdatedAt int64  `json:"updated_at"`
}

type AuditEntry struct {
	ID        int    `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Action    string `json:"action"`
	Source    string `json:"source"`
	Detail    string `json:"detail"`
}

type APIKey struct {
	ID         int    `json:"id"`
	KeyHash    string `json:"key_hash"`
	Name       string `json:"name"`
	CreatedAt  int64  `json:"created_at"`
	LastUsedAt int64  `json:"last_used_at"`
	Revoked    bool   `json:"revoked"`
}

type Agent struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	ControllerURL string `json:"controller_url"`
	APIKeyHash    string `json:"api_key_hash"`
	Status        string `json:"status"`
	LastHeartbeat int64  `json:"last_heartbeat"`
	Version       string `json:"version"`
	Metadata      string `json:"metadata"`
}

type Command struct {
	ID         int    `json:"id"`
	AgentID    int    `json:"agent_id"`
	Action     string `json:"action"`
	Payload    string `json:"payload"`
	Status     string `json:"status"`
	CreatedAt  int64  `json:"created_at"`
	ExecutedAt int64  `json:"executed_at"`
	Result     string `json:"result"`
}

type StorageProvider interface {
	Open(ctx context.Context, path string) error
	Close(ctx context.Context) error
	Migrate(ctx context.Context) error

	SaveRuntimeState(ctx context.Context, state RuntimeState) error
	LoadRuntimeState(ctx context.Context) (*RuntimeState, error)

	SaveProviderState(ctx context.Context, state ProviderState) error
	LoadProviderStates(ctx context.Context) ([]ProviderState, error)

	AppendAuditLog(ctx context.Context, entry AuditEntry) error
	QueryAuditLog(ctx context.Context, limit int) ([]AuditEntry, error)

	SaveAPIKey(ctx context.Context, key APIKey) error
	LoadAPIKeyByHash(ctx context.Context, hash string) (*APIKey, error)
	ListAPIKeys(ctx context.Context) ([]APIKey, error)
	DeleteAPIKey(ctx context.Context, id int) error
	TouchAPIKey(ctx context.Context, id int) error

	SaveAgent(ctx context.Context, agent Agent) (int, error)
	LoadAgent(ctx context.Context, id int) (*Agent, error)
	ListAgents(ctx context.Context) ([]Agent, error)
	DeleteAgent(ctx context.Context, id int) error
	UpdateAgentHeartbeat(ctx context.Context, id int, status string) error

	SaveCommand(ctx context.Context, cmd Command) (int, error)
	PendingCommands(ctx context.Context, agentID int) ([]Command, error)
	UpdateCommand(ctx context.Context, cmd Command) error
}
