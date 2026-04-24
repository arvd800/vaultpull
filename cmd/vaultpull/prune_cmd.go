package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/your-org/vaultpull/internal/envfile"
)

// runPrune implements the `vaultpull prune` sub-command.
// It removes unwanted keys from an existing .env file in-place.
func runPrune(args []string) error {
	fs := flag.NewFlagSet("prune", flag.ContinueOnError)
	envPath := fs.String("env", ".env", "path to the .env file")
	keepRaw := fs.String("keep", "", "comma-separated list of keys to keep (all others removed)")
	removeRaw := fs.String("remove", "", "comma-separated list of keys to remove")
	dryRun := fs.Bool("dry-run", false, "print what would be removed without writing changes")

	if err := fs.Parse(args); err != nil {
		return err
	}

	secrets, err := envfile.Read(*envPath)
	if err != nil {
		return fmt.Errorf("prune: read %s: %w", *envPath, err)
	}

	opts := envfile.PruneOptions{DryRun: *dryRun}
	if *keepRaw != "" {
		opts.KeepKeys = splitCSV(*keepRaw)
	}
	if *removeRaw != "" {
		opts.RemoveKeys = splitCSV(*removeRaw)
	}

	result, err := envfile.Prune(secrets, opts)
	if err != nil {
		return fmt.Errorf("prune: %w", err)
	}

	if len(result.Removed) == 0 {
		fmt.Fprintln(os.Stdout, "prune: nothing to remove")
		return nil
	}

	for _, k := range result.Removed {
		if *dryRun {
			fmt.Fprintf(os.Stdout, "[dry-run] would remove: %s\n", k)
		} else {
			fmt.Fprintf(os.Stdout, "removed: %s\n", k)
		}
	}

	if *dryRun {
		return nil
	}

	if err := envfile.Write(*envPath, result.Retained); err != nil {
		return fmt.Errorf("prune: write %s: %w", *envPath, err)
	}
	return nil
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
