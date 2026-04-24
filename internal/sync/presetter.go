package sync

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultpull/internal/envfile"
)

// ApplyPreset loads presets from dir and applies the named preset to secrets.
// If presetName is empty, secrets are returned unchanged.
// Matched keys are logged to out (defaults to os.Stdout).
func ApplyPreset(presetName, dir string, secrets map[string]string, out io.Writer) (map[string]string, error) {
	if presetName == "" {
		return secrets, nil
	}
	if out == nil {
		out = os.Stdout
	}
	presets, err := envfile.LoadPresets(dir)
	if err != nil {
		return nil, fmt.Errorf("presetter: load: %w", err)
	}
	result, err := envfile.ApplyPreset(presetName, presets, secrets)
	if err != nil {
		return nil, fmt.Errorf("presetter: apply: %w", err)
	}
	fmt.Fprintf(out, "[preset] applied %q\n", presetName)
	for k, v := range result {
		if orig, ok := secrets[k]; !ok || orig != v {
			fmt.Fprintf(out, "[preset]   %s overridden\n", k)
		}
	}
	return result, nil
}
