package sync

import (
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultpull/internal/envfile"
)

// ApplyOverrides loads an override set by name and applies it to secrets.
// If overridePath is empty or the file does not exist, secrets are returned unchanged.
func ApplyOverrides(secrets map[string]string, overridePath, setName string, out io.Writer) (map[string]string, error) {
	if out == nil {
		out = os.Stdout
	}
	if overridePath == "" || setName == "" {
		return secrets, nil
	}
	sets, err := envfile.LoadOverrides(overridePath)
	if err != nil {
		return nil, fmt.Errorf("overrider: load: %w", err)
	}
	if len(sets) == 0 {
		return secrets, nil
	}
	result, err := envfile.ApplyOverrides(secrets, sets, setName)
	if err != nil {
		return nil, fmt.Errorf("overrider: apply: %w", err)
	}
	fmt.Fprintf(out, "overrides: applied set %q (%d keys affected)\n", setName, countChanged(secrets, result))
	return result, nil
}

func countChanged(original, updated map[string]string) int {
	n := 0
	for k, v := range updated {
		if original[k] != v {
			n++
		}
	}
	for k := range updated {
		if _, ok := original[k]; !ok {
			n++
		}
	}
	return n
}
