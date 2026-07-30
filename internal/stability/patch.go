package stability

import "context"

type Patch struct {
	Version string
	Hash    string
}

type PatchManager struct {
	patches []Patch
}

func NewPatchManager() *PatchManager {
	return &PatchManager{patches: make([]Patch, 0)}
}

func (m *PatchManager) Apply(ctx context.Context, patch Patch) error {
	m.patches = append(m.patches, patch)
	return nil
}

func (m *PatchManager) List(ctx context.Context) []Patch {
	return append([]Patch(nil), m.patches...)
}
