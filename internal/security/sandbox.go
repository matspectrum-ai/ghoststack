package security

import "context"

type SandboxPolicy struct {
	AllowNetwork bool
	AllowFilesystem bool
	AllowedPaths []string
}

type Sandbox interface {
	Apply(ctx context.Context, policy SandboxPolicy) error
}

type noopSandbox struct{}

func NewNoopSandbox() Sandbox {
	return &noopSandbox{}
}

func (s *noopSandbox) Apply(ctx context.Context, policy SandboxPolicy) error {
	return nil
}
