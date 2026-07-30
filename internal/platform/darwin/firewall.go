package darwin

import (
	"context"

	common "github.com/ghoststack/ghoststack/internal/platform/common"
)

type Firewall struct{}

func NewFirewall() *Firewall {
	return &Firewall{}
}

func (f *Firewall) Apply(ctx context.Context, rules []string) error {
	return nil
}

func (f *Firewall) Flush(ctx context.Context) error {
	return nil
}

func Supported() bool {
	return common.Current() == common.PlatformDarwin
}
