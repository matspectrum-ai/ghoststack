package cost

import "context"

type CostReport struct {
	Total   float64
	ByService map[string]float64
}

type Optimizer struct{}

func NewOptimizer() *Optimizer {
	return &Optimizer{}
}

func (o *Optimizer) Analyze(ctx context.Context) (CostReport, error) {
	return CostReport{Total: 0, ByService: make(map[string]float64)}, nil
}
