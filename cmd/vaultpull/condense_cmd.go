package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/your-org/vaultpull/internal/envfile"
)

// runCondense applies condense rules from a config file to a .env file and
// writes the result back (or prints it).
//
// Usage:
//
//	vaultpull condense --config <rules.json> --env <file> [--dry-run]
func runCondense(args []string) error {
	var configPath, envPath string
	dryRun := false

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--config":
			i++
			if i < len(args) {
				configPath = args[i]
			}
		case "--env":
			i++
			if i < len(args) {
				envPath = args[i]
			}
		case "--dry-run":
			dryRun = true
		}
	}

	if configPath == "" {
		return fmt.Errorf("condense: --config is required")
	}
	if envPath == "" {
		envPath = ".env"
	}

	secrets, err := envfile.Read(envPath)
	if err != nil {
		return fmt.Errorf("condense: read env: %w", err)
	}

	cfg, err := envfile.LoadCondenseConfig(configPath)
	if err != nil {
		return fmt.Errorf("condense: load config: %w", err)
	}

	result, err := envfile.Condense(secrets, cfg)
	if err != nil {
		return fmt.Errorf("condense: apply rules: %w", err)
	}

	if dryRun {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	if err := envfile.Write(envPath, result); err != nil {
		return fmt.Errorf("condense: write env: %w", err)
	}

	fmt.Fprintf(os.Stdout, "condense: wrote %d keys to %s\n", len(result), envPath)
	return nil
}
