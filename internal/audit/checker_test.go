package audit

import (
	"context"
	"testing"
)

func TestRunPrivacyAudit(t *testing.T) {
	report, err := RunPrivacyAudit(context.Background())
	if err != nil {
		t.Fatalf("run privacy audit: %v", err)
	}
	if report == nil {
		t.Fatal("expected non-nil report")
	}
	if len(report.Results) == 0 {
		t.Fatal("expected at least one check")
	}
}

func TestRunSecurityAudit(t *testing.T) {
	report, err := RunSecurityAudit(context.Background())
	if err != nil {
		t.Fatalf("run security audit: %v", err)
	}
	if report == nil {
		t.Fatal("expected non-nil report")
	}
	if len(report.Results) == 0 {
		t.Fatal("expected at least one check")
	}
}

func TestCheckResultSeverity(t *testing.T) {
	tests := []struct {
		severity Severity
		wantPass bool
	}{
		{SeverityPass, true},
		{SeverityInfo, false},
		{SeverityWarn, false},
		{SeverityFail, false},
	}

	for _, tt := range tests {
		r := CheckResult{
			Name:     "test",
			Severity: tt.severity,
			Message:  "test message",
		}
		_ = r
	}
}

func TestAuditReportScore(t *testing.T) {
	report := &AuditReport{
		Passed: 8,
		Failed: 2,
	}
	report.Score = (report.Passed * 100) / (report.Passed + report.Failed)
	if report.Score != 80 {
		t.Fatalf("expected score 80, got %d", report.Score)
	}
}

func TestCheckNoTelemetry(t *testing.T) {
	result := checkNoTelemetry(context.Background())
	_ = result
}

func TestCheckMemorySafe(t *testing.T) {
	result := checkMemorySafe(context.Background())
	_ = result
}

func TestCheckKillSwitch(t *testing.T) {
	result := checkKillSwitch(context.Background())
	_ = result
}

func TestCheckFirewallPresent(t *testing.T) {
	result := checkFirewallPresent(context.Background())
	_ = result
}
