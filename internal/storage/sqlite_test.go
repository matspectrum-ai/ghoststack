package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func newTestProvider(t *testing.T) (*SQLiteProvider, string) {
	t.Helper()
	dir, err := os.MkdirTemp("", "ghoststack-storage-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	path := filepath.Join(dir, "test.db")
	p := NewSQLiteProvider()
	if err := p.Open(context.Background(), path); err != nil {
		t.Fatal(err)
	}
	if err := p.Migrate(context.Background()); err != nil {
		t.Fatal(err)
	}
	return p, path
}

func TestOpenClose(t *testing.T) {
	p, _ := newTestProvider(t)
	if err := p.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestMigrateIdempotent(t *testing.T) {
	p, _ := newTestProvider(t)
	if err := p.Migrate(context.Background()); err != nil {
		t.Fatal("second migrate failed:", err)
	}
}

func TestSaveLoadRuntimeState(t *testing.T) {
	p, _ := newTestProvider(t)
	defer p.Close(context.Background())
	ctx := context.Background()

	state := RuntimeState{
		Status:    "running",
		Mode:      "vpn",
		StartedAt: 1000,
		UpdatedAt: 2000,
	}

	if err := p.SaveRuntimeState(ctx, state); err != nil {
		t.Fatal("save:", err)
	}

	loaded, err := p.LoadRuntimeState(ctx)
	if err != nil {
		t.Fatal("load:", err)
	}
	if loaded == nil {
		t.Fatal("expected non-nil")
	}
	if loaded.Status != "running" || loaded.Mode != "vpn" {
		t.Fatalf("got status=%q mode=%q", loaded.Status, loaded.Mode)
	}
}

func TestLoadRuntimeStateEmpty(t *testing.T) {
	p, _ := newTestProvider(t)
	defer p.Close(context.Background())

	state, err := p.LoadRuntimeState(context.Background())
	if err != nil {
		t.Fatal("load:", err)
	}
	if state == nil {
		t.Fatal("expected default state (inserted by migration), got nil")
	}
}

func TestSaveLoadProviderStates(t *testing.T) {
	p, _ := newTestProvider(t)
	defer p.Close(context.Background())
	ctx := context.Background()

	states := []ProviderState{
		{Name: "wireguard", State: "running", Config: `{"port":51820}`, UpdatedAt: time.Now().Unix()},
		{Name: "tor", State: "stopped", Config: `{}`, UpdatedAt: time.Now().Unix()},
	}
	for _, s := range states {
		if err := p.SaveProviderState(ctx, s); err != nil {
			t.Fatal("save:", err)
		}
	}

	loaded, err := p.LoadProviderStates(ctx)
	if err != nil {
		t.Fatal("load:", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 states, got %d", len(loaded))
	}
}

func TestSaveProviderStateUpdate(t *testing.T) {
	p, _ := newTestProvider(t)
	defer p.Close(context.Background())
	ctx := context.Background()

	if err := p.SaveProviderState(ctx, ProviderState{Name: "wg", State: "running", Config: `{}`, UpdatedAt: 1}); err != nil {
		t.Fatal("save:", err)
	}
	if err := p.SaveProviderState(ctx, ProviderState{Name: "wg", State: "stopped", Config: `{}`, UpdatedAt: 2}); err != nil {
		t.Fatal("update:", err)
	}

	states, err := p.LoadProviderStates(ctx)
	if err != nil {
		t.Fatal("load:", err)
	}
	if len(states) != 1 || states[0].State != "stopped" {
		t.Fatalf("expected stopped, got %+v", states)
	}
}

func TestAuditLogAppendQuery(t *testing.T) {
	p, _ := newTestProvider(t)
	defer p.Close(context.Background())
	ctx := context.Background()

	entries := []AuditEntry{
		{Timestamp: 100, Action: "start", Source: "cli", Detail: "started daemon"},
		{Timestamp: 200, Action: "stop", Source: "cli", Detail: "stopped daemon"},
	}
	for _, e := range entries {
		if err := p.AppendAuditLog(ctx, e); err != nil {
			t.Fatal("append:", err)
		}
	}

	loaded, err := p.QueryAuditLog(ctx, 10)
	if err != nil {
		t.Fatal("query:", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(loaded))
	}
}

func TestAuditLogLimit(t *testing.T) {
	p, _ := newTestProvider(t)
	defer p.Close(context.Background())
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		_ = p.AppendAuditLog(ctx, AuditEntry{Timestamp: int64(i), Action: "test", Detail: ""})
	}

	loaded, err := p.QueryAuditLog(ctx, 3)
	if err != nil {
		t.Fatal("query:", err)
	}
	if len(loaded) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(loaded))
	}
}
