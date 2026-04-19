package envfile

import (
	"encoding/json"
	"fmt"
	"os"
)

// AliasMap maps alias names to canonical secret keys.
type AliasMap map[string]string

// SaveAliases writes an alias map to a JSON file.
func SaveAliases(path string, aliases AliasMap) error {
	if path == "" {
		return nil
	}
	data, err := json.MarshalIndent(aliases, "", "  ")
	if err != nil {
		return fmt.Errorf("alias: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadAliases reads an alias map from a JSON file.
// Returns an empty map if the file does not exist.
func LoadAliases(path string) (AliasMap, error) {
	if path == "" {
		return AliasMap{}, nil
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return AliasMap{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("alias: read: %w", err)
	}
	var aliases AliasMap
	if err := json.Unmarshal(data, &aliases); err != nil {
		return nil, fmt.Errorf("alias: unmarshal: %w", err)
	}
	return aliases, nil
}

// ApplyAliases returns a new map that includes aliased keys.
// For each alias -> canonical pair, if the canonical key exists in secrets,
// the alias key is added with the canonical value.
func ApplyAliases(secrets map[string]string, aliases AliasMap) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for alias, canonical := range aliases {
		if val, ok := secrets[canonical]; ok {
			out[alias] = val
		}
	}
	return out
}
