package release

import "context"

type ReleaseCandidate struct {
	Version string
	Checks  []string
}

type ReleaseManager struct{}

func NewReleaseManager() *ReleaseManager {
	return &ReleaseManager{}
}

func (m *ReleaseManager) Prepare(ctx context.Context, version string) (ReleaseCandidate, error) {
	return ReleaseCandidate{
		Version: version,
		Checks:  []string{"build", "test", "lint"},
	}, nil
}
