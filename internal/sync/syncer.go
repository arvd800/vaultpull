package sync

import (
	"fmt"

	"github.com/user/vaultpull/internal/envfile"
	"github.com/user/vaultpull/internal/vault"
)

// Syncer orchestrates reading secrets from Vault and writing them to a .env file.
type Syncer struct {
	vaultClient *vault.Client
	secretPath  string
	outputPath  string
}

// New creates a new Syncer with the provided Vault connection details.
func New(vaultAddr, token, secretPath, outputPath string) (*Syncer, error) {
	client, err := vault.NewClient(vaultAddr, token)
	if err != nil {
		return nil, fmt.Errorf("syncer: create vault client: %w", err)
	}
	return &Syncer{
		vaultClient: client,
		secretPath:  secretPath,
		outputPath:  outputPath,
	}, nil
}

// Run fetches secrets from Vault and writes them to the configured output file.
func (s *Syncer) Run() error {
	secrets, err := s.vaultClient.ReadSecret(s.secretPath)
	if err != nil {
		return fmt.Errorf("syncer: read secret: %w", err)
	}

	if err := envfile.Write(s.outputPath, secrets); err != nil {
		return fmt.Errorf("syncer: write env file: %w", err)
	}

	return nil
}
