package security

import (
	"context"
	"sync"
	"time"
)

type auditLogger struct {
	mu      sync.RWMutex
	entries []AuditEntry
	maxSize int
}

func newAuditLogger(maxSize int) AuditLogger {
	if maxSize <= 0 {
		maxSize = 1024
	}
	return &auditLogger{maxSize: maxSize}
}

func (a *auditLogger) Log(ctx context.Context, entry AuditEntry) error {
	if entry.Action == "" {
		return nil
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if entry.Timestamp == "" {
		entry.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}

	a.entries = append(a.entries, entry)
	if len(a.entries) > a.maxSize {
		a.entries = a.entries[len(a.entries)-a.maxSize:]
	}

	return nil
}

func (a *auditLogger) List(ctx context.Context, limit int) ([]AuditEntry, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if limit <= 0 || limit > len(a.entries) {
		limit = len(a.entries)
	}

	out := make([]AuditEntry, limit)
	copy(out, a.entries[:limit])
	return out, nil
}
