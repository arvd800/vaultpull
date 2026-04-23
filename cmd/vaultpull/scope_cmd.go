package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/vaultpull/internal/envfile"
)

// runScope handles the `vaultpull scope` subcommand.
// Usage:
//
//	vaultpull scope apply  <name> --scopes=<file> --env=<file>
//	vaultpull scope list              --scopes=<file>
func runScope(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("scope: subcommand required (apply|list)")
	}

	scopesFile := os.Getenv("VAULTPULL_SCOPES_FILE")
	if scopesFile == "" {
		scopesFile = ".vaultpull.scopes.json"
	}

	switch args[0] {
	case "list":
		return runScopeList(scopesFile)
	case "apply":
		if len(args) < 2 {
			return fmt.Errorf("scope apply: scope name required")
		}
		envFile := os.Getenv("VAULTPULL_ENV_FILE")
		if envFile == "" {
			envFile = ".env"
		}
		return runScopeApply(args[1], scopesFile, envFile)
	default:
		return fmt.Errorf("scope: unknown subcommand %q", args[0])
	}
}

func runScopeList(scopesFile string) error {
	scopes, err := envfile.LoadScopes(scopesFile)
	if err != nil {
		return fmt.Errorf("scope list: %w", err)
	}
	names := envfile.ListScopes(scopes)
	if len(names) == 0 {
		fmt.Println("no scopes defined")
		return nil
	}
	for _, name := range names {
		s := scopes[name]
		fmt.Printf("%-20s keys: %s\n", name, strings.Join(s.Keys, ", "))
	}
	return nil
}

func runScopeApply(name, scopesFile, envFile string) error {
	scopes, err := envfile.LoadScopes(scopesFile)
	if err != nil {
		return fmt.Errorf("scope apply: load scopes: %w", err)
	}

	secrets, err := envfile.Read(envFile)
	if err != nil {
		return fmt.Errorf("scope apply: read env: %w", err)
	}

	filtered, err := envfile.ApplyScope(secrets, scopes, name)
	if err != nil {
		return fmt.Errorf("scope apply: %w", err)
	}

	outFile := fmt.Sprintf(".env.%s", name)
	if err := envfile.Write(outFile, filtered); err != nil {
		return fmt.Errorf("scope apply: write: %w", err)
	}

	fmt.Printf("scope %q applied → %s (%d keys)\n", name, outFile, len(filtered))
	return nil
}
