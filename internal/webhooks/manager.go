package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Webhook struct {
	ID      string
	URL     string
	Secret  string
	Events  []string
}

type WebhookManager struct {
	mu       sync.RWMutex
	webhooks map[string]Webhook
}

func NewWebhookManager() *WebhookManager {
	return &WebhookManager{webhooks: make(map[string]Webhook)}
}

func (m *WebhookManager) Register(ctx context.Context, webhook Webhook) error {
	if webhook.ID == "" {
		return fmt.Errorf("webhook id must not be empty")
	}
	if webhook.URL == "" {
		return fmt.Errorf("webhook url must not be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.webhooks[webhook.ID] = webhook
	return nil
}

func (m *WebhookManager) Dispatch(ctx context.Context, event string, payload any) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, webhook := range m.webhooks {
		if !contains(webhook.Events, event) {
			continue
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhook.URL, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Body = nil

		_ = body
		_ = req
	}

	return nil
}

func (m *WebhookManager) List(ctx context.Context) ([]Webhook, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]Webhook, 0, len(m.webhooks))
	for _, w := range m.webhooks {
		out = append(out, w)
	}
	return out, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
