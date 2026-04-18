package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

// Config holds the vaultpull configuration.
type Config struct {
	VaultAddr  string `mapstructure:"vault_addr"`
	VaultToken string `mapstructure:"vault_token"`
	SecretPath string `mapstructure:"secret_path"`
	OutputFile string `mapstructure:"output_file"`
	Namespace  string `mapstructure:"namespace"`
}

// Load reads configuration from file and environment variables.
// Configuration is resolved in the following order (highest to lowest precedence):
//  1. VAULTPULL_* environment variables
//  2. Config file (.vaultpull.yaml)
//  3. Native Vault environment variables (VAULT_ADDR, VAULT_TOKEN)
//  4. Built-in defaults
func Load(cfgFile string) (*Config, error) {
	v := viper.New()

	v.SetDefault("output_file", ".env")
	v.SetDefault("vault_addr", "http://127.0.0.1:8200")

	v.SetEnvPrefix("VAULTPULL")
	v.AutomaticEnv()

	// Also respect native Vault env vars.
	if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		v.SetDefault("vault_addr", addr)
	}
	if token := os.Getenv("VAULT_TOKEN"); token != "" {
		v.SetDefault("vault_token", token)
	}

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigName(".vaultpull")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME")
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return nil, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validate ensures all required configuration fields are populated.
func (c *Config) validate() error {
	if c.VaultAddr == "" {
		return errors.New("vault_addr is required")
	}
	if c.VaultToken == "" {
		return errors.New("vault_token is required (set VAULT_TOKEN or vault_token in config)")
	}
	if c.SecretPath == "" {
		return errors.New("secret_path is required")
	}
	return nil
}
