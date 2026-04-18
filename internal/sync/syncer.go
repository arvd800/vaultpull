package sync

import (
	"context"
	"fmt"
	"log"

	"github.com/user/vaultpull/internal/envfile"
	"github.com/user/vaultpull/internal/vault"
)

// Syncer orchestrates reading secrets from Vault and writing them to a .env file.
type Syncer struct {
	client     *vault.Client
	secretPath string
	envPath    string
	backup     bool
}

// New creates a new Syncer.
func New(client *vault.Client, secretPath, envPath string, backup bool) *Syncer {
	return &Syncer{
		client:     client,
		secretPath: secretPath,
		envPath:    envPath,
		backup:     backup,
	}
}

// Run fetches secrets and writes them to the env file, optionally showing a diff.
func (s *Syncer) Run(ctx context.Context) error {
	secrets, err := vault.ReadSecretAuto(ctx, s.client, s.secretPath)
	if err != nil {
		return fmt.Errorf("reading secret: %w", err)
	}

	existing, err := envfile.Read(s.envPath)
	if err != nil {
		existing = map[string]string{}
	}

	diff := envfile.Diff(existing, secrets)
	if !diff.HasChanges() {
		log.Println("vaultpull: no changes detected")
		return nil
	}

	logDiff(diff)

	var backupPath string
	if s.backup {
		backupPath, err = envfile.Backup(s.envPath)
		if err != nil {
			return fmt.Errorf("creating backup: %w", err)
		}
	}

	merged := envfile.Merge(existing, secrets)
	if err := envfile.Write(s.envPath, merged); err != nil {
		return fmt.Errorf("writing env file: %w", err)
	}

	if backupPath != "" {
		_ = envfile.RemoveBackup(backupPath)
	}

	return nil
}

func logDiff(d envfile.DiffResult) {
	for k := range d.Added {
		log.Printf("  + %s", k)
	}
	for k := range d.Changed {
		log.Printf("  ~ %s", k)
	}
	for k := range d.Removed {
		log.Printf("  - %s", k)
	}
}
