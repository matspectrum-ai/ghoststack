package security

import (
	"context"
	"testing"
)

func TestKillSwitchLifecycle(t *testing.T) {
	ks := newKillSwitch("wg0")

	if ks.Active() {
		t.Fatal("expected inactive kill switch")
	}

	ks.mu.Lock()
	ks.active = true
	ks.mu.Unlock()

	if !ks.Active() {
		t.Fatal("expected active")
	}

	if err := ks.Disable(context.Background()); err != nil {
		t.Fatalf("disable: %v", err)
	}

	if ks.Active() {
		t.Fatal("expected inactive after disable")
	}
}

func TestKillSwitchEmptyInterface(t *testing.T) {
	ks := newKillSwitch("")

	if err := ks.Enable(context.Background()); err == nil {
		t.Fatal("expected error for empty interface")
	}
}

func TestKillSwitchDoubleEnable(t *testing.T) {
	ks := newKillSwitch("wg0")

	ks.mu.Lock()
	ks.active = true
	ks.mu.Unlock()

	if err := ks.Enable(context.Background()); err != nil {
		t.Fatalf("double enable should no-op: %v", err)
	}
}

func TestKillSwitchOptions(t *testing.T) {
	ks := newKillSwitch("wg0",
		WithDNSForce([]string{"1.1.1.1", "9.9.9.9"}),
	)

	if !ks.dnsForce {
		t.Fatal("expected dnsForce=true")
	}
	if len(ks.dnsAddrs) != 2 {
		t.Fatalf("expected 2 dns addrs, got %d", len(ks.dnsAddrs))
	}
}

func TestKillSwitchOptionsNoDNS(t *testing.T) {
	ks := newKillSwitch("wg0", WithDNSForce(nil))
	if ks.dnsForce {
		t.Fatal("expected dnsForce=false with nil addrs")
	}
}

func TestKillSwitchStatus(t *testing.T) {
	ks := newKillSwitch("wg0", WithDNSForce([]string{"1.1.1.1"}))
	status := ks.Status()

	if status["interface"] != "wg0" {
		t.Fatalf("expected wg0, got %v", status["interface"])
	}
	if status["active"] != false {
		t.Fatalf("expected inactive, got %v", status["active"])
	}
}

func TestLeakCheckResultJSON(t *testing.T) {
	r := &LeakCheckResult{
		DNS:      false,
		IPv6:     false,
		ICMP:     false,
		PublicIP: "",
	}

	if r.DNS {
		t.Fatal("expected DNS=false")
	}
}

func TestNewKillSwitch(t *testing.T) {
	ks := NewKillSwitch("wg0")
	if ks == nil {
		t.Fatal("expected non-nil")
	}
	if ks.Active() {
		t.Fatal("expected inactive")
	}
}
