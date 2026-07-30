package api

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ghoststack/ghoststack/internal/storage"
)

const keyPrefix = "gs_"
const keyBytes = 40

type APIKeyStore struct {
	store storage.StorageProvider
}

func NewAPIKeyStore(store storage.StorageProvider) *APIKeyStore {
	return &APIKeyStore{store: store}
}

func (ks *APIKeyStore) Generate(ctx context.Context, name string) (string, error) {
	buf := make([]byte, keyBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("rand: %w", err)
	}

	raw := keyPrefix + hex.EncodeToString(buf)
	hash := sha256Hex(raw)
	now := time.Now().Unix()

	err := ks.store.SaveAPIKey(ctx, storage.APIKey{
		KeyHash:   hash,
		Name:      name,
		CreatedAt: now,
	})
	if err != nil {
		return "", fmt.Errorf("save key: %w", err)
	}

	return raw, nil
}

func (ks *APIKeyStore) Validate(ctx context.Context, raw string) (*storage.APIKey, bool) {
	if !strings.HasPrefix(raw, keyPrefix) {
		return nil, false
	}

	hash := sha256Hex(raw)
	key, err := ks.store.LoadAPIKeyByHash(ctx, hash)
	if err != nil || key == nil || key.Revoked {
		return nil, false
	}

	ks.store.TouchAPIKey(ctx, key.ID)
	return key, true
}

func (ks *APIKeyStore) List(ctx context.Context) ([]storage.APIKey, error) {
	return ks.store.ListAPIKeys(ctx)
}

func (ks *APIKeyStore) Revoke(ctx context.Context, name string) error {
	keys, err := ks.store.ListAPIKeys(ctx)
	if err != nil {
		return err
	}
	for _, k := range keys {
		if k.Name == name {
			return ks.store.DeleteAPIKey(ctx, k.ID)
		}
	}
	return fmt.Errorf("key not found: %s", name)
}

func AuthMiddleware(store *APIKeyStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, "invalid authorization format", http.StatusUnauthorized)
				return
			}

			if _, ok := store.Validate(r.Context(), parts[1]); !ok {
				http.Error(w, "invalid or revoked api key", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
