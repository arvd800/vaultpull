package envfile

import (
	"encoding/json"
	"fmt"
	"os"
)

// Override represents a single key override with an optional condition.
type Override struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Condition string `json:"condition,omitempty"` // e.g. "missing", "always"
}

// OverrideSet is a named collection of overrides.
type OverrideSet struct {
	Name      string     `json:"name"`
	Overrides []Override `json:"overrides"`
}

// SaveOverrides writes override sets to a JSON file.
func SaveOverrides(path string, sets []OverrideSet) error {
	if path == "" {
		return nil
	}
	data, err := json.MarshalIndent(sets, "", "  ")
	if err != nil {
		return fmt.Errorf("override: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadOverrides reads override sets from a JSON file.
func LoadOverrides(path string) ([]OverrideSet, error) {
	if path == "" {
		return nil, nil
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("override: read: %w", err)
	}
	var sets []OverrideSet
	if err := json.Unmarshal(data, &sets); err != nil {
		return nil, fmt.Errorf("override: unmarshal: %w", err)
	}
	return sets, nil
}

// ApplyOverrides applies a named override set to secrets.
// condition "missing" only sets the key if it is absent; "always" (default) overwrites.
func ApplyOverrides(secrets map[string]string, sets []OverrideSet, name string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, set := range sets {
		if set.Name != name {
			continue
		}
		for _, ov := range set.Overrides {
			switch ov.Condition {
			case "missing":
				if _, exists := out[ov.Key]; !exists {
					out[ov.Key] = ov.Value
				}
			case "", "always":
				out[ov.Key] = ov.Value
			default:
				return nil, fmt.Errorf("override: unknown condition %q for key %q", ov.Condition, ov.Key)
			}
		}
		return out, nil
	}
	return out, fmt.Errorf("override: set %q not found", name)
}
