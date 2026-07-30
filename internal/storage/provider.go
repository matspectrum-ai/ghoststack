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
}
