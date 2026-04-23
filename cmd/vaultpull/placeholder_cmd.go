package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/your-org/vaultpull/internal/envfile"
)

// runPlaceholder resolves PLACEHOLDER:<KEY> sentinels in a target .env file
// using values from a source .env file.
//
// Usage:
//
//	vaultpull placeholder --src <source.env> --dst <target.env> [--fail-on-unresolved] [--list]
func runPlaceholder(args []string) error {
	var (
		srcPath         string
		dstPath         string
		prefix          string
		failOnUnresolved bool
		listOnly        bool
	)

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--src":
			i++
			if i < len(args) {
				srcPath = args[i]
			}
		case "--dst":
			i++
			if i < len(args) {
				dstPath = args[i]
			}
		case "--prefix":
			i++
			if i < len(args) {
				prefix = args[i]
			}
		case "--fail-on-unresolved":
			failOnUnresolved = true
		case "--list":
			listOnly = true
		}
	}

	if dstPath == "" {
		return fmt.Errorf("--dst is required")
	}

	dst, err := envfile.Read(dstPath)
	if err != nil {
		return fmt.Errorf("reading dst: %w", err)
	}

	if listOnly {
		keys := envfile.ListPlaceholders(dst, prefix)
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(keys)
	}

	if srcPath == "" {
		return fmt.Errorf("--src is required when not using --list")
	}

	src, err := envfile.Read(srcPath)
	if err != nil {
		return fmt.Errorf("reading src: %w", err)
	}

	cfg := envfile.PlaceholderConfig{
		Prefix:          prefix,
		FailOnUnresolved: failOnUnresolved,
	}

	out, err := envfile.ResolvePlaceholders(dst, src, cfg)
	if err != nil {
		return err
	}

	if err := envfile.Write(dstPath, out); err != nil {
		return fmt.Errorf("writing resolved env: %w", err)
	}

	fmt.Fprintf(os.Stdout, "placeholder resolution complete → %s\n", dstPath)
	return nil
}
