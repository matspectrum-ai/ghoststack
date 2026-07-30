package security

import "context"

type AuditEntry struct {
	ID        string
	Timestamp string
	Action    string
	Source    string
	Detail    string
}

type AuditLogger interface {
	Log(ctx context.Context, entry AuditEntry) error
	List(ctx context.Context, limit int) ([]AuditEntry, error)
}
