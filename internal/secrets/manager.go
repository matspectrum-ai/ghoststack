package secrets

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/crypto/argon2"
)

type SecretsManager struct {
	mu      sync.RWMutex
	store   map[string]string
	homeDir string
}

type encryptedStore struct {
	Salt       []byte `json:"salt"`
	Nonce      []byte `json:"nonce"`
	Ciphertext []byte `json:"ciphertext"`
}

func NewSecretsManager(homeDir string) *SecretsManager {
	if homeDir == "" {
		homeDir = filepath.Join(os.Getenv("HOME"), ".ghoststack")
	}
	return &SecretsManager{
		store:   make(map[string]string),
		homeDir: homeDir,
	}
}

func deriveKey(passphrase string, salt []byte) []byte {
	return argon2.IDKey([]byte(passphrase), salt, 1, 64*1024, 4, 32)
}

func (sm *SecretsManager) Init(passphrase string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	storeFile := filepath.Join(sm.homeDir, "secrets.enc")
	if _, err := os.Stat(storeFile); err == nil {
		return fmt.Errorf("secrets already initialized: %s", storeFile)
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("generate salt: %w", err)
	}

	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("generate nonce: %w", err)
	}

	key := deriveKey(passphrase, salt)

	emptyData := []byte("{}")
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("aes: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("gcm: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, emptyData, nil)

	enc := encryptedStore{Salt: salt, Nonce: nonce, Ciphertext: ciphertext}
	data, err := json.Marshal(enc)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	if err := os.MkdirAll(sm.homeDir, 0700); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	if err := os.WriteFile(storeFile, data, 0600); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	sm.store = make(map[string]string)
	return nil
}

func (sm *SecretsManager) Load(passphrase string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	storeFile := filepath.Join(sm.homeDir, "secrets.enc")
	data, err := os.ReadFile(storeFile)
	if err != nil {
		return fmt.Errorf("read secrets: %w", err)
	}

	var enc encryptedStore
	if err := json.Unmarshal(data, &enc); err != nil {
		return fmt.Errorf("parse secrets: %w", err)
	}

	key := deriveKey(passphrase, enc.Salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("aes: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("gcm: %w", err)
	}

	plaintext, err := aesgcm.Open(nil, enc.Nonce, enc.Ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decrypt: wrong passphrase or corrupted store")
	}

	var store map[string]string
	if err := json.Unmarshal(plaintext, &store); err != nil {
		return fmt.Errorf("parse store: %w", err)
	}

	sm.store = store
	return nil
}

func (sm *SecretsManager) Set(name, value string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if value == "" {
		return fmt.Errorf("secret value must not be empty")
	}

	sm.store[name] = value
	return nil
}

func (sm *SecretsManager) Get(name string) (string, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	val, ok := sm.store[name]
	if !ok {
		return "", fmt.Errorf("secret not found: %s", name)
	}
	return val, nil
}

func (sm *SecretsManager) Delete(name string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.store, name)
	return nil
}

func (sm *SecretsManager) List() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	keys := make([]string, 0, len(sm.store))
	for k := range sm.store {
		keys = append(keys, k)
	}
	return keys
}

func (sm *SecretsManager) Save(passphrase string) error {
	sm.mu.RLock()
	store := make(map[string]string, len(sm.store))
	for k, v := range sm.store {
		store[k] = v
	}
	sm.mu.RUnlock()

	storeFile := filepath.Join(sm.homeDir, "secrets.enc")
	existing, err := os.ReadFile(storeFile)
	if err != nil {
		return fmt.Errorf("read existing secrets: %w", err)
	}

	var enc encryptedStore
	if err := json.Unmarshal(existing, &enc); err != nil {
		return fmt.Errorf("parse existing: %w", err)
	}

	plaintext, err := json.Marshal(store)
	if err != nil {
		return fmt.Errorf("marshal store: %w", err)
	}

	key := deriveKey(passphrase, enc.Salt)

	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("generate nonce: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("aes: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("gcm: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	newEnc := encryptedStore{Salt: enc.Salt, Nonce: nonce, Ciphertext: ciphertext}
	out, err := json.Marshal(newEnc)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	return os.WriteFile(storeFile, out, 0600)
}

func (sm *SecretsManager) Exists(name string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	_, ok := sm.store[name]
	return ok
}

func (sm *SecretsManager) HomeDir() string {
	return sm.homeDir
}

func KeyHash(passphrase string) string {
	salt := make([]byte, 16)
	rand.Read(salt)
	key := deriveKey(passphrase, salt)
	return hex.EncodeToString(key[:8])
}
