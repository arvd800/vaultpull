package sync

import (
	"fmt"
	"log"
	"time"

	"github.com/user/vaultpull/internal/envfile"
	"github.com/user/vaultpull/internal/vault"
)

// Options configures Syncer behaviour.
type Options struct {
	OutputFile  string
	Backup      bool
	Passphrase  string
	StripPrefix string
	Include     []string
	Exclude     []string
	TTL         time.Duration
	TTLFile     string
}

// Syncer orchestrates pulling secrets from Vault and writing them locally.
type Syncer struct {
	client *vault.Client
	opts   Options
}

// New creates a Syncer.
func New(client *vault.Client, opts Options) *Syncer {
	return &Syncer{client: client, opts: opts}
}

// Run pulls secrets from path and writes them to the output file.
func (s *Syncer) Run(secretPath string) error {
	secrets, err := vault.ReadSecretAuto(s.client, secretPath)
	if err != nil {
		return fmt.Errorf("read secret: %w", err)
	}

	if len(s.opts.Include) > 0 || s.opts.StripPrefix != "" || len(s.opts.Exclude) > 0 {
		secrets = envfile.Filter(secrets, envfile.FilterOptions{
			IncludePrefixes: s.opts.Include,
			ExcludeKeys:     s.opts.Exclude,
			StripPrefix:     s.opts.StripPrefix,
		})
	}

	if err := envfile.Validate(secrets); err != nil {
		return fmt.Errorf("validate secrets: %w", err)
	}

	existing, _ := envfile.Read(s.opts.OutputFile)
	merged := envfile.Merge(existing, secrets)

	diff := envfile.Diff(existing, merged)
	logDiff(diff)

	if s.opts.Backup {
		if bp, err := envfile.Backup(s.opts.OutputFile); err == nil {
			log.Printf("backup created: %s", bp)
		}
	}

	if s.opts.Passphrase != "" {
		if err := envfile.WriteEncryptedFile(s.opts.OutputFile, merged, s.opts.Passphrase); err != nil {
			return fmt.Errorf("write encrypted: %w", err)
		}
	} else {
		if err := envfile.Write(s.opts.OutputFile, merged); err != nil {
			return fmt.Errorf("write env: %w", err)
		}
	}

	if s.opts.TTL > 0 && s.opts.TTLFile != "" {
		if err := envfile.SaveTTL(s.opts.TTLFile, s.opts.TTL); err != nil {
			log.Printf("warn: could not save TTL record: %v", err)
		}
	}

	return nil
}

func logDiff(diff []envfile.DiffEntry) {
	for _, d := range diff {
		log.Printf("[%s] %s", d.Action, d.Key)
	}
}
