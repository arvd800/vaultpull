package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://vault.example.com:8200")
	t.Setenv("VAULT_TOKEN", "s.testtoken")
	t.Setenv("VAULTPULL_SECRET_PATH", "secret/data/myapp")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.VaultAddr != "http://vault.example.com:8200" {
		t.Errorf("expected vault addr from env, got %q", cfg.VaultAddr)
	}
	if cfg.VaultToken != "s.testtoken" {
		t.Errorf("expected vault token from env, got %q", cfg.VaultToken)
	}
	if cfg.SecretPath != "secret/data/myapp" {
		t.Errorf("expected secret path, got %q", cfg.SecretPath)
	}
	if cfg.OutputFile != ".env" {
		t.Errorf("expected default output file .env, got %q", cfg.OutputFile)
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	content := []byte("vault_addr: http://localhost:8200\nvault_token: s.filetoken\nsecret_path: secret/data/test\noutput_file: test.env\n")
	if err := os.WriteFile(cfgPath, content, 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.OutputFile != "test.env" {
		t.Errorf("expected test.env, got %q", cfg.OutputFile)
	}
	if cfg.SecretPath != "secret/data/test" {
		t.Errorf("expected secret/data/test, got %q", cfg.SecretPath)
	}
}

func TestValidateMissingToken(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "")
	t.Setenv("VAULTPULL_VAULT_TOKEN", "")

	cfg := &Config{
		VaultAddr:  "http://localhost:8200",
		SecretPath: "secret/data/app",
	}
	if err := cfg.validate(); err == nil {
		t.Error("expected validation error for missing token")
	}
}

func TestValidateMissingSecretPath(t *testing.T) {
	cfg := &Config{
		VaultAddr:  "http://localhost:8200",
		VaultToken: "s.token",
	}
	if err := cfg.validate(); err == nil {
		t.Error("expected validation error for missing secret_path")
	}
}
