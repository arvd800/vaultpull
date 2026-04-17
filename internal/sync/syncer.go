package sync

import (
	"fmt"
	"log"

	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/envfile"
	"github.com/your-org/vaultpull/internal/vault"
)

// Syncer orchestrates fetching secrets from Vault and writing them to a .env file.
type Syncer struct {
	cfg    *config.Config
	client *vault.Client
}

// New creates a Syncer from the provided config.
func New(cfg *config.Config) (*Syncer, error) {
	c, err := vault.NewClient(cfg.VaultAddr, cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("syncer: %w", err)
	}
	return &Syncer{cfg: cfg, client: c}, nil
}

// Run fetches secrets and writes them to the configured output path.
func (s *Syncer) Run() error {
	log.Printf("fetching secrets from %s", s.cfg.SecretPath)

	secrets, err := s.client.ReadSecret(s.cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("syncer: read secret: %w", err)
	}

	output := s.cfg.OutputFile
	if output == "" {
		output = ".env"
	}

	if err := envfile.Write(output, secrets); err != nil {
		return fmt.Errorf("syncer: write env file: %w", err)
	}

	log.Printf("wrote %d secrets to %s", len(secrets), output)
	return nil
}
