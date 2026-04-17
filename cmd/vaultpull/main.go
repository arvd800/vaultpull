package main

import (
	"fmt"
	"os"

	"github.com/user/vaultpull/internal/config"
	"github.com/user/vaultpull/internal/sync"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	outputPath := cfg.OutputFile
	if outputPath == "" {
		outputPath = ".env"
	}

	s, err := sync.New(cfg.VaultAddr, cfg.Token, cfg.SecretPath, outputPath)
	if err != nil {
		return fmt.Errorf("init syncer: %w", err)
	}

	if err := s.Run(); err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	fmt.Printf("secrets written to %s\n", outputPath)
	return nil
}
