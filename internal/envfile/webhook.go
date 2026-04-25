package envfile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookEvent represents the payload sent to a webhook endpoint.
type WebhookEvent struct {
	Event     string            `json:"event"`
	Timestamp time.Time         `json:"timestamp"`
	Secrets   map[string]string `json:"secrets,omitempty"`
	Diff      []DiffEntry       `json:"diff,omitempty"`
	Source    string            `json:"source,omitempty"`
}

// WebhookConfig holds configuration for outbound webhook delivery.
type WebhookConfig struct {
	URL        string
	Secret     string
	TimeoutSec int
}

// SendWebhook delivers a WebhookEvent as a JSON POST to the configured URL.
// It returns an error if the request fails or the server responds with a
// non-2xx status code.
func SendWebhook(cfg WebhookConfig, event WebhookEvent) error {
	if cfg.URL == "" {
		return nil
	}

	event.Timestamp = time.Now().UTC()

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	timeout := cfg.TimeoutSec
	if timeout <= 0 {
		timeout = 10
	}

	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}

	req, err := http.NewRequest(http.MethodPost, cfg.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if cfg.Secret != "" {
		req.Header.Set("X-Vaultpull-Secret", cfg.Secret)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: server returned %d", resp.StatusCode)
	}
	return nil
}
