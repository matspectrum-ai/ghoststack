package security

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type StructuredAuditEntry struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Source    string    `json:"source"`
	Detail    string    `json:"detail"`
	Result    string    `json:"result"`
}

type StructuredAuditLogger struct {
	mu      sync.RWMutex
	entries []StructuredAuditEntry
	output  string
}

func NewStructuredAuditLogger(output string) *StructuredAuditLogger {
	if output == "" {
		output = "/var/log/ghoststack/audit.json"
	}
	return &StructuredAuditLogger{output: output}
}

func (l *StructuredAuditLogger) Log(ctx context.Context, entry StructuredAuditEntry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	if entry.ID == "" {
		entry.ID = fmt.Sprintf("%d", entry.Timestamp.UnixNano())
	}

	l.entries = append(l.entries, entry)

	if len(l.entries) >= 100 {
		if err := l.flush(); err != nil {
			return err
		}
	}

	return nil
}

func (l *StructuredAuditLogger) List(ctx context.Context, limit int) ([]StructuredAuditEntry, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if limit <= 0 || limit > len(l.entries) {
		limit = len(l.entries)
	}

	start := len(l.entries) - limit
	if start < 0 {
		start = 0
	}

	out := make([]StructuredAuditEntry, limit)
	copy(out, l.entries[start:])
	return out, nil
}

func (l *StructuredAuditLogger) flush() error {
	remaining := l.entries
	l.entries = nil

	f, err := os.OpenFile(l.output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		l.entries = append(l.entries, remaining...)
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	for _, entry := range remaining {
		if err := enc.Encode(entry); err != nil {
			l.entries = append(remaining, l.entries...)
			return err
		}
	}
	return nil
}
