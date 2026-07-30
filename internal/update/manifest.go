package update

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type UpdateState string

const (
	UpdateStateIdle        UpdateState = "idle"
	UpdateStateChecking    UpdateState = "checking"
	UpdateStateDownloading UpdateState = "downloading"
	UpdateStateVerifying   UpdateState = "verifying"
	UpdateStateMigrating   UpdateState = "migrating"
	UpdateStateInstalling  UpdateState = "installing"
	UpdateStateValidating  UpdateState = "validating"
	UpdateStateCompleted   UpdateState = "completed"
	UpdateStateFailed      UpdateState = "failed"
)

type UpdateManifest struct {
	Version     string
	Requires    ComponentVersion
	Migrations  []string
	Checksum    string
	URL         string
	Size        int64
	PublishedAt string
}

type UpdateCheckResult struct {
	State     UpdateState
	Available bool
	Manifest  *UpdateManifest
	Error     error
}

func (m UpdateManifest) VerifyChecksum(data []byte) error {
	hash := sha256.Sum256(data)
	actual := fmt.Sprintf("%x", hash[:])
	if actual != m.Checksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", m.Checksum, actual)
	}
	return nil
}

func (m UpdateManifest) Validate() error {
	if m.Version == "" {
		return fmt.Errorf("version is required")
	}
	if m.Checksum == "" {
		return fmt.Errorf("checksum is required")
	}
	return nil
}

func (m UpdateManifest) String() string {
	b, _ := json.MarshalIndent(m, "", "  ")
	return string(b)
}
