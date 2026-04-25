package sync

import (
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultpull/internal/envfile"
)

// WebhookDispatchConfig holds settings for dispatching sync events.
type WebhookDispatchConfig struct {
	Webhook envfile.WebhookConfig
	Source  string
	Out     io.Writer
}

// DispatchSyncEvent sends a webhook event after a sync operation, including
// the computed diff entries. It logs the outcome to cfg.Out (defaults to
// os.Stdout). A missing webhook URL is silently skipped.
func DispatchSyncEvent(cfg WebhookDispatchConfig, secrets map[string]string, diff []envfile.DiffEntry) error {
	out := cfg.Out
	if out == nil {
		out = os.Stdout
	}

	if cfg.Webhook.URL == "" {
		return nil
	}

	event := envfile.WebhookEvent{
		Event:  "sync",
		Source: cfg.Source,
		Diff:   diff,
	}

	if err := envfile.SendWebhook(cfg.Webhook, event); err != nil {
		return fmt.Errorf("dispatch webhook: %w", err)
	}

	fmt.Fprintf(out, "webhook: dispatched sync event to %s (%d diff entries)\n", cfg.Webhook.URL, len(diff))
	return nil
}
