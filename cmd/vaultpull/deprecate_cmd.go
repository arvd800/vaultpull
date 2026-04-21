package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/your-org/vaultpull/internal/envfile"
)

// runDeprecate handles the `deprecate` subcommand.
// Usage:
//
//	vaultpull deprecate --key OLD_KEY --reason "use NEW_KEY" [--replaced-by NEW_KEY] [--file .deprecations.json]
//	vaultpull deprecate --check --env .env [--file .deprecations.json]
func runDeprecate(args []string) error {
	fs := flag.NewFlagSet("deprecate", flag.ContinueOnError)

	key := fs.String("key", "", "key to mark as deprecated")
	reason := fs.String("reason", "", "reason for deprecation")
	replacedBy := fs.String("replaced-by", "", "replacement key name (optional)")
	check := fs.Bool("check", false, "check current env file for deprecated keys")
	envPath := fs.String("env", ".env", "path to .env file (used with --check)")
	filePath := fs.String("file", ".deprecations.json", "path to deprecations metadata file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *check {
		return runDeprecateCheck(*envPath, *filePath)
	}

	if *key == "" || *reason == "" {
		return fmt.Errorf("deprecate: --key and --reason are required")
	}

	dm, err := envfile.LoadDeprecations(*filePath)
	if err != nil {
		return fmt.Errorf("deprecate: load: %w", err)
	}

	dm = envfile.MarkDeprecated(dm, *key, *reason, *replacedBy)

	if err := envfile.SaveDeprecations(*filePath, dm); err != nil {
		return fmt.Errorf("deprecate: save: %w", err)
	}

	fmt.Fprintf(os.Stdout, "marked %q as deprecated\n", *key)
	return nil
}

func runDeprecateCheck(envPath, filePath string) error {
	secrets, err := envfile.Read(envPath)
	if err != nil {
		return fmt.Errorf("deprecate check: read env: %w", err)
	}

	dm, err := envfile.LoadDeprecations(filePath)
	if err != nil {
		return fmt.Errorf("deprecate check: load: %w", err)
	}

	warnings := envfile.CheckDeprecations(secrets, dm)
	if len(warnings) == 0 {
		fmt.Fprintln(os.Stdout, "no deprecated keys found")
		return nil
	}

	for _, w := range warnings {
		fmt.Fprintf(os.Stderr, "WARN: %s\n", w)
	}
	return fmt.Errorf("deprecated keys detected: %d warning(s)", len(warnings))
}
