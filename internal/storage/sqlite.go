package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type SQLiteProvider struct {
	db   *sql.DB
	path string
}

func NewSQLiteProvider() *SQLiteProvider {
	return &SQLiteProvider{}
}

func (p *SQLiteProvider) Open(ctx context.Context, path string) error {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("home dir: %w", err)
		}
		path = filepath.Join(home, ".ghoststack", "ghost.db")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	db, err := sql.Open("sqlite", path+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return fmt.Errorf("open sqlite: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return fmt.Errorf("ping sqlite: %w", err)
	}

	p.db = db
	p.path = path
	return nil
}

func (p *SQLiteProvider) Close(ctx context.Context) error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (p *SQLiteProvider) Migrate(ctx context.Context) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS runtime_state (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			status TEXT NOT NULL DEFAULT 'stopped',
			mode TEXT NOT NULL DEFAULT '',
			started_at INTEGER NOT NULL DEFAULT 0,
			updated_at INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS provider_states (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			state TEXT NOT NULL DEFAULT 'stopped',
			config TEXT NOT NULL DEFAULT '{}',
			updated_at INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS audit_log (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp INTEGER NOT NULL,
			action TEXT NOT NULL,
			source TEXT NOT NULL DEFAULT '',
			detail TEXT NOT NULL DEFAULT ''
		)`,
		`INSERT OR IGNORE INTO runtime_state (id, status, mode, started_at, updated_at)
		 VALUES (1, 'stopped', '', 0, ?)`,
	}

	for i, m := range migrations {
		if _, err := p.db.ExecContext(ctx, m, time.Now().Unix()); err != nil {
			return fmt.Errorf("migration %d: %w", i, err)
		}
	}

	return nil
}

func (p *SQLiteProvider) SaveRuntimeState(ctx context.Context, state RuntimeState) error {
	_, err := p.db.ExecContext(ctx,
		`INSERT INTO runtime_state (id, status, mode, started_at, updated_at)
		 VALUES (1, ?, ?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET
		   status = excluded.status,
		   mode = excluded.mode,
		   started_at = excluded.started_at,
		   updated_at = excluded.updated_at`,
		state.Status, state.Mode, state.StartedAt, state.UpdatedAt,
	)
	return err
}

func (p *SQLiteProvider) LoadRuntimeState(ctx context.Context) (*RuntimeState, error) {
	row := p.db.QueryRowContext(ctx,
		`SELECT id, status, mode, started_at, updated_at FROM runtime_state WHERE id = 1`)

	var state RuntimeState
	err := row.Scan(&state.ID, &state.Status, &state.Mode, &state.StartedAt, &state.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load runtime state: %w", err)
	}
	return &state, nil
}

func (p *SQLiteProvider) SaveProviderState(ctx context.Context, state ProviderState) error {
	_, err := p.db.ExecContext(ctx,
		`INSERT INTO provider_states (name, state, config, updated_at)
		 VALUES (?, ?, ?, ?)
		 ON CONFLICT(name) DO UPDATE SET
		   state = excluded.state,
		   config = excluded.config,
		   updated_at = excluded.updated_at`,
		state.Name, state.State, state.Config, state.UpdatedAt,
	)
	return err
}

func (p *SQLiteProvider) LoadProviderStates(ctx context.Context) ([]ProviderState, error) {
	rows, err := p.db.QueryContext(ctx,
		`SELECT id, name, state, config, updated_at FROM provider_states ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("query provider states: %w", err)
	}
	defer rows.Close()

	var states []ProviderState
	for rows.Next() {
		var s ProviderState
		if err := rows.Scan(&s.ID, &s.Name, &s.State, &s.Config, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		states = append(states, s)
	}
	return states, rows.Err()
}

func (p *SQLiteProvider) AppendAuditLog(ctx context.Context, entry AuditEntry) error {
	_, err := p.db.ExecContext(ctx,
		`INSERT INTO audit_log (timestamp, action, source, detail)
		 VALUES (?, ?, ?, ?)`,
		entry.Timestamp, entry.Action, entry.Source, entry.Detail,
	)
	return err
}

func (p *SQLiteProvider) QueryAuditLog(ctx context.Context, limit int) ([]AuditEntry, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := p.db.QueryContext(ctx,
		`SELECT id, timestamp, action, source, detail
		 FROM audit_log ORDER BY id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("query audit log: %w", err)
	}
	defer rows.Close()

	var entries []AuditEntry
	for rows.Next() {
		var e AuditEntry
		if err := rows.Scan(&e.ID, &e.Timestamp, &e.Action, &e.Source, &e.Detail); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
