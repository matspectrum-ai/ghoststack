package policy

import "context"

type Policy struct {
	ID     string
	Effect string
	Rules  []string
}

type PolicyEngine struct {
	policies []Policy
}

func NewPolicyEngine() *PolicyEngine {
	return &PolicyEngine{}
}

func (e *PolicyEngine) AddPolicy(ctx context.Context, p Policy) error {
	e.policies = append(e.policies, p)
	return nil
}

func (e *PolicyEngine) Evaluate(ctx context.Context, subject, action string) bool {
	for _, p := range e.policies {
		if p.Effect == "deny" && contains(p.Rules, action) {
			return false
		}
	}
	return true
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
