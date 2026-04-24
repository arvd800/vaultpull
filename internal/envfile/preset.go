package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Preset represents a named collection of key-value overrides.
type Preset struct {
	Name   string            `json:"name"`
	Values map[string]string `json:"values"`
}

// PresetFile holds all named presets for a project.
type PresetFile struct {
	Presets []Preset `json:"presets"`
}

const presetFileName = ".vaultpull.presets.json"

// SavePresets writes presets to a JSON file at dir.
func SavePresets(dir string, presets []Preset) error {
	if dir == "" {
		return nil
	}
	pf := PresetFile{Presets: presets}
	data, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		return fmt.Errorf("preset: marshal: %w", err)
	}
	path := filepath.Join(dir, presetFileName)
	return os.WriteFile(path, data, 0o600)
}

// LoadPresets reads presets from dir. Returns empty slice if file does not exist.
func LoadPresets(dir string) ([]Preset, error) {
	path := filepath.Join(dir, presetFileName)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Preset{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("preset: read: %w", err)
	}
	var pf PresetFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("preset: unmarshal: %w", err)
	}
	return pf.Presets, nil
}

// ApplyPreset merges the named preset's values into secrets (overwriting).
// Returns an error if the preset name is not found.
func ApplyPreset(name string, presets []Preset, secrets map[string]string) (map[string]string, error) {
	for _, p := range presets {
		if p.Name != name {
			continue
		}
		out := make(map[string]string, len(secrets)+len(p.Values))
		for k, v := range secrets {
			out[k] = v
		}
		for k, v := range p.Values {
			out[k] = v
		}
		return out, nil
	}
	return nil, fmt.Errorf("preset: %q not found", name)
}
