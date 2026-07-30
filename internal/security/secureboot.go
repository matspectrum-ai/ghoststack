package security

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
	"time"
)

type SecureBoot struct {
	mu           sync.RWMutex
	expectedHash string
	secrets      map[string]string
	audit        *StructuredAuditLogger
}

func NewSecureBoot(expectedHash string) *SecureBoot {
	return &SecureBoot{
		expectedHash: expectedHash,
		secrets:      make(map[string]string),
		audit:        NewStructuredAuditLogger(""),
	}
}

func (b *SecureBoot) Verify(ctx context.Context, target string) ([]string, error) {
	var failures []string

	if b.expectedHash == "" {
		return []string{"secure boot: no expected hash configured"}, nil
	}

	data, err := os.ReadFile(target)
	if err != nil {
		return []string{fmt.Sprintf("secure boot: read %s: %v", target, err)}, nil
	}

	hash := sha256.Sum256(data)
	got := hex.EncodeToString(hash[:])

	if got != b.expectedHash {
		failures = append(failures, fmt.Sprintf(
			"secure boot: hash mismatch for %s\n  expected: %s\n  got:      %s",
			target, b.expectedHash, got))
	}

	if len(failures) > 0 {
		return failures, fmt.Errorf("secure boot: %d failure(s)", len(failures))
	}

	return nil, nil
}

func (b *SecureBoot) ComputeHash(ctx context.Context, target string) (string, error) {
	data, err := os.ReadFile(target)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", target, err)
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

type SecretRotationRecord struct {
	Timestamp time.Time
	Name      string
	PrevHash  string
	NewHash   string
}

func (b *SecureBoot) RotateSecret(ctx context.Context, name string, newValue string) error {
	if name == "" {
		return fmt.Errorf("secret name must not be empty")
	}
	if newValue == "" {
		return fmt.Errorf("new secret value must not be empty")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.secrets[name] = newValue

	b.audit.Log(ctx, StructuredAuditEntry{
		Action:    "secret_rotate",
		Source:    "secureboot",
		Detail:    fmt.Sprintf("rotated secret %s", name),
		Result:    "success",
		ID:        fmt.Sprintf("rotate-%s-%d", name, time.Now().UnixNano()),
		Timestamp: time.Now(),
	})

	return nil
}

func (b *SecureBoot) GetSecret(name string) (string, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	val, ok := b.secrets[name]
	return val, ok
}

func (b *SecureBoot) SetExpectedHash(hash string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.expectedHash = hash
}
