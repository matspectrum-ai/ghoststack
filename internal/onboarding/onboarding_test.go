package onboarding

import "testing"

func TestOnboardingFlow(t *testing.T) {
	o := NewOnboarding([]Step{{ID: "1", Title: "init"}, {ID: "2", Title: "config"}})
	step, err := o.Next()
	if err != nil {
		t.Fatalf("next: %v", err)
	}
	if step.ID != "1" {
		t.Fatalf("expected step 1, got %s", step.ID)
	}
	if err := o.Complete("1"); err != nil {
		t.Fatalf("complete: %v", err)
	}
	step, _ = o.Next()
	if step.ID != "2" {
		t.Fatalf("expected step 2, got %s", step.ID)
	}
}
