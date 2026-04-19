package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/your-org/vaultpull/internal/envfile"
	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/vault"
)

func runExport(args []string) error {
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	format := fs.String("format", "dotenv", "output format: dotenv, json, export")
	output := fs.String("out", "", "output file path (default: stdout)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := vault.ReadSecretAuto(client, cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("read secret: %w", err)
	}

	fmt := envfile.ExportFormat(*format)

	if *output != "" {
		if err := envfile.ExportToFile(secrets, fmt, *output); err != nil {
			return fmt.Errorf("export to file: %w", err)
		}
		fmt.Fprintf(os.Stdout, "exported %d keys to %s\n", len(secrets), *output)
		return nil
	}

	out, err := envfile.Export(secrets, fmt)
	if err != nil {
		return err
	}
	_, err = os.Stdout.WriteString(out)
	return err
}
