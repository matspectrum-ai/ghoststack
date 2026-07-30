package security

import "context"

type KillSwitch interface {
	Enable(ctx context.Context) error
	Disable(ctx context.Context) error
	Active() bool
}
