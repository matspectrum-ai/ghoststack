package networking

import (
	"context"
	"testing"
)

func TestFirewallAllowDropList(t *testing.T) {
	fw := newFirewall().(*firewall)

	fw.mu.Lock()
	fw.started = true
	fw.mu.Unlock()

	if err := fw.Allow(context.Background(), ""); err == nil {
		t.Fatal("expected error for empty rule")
	}

	if err := fw.Allow(context.Background(), "allow ssh"); err != nil {
		t.Fatalf("allow: %v", err)
	}

	list, err := fw.List(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(list))
	}

	if err := fw.Drop(context.Background(), "allow ssh"); err != nil {
		t.Fatalf("drop: %v", err)
	}

	list, err = fw.List(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected 0 rules, got %d", len(list))
	}
}
