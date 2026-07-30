package scaling

import "context"

type LoadBalancer struct{}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{}
}

func (l *LoadBalancer) Route(ctx context.Context, target string) error {
	return nil
}
