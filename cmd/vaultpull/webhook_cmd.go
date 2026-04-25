package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/your-org/vaultpull/internal/envfile"
	syncer "github.com/your-org/vaultpull/internal/sync"
)

// runWebhook sends a test or manual webhook event to the configured endpoint.
// Usage: vaultpull webhook --url <url> [--secret <secret>] [--source <source>]
func runWebhook(args []string) error {
	fs := flag.NewFlagSet("webhook", flag.ContinueOnError)
	url := fs.String("url", "", "webhook endpoint URL (required)")
	secret := fs.String("secret", "", "optional shared secret header value")
	source := fs.String("source", "manual", "event source label")
	timeout := fs.Int("timeout", 10, "HTTP timeout in seconds")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *url == "" {
		return fmt.Errorf("webhook: --url is required")
	}

	cfg := syncer.WebhookDispatchConfig{
		Webhook: envfile.WebhookConfig{
			URL:        *url,
			Secret:     *secret,
			TimeoutSec: *timeout,
		},
		Source: *source,
		Out:    os.Stdout,
	}

	if err := syncer.DispatchSyncEvent(cfg, nil, nil); err != nil {
		return fmt.Errorf("webhook dispatch failed: %w", err)
	}
	return nil
}
