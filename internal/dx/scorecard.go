package dx

import "fmt"

type Scorecard struct {
	Total      int
	Dimensions map[string]int
}

func NewScorecard() *Scorecard {
	return &Scorecard{Dimensions: make(map[string]int)}
}

func (s *Scorecard) Rate(dimension string, score int) {
	if score < 0 {
		score = 0
	}
	if score > 10 {
		score = 10
	}
	s.Dimensions[dimension] = score
	s.Total = 0
	for _, v := range s.Dimensions {
		s.Total += v
	}
}

func (s *Scorecard) Average() float64 {
	if len(s.Dimensions) == 0 {
		return 0
	}
	return float64(s.Total) / float64(len(s.Dimensions))
}

func (s *Scorecard) Report() string {
	report := "DX Scorecard\n"
	for dim, score := range s.Dimensions {
		report += fmt.Sprintf("- %s: %d/10\n", dim, score)
	}
	report += fmt.Sprintf("Average: %.1f/10\n", s.Average())
	return report
}
