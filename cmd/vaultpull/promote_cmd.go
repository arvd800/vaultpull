package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/your-org/vaultpull/internal/envfile"
)

// runPromote promotes secrets from one .env file to another.
// Usage: vaultpull promote <src.env> <dst.env> [--skip-existing] [--keys KEY1,KEY2]
func runPromote(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: promote <src.env> <dst.env> [--skip-existing] [--keys KEY1,KEY2]")
	}

	srcPath := args[0]
	dstPath := args[1]

	opts := envfile.PromoteOptions{}
	for i := 2; i < len(args); i++ {
		switch {
		case args[i] == "--skip-existing":
			opts.SkipExisting = true
		case strings.HasPrefix(args[i], "--keys="):
			raw := strings.TrimPrefix(args[i], "--keys=")
			opts.Keys = strings.Split(raw, ",")
		}
	}

	src, err := envfile.Read(srcPath)
	if err != nil {
		return fmt.Errorf("reading src: %w", err)
	}

	dst, err := envfile.Read(dstPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading dst: %w", err)
	}
	if dst == nil {
		dst = map[string]string{}
	}

	out, result, err := envfile.Promote(src, dst, opts)
	if err != nil {
		return err
	}

	if err := envfile.Write(dstPath, out); err != nil {
		return fmt.Errorf("writing dst: %w", err)
	}

	fmt.Printf("promoted %d key(s), skipped %d key(s)\n", len(result.Promoted), len(result.Skipped))
	return nil
}
