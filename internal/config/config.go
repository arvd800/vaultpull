package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all vaultpull runtime configuration.
type Config struct {
	VaultAddr   string        `yaml:"vault_addr"`
	VaultToken  string        `yaml:"vault_token"`
	SecretPath  string        `yaml:"secret_path"`
	OutputFile  string        `yaml:"output_file"`
	Backup      bool          `yaml:"backup"`
	KVVersion   int           `yaml:"kv_version"`
	Passphrase  string        `yaml:"passphrase"`
	TTL         time.Duration `yaml:"ttl"`
	TTLFile     string        `yaml:"ttl_file"`
	StripPrefix string        `yaml:"strip_prefix"`
	Include     []string      `yaml:"include"`
	Exclude     []string      `yaml:"exclude"`
}

// Load reads config from a YAML file, then overrides with environment variables.
func Load(path string) (*Config, error) {
	cfg := &Config{
		OutputFile: ".env",
		KVVersion:  2,
		TTLFile:    ".vaultpull.ttl.json",
	}

	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read config: %w", err)
		}
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parse config: %w", err)
		}
	}

	if v := os.Getenv("VAULT_ADDR"); v != "" {
		cfg.VaultAddr = v
	}
	if v := os.Getenv("VAULT_TOKEN"); v != "" {
		cfg.VaultToken = v
	}
	if v := os.Getenv("VAULTPULL_SECRET_PATH"); v != "" {
		cfg.SecretPath = v
	}
	if v := os.Getenv("VAULTPULL_OUTPUT"); v != "" {
		cfg.OutputFile = v
	}
	if v := os.Getenv("VAULTPULL_PASSPHRASE"); v != "" {
		cfg.Passphrase = v
	}

	return cfg, cfg.validate()
}

func (c *Config) validate() error {
	if c.VaultToken == "" {
		return fmt.Errorf("vault_token is required (set VAULT_TOKEN or vault_token in config)")
	}
	if c.SecretPath == "" {
		return fmt.Errorf("secret_path is required")
	}
	return nil
}
