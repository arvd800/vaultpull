package main

import (
	"fmt"
	"os"

	"github.com/yourusername/vaultpull/internal/envfile"
	"github.com/yourusername/vaultpull/internal/sync"
)

// runLint reads an existing .env file and runs lint checks against it.
// Usage: vaultpull lint [--fail] [--file <path>]
func runLint(args []string) error {
	filePath := ".env"
	failOnWarnings := false

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--fail":
			failOnWarnings = true
		case "--file":
			if i+1 >= len(args) {
				return fmt.Errorf("--file requires a path argument")
			}
			i++
			filePath = args[i]
		}
	}

	secrets, err := envfile.Read(filePath)
	if err != nil {
		return fmt.Errorf("lint: reading %s: %w", filePath, err)
	}

	return sync.RunLint(secrets, sync.LintConfig{
		FailOnWarnings: failOnWarnings,
		Output:         os.Stdout,
	})
}
