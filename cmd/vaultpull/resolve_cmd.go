package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/your-org/vaultpull/internal/envfile"
)

// runResolve reads a .env file, performs variable interpolation, and prints
// or writes the resolved secrets.
//
// Usage: vaultpull resolve [--file .env] [--out .env.resolved] [--allow-missing] [--fallback-env] [--format dotenv|json]
func runResolve(args []string) error {
	filePath := ".env"
	outPath := ""
	allowMissing := false
	fallbackEnv := false
	format := "dotenv"

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--file":
			i++
			if i < len(args) {
				filePath = args[i]
			}
		case "--out":
			i++
			if i < len(args) {
				outPath = args[i]
			}
		case "--allow-missing":
			allowMissing = true
		case "--fallback-env":
			fallbackEnv = true
		case "--format":
			i++
			if i < len(args) {
				format = args[i]
			}
		}
	}

	secrets, err := envfile.Read(filePath)
	if err != nil {
		return fmt.Errorf("read %q: %w", filePath, err)
	}

	resolved, err := envfile.Resolve(secrets, envfile.ResolveOptions{
		AllowMissing:  allowMissing,
		FallbackToEnv: fallbackEnv,
	})
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	if outPath != "" {
		return envfile.Write(outPath, resolved)
	}

	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(resolved)
	default:
		for _, line := range envfile.Format(resolved) {
			fmt.Println(line)
		}
	}
	return nil
}
