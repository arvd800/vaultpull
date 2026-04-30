package sync

import (
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultpull/internal/envfile"
)

// ApplyCondense loads a condense config from configPath and applies it to
// secrets. If configPath is empty, secrets are returned unchanged.
// Warnings about skipped empty rules are written to w (defaults to os.Stderr).
func ApplyCondense(secrets map[string]string, configPath string, w io.Writer) (map[string]string, error) {
	if w == nil {
		w = os.Stderr
	}
	if configPath == "" {
		return secrets, nil
	}
	cfg, err := envfile.LoadCondenseConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("condenser: load config: %w", err)
	}
	if len(cfg.Rules) == 0 {
		fmt.Fprintln(w, "condenser: no rules defined, skipping")
		return secrets, nil
	}
	result, err := envfile.Condense(secrets, cfg)
	if err != nil {
		return nil, fmt.Errorf("condenser: apply: %w", err)
	}
	outKeys := envfile.ListCondenseOutputKeys(cfg)
	for _, k := range outKeys {
		if _, ok := result[k]; ok {
			fmt.Fprintf(w, "condenser: produced key %q\n", k)
		}
	}
	return result, nil
}
