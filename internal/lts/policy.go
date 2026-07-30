package lts

import (
	"context"
	"fmt"
)

type SupportPolicy struct {
	Version  string
	EOL      string
	Backport bool
}

type PolicyManager struct {
	policies []SupportPolicy
}

func NewPolicyManager() *PolicyManager {
	return &PolicyManager{policies: make([]SupportPolicy, 0)}
}

func (m *PolicyManager) AddPolicy(ctx context.Context, policy SupportPolicy) error {
	m.policies = append(m.policies, policy)
	return nil
}

func (m *PolicyManager) Current() string {
	return "v0.7.0"
}

func (m *PolicyManager) Report() string {
	report := "LTS Policies\n"
	for _, p := range m.policies {
		report += fmt.Sprintf("- %s: EOL=%s backport=%v\n", p.Version, p.EOL, p.Backport)
	}
	return report
}
