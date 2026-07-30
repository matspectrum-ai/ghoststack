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
		`CREATE TABLE IF NOT EXISTS api_keys (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key_hash TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL UNIQUE,
			created_at INTEGER NOT NULL,
			last_used_at INTEGER NOT NULL DEFAULT 0,
			revoked INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS agents (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			controller_url TEXT NOT NULL,
			api_key_hash TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'offline',
			last_heartbeat INTEGER NOT NULL DEFAULT 0,
			version TEXT NOT NULL DEFAULT '',
			metadata TEXT NOT NULL DEFAULT '{}'
		)`,
		`CREATE TABLE IF NOT EXISTS commands (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			agent_id INTEGER NOT NULL REFERENCES agents(id),
			action TEXT NOT NULL,
			payload TEXT NOT NULL DEFAULT '{}',
			status TEXT NOT NULL DEFAULT 'pending',
			created_at INTEGER NOT NULL,
			executed_at INTEGER NOT NULL DEFAULT 0,
			result TEXT NOT NULL DEFAULT ''
		)`,
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

func (p *SQLiteProvider) SaveAPIKey(ctx context.Context, key APIKey) error {
	_, err := p.db.ExecContext(ctx,
		`INSERT INTO api_keys (key_hash, name, created_at, last_used_at, revoked)
		 VALUES (?, ?, ?, ?, ?)`,
		key.KeyHash, key.Name, key.CreatedAt, key.LastUsedAt, boolToInt(key.Revoked),
	)
	return err
}

func (p *SQLiteProvider) LoadAPIKeyByHash(ctx context.Context, hash string) (*APIKey, error) {
	row := p.db.QueryRowContext(ctx,
		`SELECT id, key_hash, name, created_at, last_used_at, revoked FROM api_keys WHERE key_hash = ?`, hash)

	var key APIKey
	var revoked int
	err := row.Scan(&key.ID, &key.KeyHash, &key.Name, &key.CreatedAt, &key.LastUsedAt, &revoked)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load api key: %w", err)
	}
	key.Revoked = revoked != 0
	return &key, nil
}

func (p *SQLiteProvider) ListAPIKeys(ctx context.Context) ([]APIKey, error) {
	rows, err := p.db.QueryContext(ctx,
		`SELECT id, key_hash, name, created_at, last_used_at, revoked FROM api_keys ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list api keys: %w", err)
	}
	defer rows.Close()

	var keys []APIKey
	for rows.Next() {
		var k APIKey
		var revoked int
		if err := rows.Scan(&k.ID, &k.KeyHash, &k.Name, &k.CreatedAt, &k.LastUsedAt, &revoked); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		k.Revoked = revoked != 0
		keys = append(keys, k)
	}
	return keys, rows.Err()
}

func (p *SQLiteProvider) DeleteAPIKey(ctx context.Context, id int) error {
	_, err := p.db.ExecContext(ctx, `DELETE FROM api_keys WHERE id = ?`, id)
	return err
}

func (p *SQLiteProvider) TouchAPIKey(ctx context.Context, id int) error {
	now := time.Now().Unix()
	_, err := p.db.ExecContext(ctx,
		`UPDATE api_keys SET last_used_at = ? WHERE id = ?`, now, id)
	return err
}

func (p *SQLiteProvider) SaveAgent(ctx context.Context, agent Agent) (int, error) {
	res, err := p.db.ExecContext(ctx,
		`INSERT INTO agents (name, controller_url, api_key_hash, status, last_heartbeat, version, metadata)
		 VALUES (?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(name) DO UPDATE SET
		   controller_url = excluded.controller_url,
		   status = excluded.status,
		   last_heartbeat = excluded.last_heartbeat,
		   version = excluded.version,
		   metadata = excluded.metadata`,
		agent.Name, agent.ControllerURL, agent.APIKeyHash,
		agent.Status, agent.LastHeartbeat, agent.Version, agent.Metadata,
	)
	if err != nil {
		return 0, fmt.Errorf("save agent: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insert id: %w", err)
	}
	if id == 0 {
		row := p.db.QueryRowContext(ctx,
			`SELECT id FROM agents WHERE name = ?`, agent.Name)
		if err := row.Scan(&id); err != nil {
			return 0, fmt.Errorf("select agent id: %w", err)
		}
	}
	return int(id), nil
}

func (p *SQLiteProvider) LoadAgent(ctx context.Context, id int) (*Agent, error) {
	row := p.db.QueryRowContext(ctx,
		`SELECT id, name, controller_url, api_key_hash, status, last_heartbeat, version, metadata
		 FROM agents WHERE id = ?`, id)

	var a Agent
	err := row.Scan(&a.ID, &a.Name, &a.ControllerURL, &a.APIKeyHash,
		&a.Status, &a.LastHeartbeat, &a.Version, &a.Metadata)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load agent: %w", err)
	}
	return &a, nil
}

func (p *SQLiteProvider) ListAgents(ctx context.Context) ([]Agent, error) {
	rows, err := p.db.QueryContext(ctx,
		`SELECT id, name, controller_url, api_key_hash, status, last_heartbeat, version, metadata
		 FROM agents ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list agents: %w", err)
	}
	defer rows.Close()

	var agents []Agent
	for rows.Next() {
		var a Agent
		if err := rows.Scan(&a.ID, &a.Name, &a.ControllerURL, &a.APIKeyHash,
			&a.Status, &a.LastHeartbeat, &a.Version, &a.Metadata); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		agents = append(agents, a)
	}
	return agents, rows.Err()
}

func (p *SQLiteProvider) DeleteAgent(ctx context.Context, id int) error {
	_, err := p.db.ExecContext(ctx, `DELETE FROM agents WHERE id = ?`, id)
	return err
}

func (p *SQLiteProvider) UpdateAgentHeartbeat(ctx context.Context, id int, status string) error {
	now := time.Now().Unix()
	_, err := p.db.ExecContext(ctx,
		`UPDATE agents SET status = ?, last_heartbeat = ? WHERE id = ?`,
		status, now, id)
	return err
}

func (p *SQLiteProvider) SaveCommand(ctx context.Context, cmd Command) (int, error) {
	res, err := p.db.ExecContext(ctx,
		`INSERT INTO commands (agent_id, action, payload, status, created_at)
		 VALUES (?, ?, ?, 'pending', ?)`,
		cmd.AgentID, cmd.Action, cmd.Payload, cmd.CreatedAt)
	if err != nil {
		return 0, fmt.Errorf("save command: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insert id: %w", err)
	}
	return int(id), nil
}

func (p *SQLiteProvider) PendingCommands(ctx context.Context, agentID int) ([]Command, error) {
	rows, err := p.db.QueryContext(ctx,
		`SELECT id, agent_id, action, payload, status, created_at, executed_at, result
		 FROM commands WHERE agent_id = ? AND status = 'pending' ORDER BY created_at`, agentID)
	if err != nil {
		return nil, fmt.Errorf("pending commands: %w", err)
	}
	defer rows.Close()

	var cmds []Command
	for rows.Next() {
		var c Command
		if err := rows.Scan(&c.ID, &c.AgentID, &c.Action, &c.Payload, &c.Status,
			&c.CreatedAt, &c.ExecutedAt, &c.Result); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		cmds = append(cmds, c)
	}
	return cmds, rows.Err()
}

func (p *SQLiteProvider) UpdateCommand(ctx context.Context, cmd Command) error {
	_, err := p.db.ExecContext(ctx,
		`UPDATE commands SET status = ?, executed_at = ?, result = ? WHERE id = ?`,
		cmd.Status, cmd.ExecutedAt, cmd.Result, cmd.ID)
	return err
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
