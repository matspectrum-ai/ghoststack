package final

import "context"

type ReleaseCheck struct {
	Name   string
	Passed bool
}

type ReleaseManager struct{}

func NewReleaseManager() *ReleaseManager {
	return &ReleaseManager{}
}

func (m *ReleaseManager) Prepare(ctx context.Context) ([]ReleaseCheck, error) {
	return []ReleaseCheck{
		{Name: "build", Passed: true},
		{Name: "tests", Passed: true},
	}, nil
}
