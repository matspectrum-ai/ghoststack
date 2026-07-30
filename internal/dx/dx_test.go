package dx

import "testing"

func TestScorecardRateAndAverage(t *testing.T) {
	s := NewScorecard()
	s.Rate("docs", 8)
	s.Rate("api", 7)
	if s.Average() != 7.5 {
		t.Fatalf("expected 7.5, got %f", s.Average())
	}
	report := s.Report()
	if report == "" {
		t.Fatal("expected non-empty report")
	}
}

func TestScorecardClamp(t *testing.T) {
	s := NewScorecard()
	s.Rate("docs", -5)
	s.Rate("api", 20)
	if s.Dimensions["docs"] != 0 {
		t.Fatalf("expected clamped to 0")
	}
	if s.Dimensions["api"] != 10 {
		t.Fatalf("expected clamped to 10")
	}
}

func TestUserError(t *testing.T) {
	err := Wrap(nil, "E1", "msg")
	if err.Code != "E1" {
		t.Fatalf("unexpected code: %s", err.Code)
	}
	if !IsUserError(err) {
		t.Fatal("expected user error")
	}
}
