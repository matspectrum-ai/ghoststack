package accessibility

import "context"

type A11yChecker struct{}

func NewA11yChecker() *A11yChecker {
	return &A11yChecker{}
}

func (c *A11yChecker) Check(ctx context.Context, target string) []string {
	return []string{"ok"}
}
