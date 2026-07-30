package security

import (
	"context"
	"testing"
)

func TestCheckLeaks(t *testing.T) {
	result, err := CheckLeaks(context.Background(), "ghost0")
	if err != nil {
		t.Fatalf("check leaks should not error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestCheckDNSLeak(t *testing.T) {
	result := checkDNSLeak(context.Background())
	_ = result
}

func TestCheckIPv6Leak(t *testing.T) {
	result := checkIPv6Leak(context.Background())
	_ = result
}

func TestCheckICMPLeak(t *testing.T) {
	result := checkICMPLeak(context.Background())
	_ = result
}

func TestGetPublicIP(t *testing.T) {
	ip, err := getPublicIP(context.Background())
	if err != nil {
		_ = ip
	}
}
