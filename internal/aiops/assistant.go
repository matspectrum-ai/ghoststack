package aiops

import "context"

type Assistant interface {
	Suggest(ctx context.Context, issue string) (string, error)
}

type noopAssistant struct{}

func NewNoopAssistant() Assistant {
	return &noopAssistant{}
}

func (a *noopAssistant) Suggest(ctx context.Context, issue string) (string, error) {
	return "check logs and restart", nil
}
