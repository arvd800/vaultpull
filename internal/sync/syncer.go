package sync

import (
	"fmt"

	"github.com/user/vaultpull/internal/envfile"
)

// VaultReader reads secrets from Vault.
type VaultReader interface {
	ReadSecret(path string) (map[string]string, error)
}

// Syncer orchestrates pulling secrets and writing the env file.
type Syncer struct {
	vault      VaultReader
	secretPath string
	envPath    string
	backup     bool
}

// New creates a Syncer.
func New(v VaultReader, secretPath, envPath string, backup bool) *Syncer {
	return &Syncer{vault: v, secretPath: secretPath, envPath: envPath, backup: backup}
}

// Run pulls secrets and writes/merges the env file.
// If backup is enabled it creates a backup before writing.
func (s *Syncer) Run() error {
	secrets, err := s.vault.ReadSecret(s.secretPath)
	if err != nil {
		return fmt.Errorf("read secret %s: %w", s.secretPath, err)
	}

	var backupPath string
	if s.backup {
		backupPath, err = envfile.Backup(s.envPath)
		if err != nil {
			return fmt.Errorf("backup %s: %w", s.envPath, err)
		}
	}

	existing, err := envfile.Read(s.envPath)
	if err != nil {
		_ = envfile.RemoveBackup(backupPath)
		return fmt.Errorf("read env file: %w", err)
	}

	merged := envfile.Merge(existing, secrets)

	if err := envfile.Write(s.envPath, merged); err != nil {
		_ = envfile.RemoveBackup(backupPath)
		return fmt.Errorf("write env file: %w", err)
	}

	return nil
}
