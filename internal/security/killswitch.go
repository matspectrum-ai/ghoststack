package security

import "context"

type KillSwitch interface {
	Enable(ctx context.Context) error
	Disable(ctx context.Context) error
	Active() bool
	Status() map[string]any
	RunLeakTest(ctx context.Context) (*LeakCheckResult, error)
}

func NewKillSwitch(iface string, opts ...KillSwitchOption) KillSwitch {
	return newKillSwitch(iface, opts...)
}
