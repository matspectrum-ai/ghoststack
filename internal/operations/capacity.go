package operations

import "context"

type CapacityReport struct {
	CPU     float64
	Memory  float64
	Network float64
}

type CapacityPlanner struct{}

func NewCapacityPlanner() *CapacityPlanner {
	return &CapacityPlanner{}
}

func (p *CapacityPlanner) Plan(ctx context.Context) (CapacityReport, error) {
	return CapacityReport{CPU: 0, Memory: 0, Network: 0}, nil
}
