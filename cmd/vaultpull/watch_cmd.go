package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/your-org/vaultpull/"
	"github."
/vault"
)

Watch starts the watch-. It exits on SIGINT/SIGTERM.
func runWatch(cfg *config.Config, interval time.Duration) error {
	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	s := sync.New(client, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("\n[watch] shutting down")
		cancel()
	}()

	return sync.WatchAndSync(ctx, s, cfg.EnvFile, interval)
}
