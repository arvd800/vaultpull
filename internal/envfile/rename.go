package envfile

import (
	"encoding/json"
	"fmt"
	"os"
)

// RenameRule maps an old key name to a new key name.
type RenameRule struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// RenameMap holds a list of rename rules persisted to disk.
type RenameMap struct {
	Rules []RenameRule `json:"rules"`
}

// SaveRenames writes rename rules to the given path as JSON.
func SaveRenames(path string, rm RenameMap) error {
	if path == "" {
		return nil
	}
	data, err := json.MarshalIndent(rm, "", "  ")
	if err != nil {
		return fmt.Errorf("rename: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadRenames reads rename rules from the given path.
// Returns an empty RenameMap if the file does not exist.
func LoadRenames(path string) (RenameMap, error) {
	if path == "" {
		return RenameMap{}, nil
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return RenameMap{}, nil
	}
	if err != nil {
		return RenameMap{}, fmt.Errorf("rename: read: %w", err)
	}
	var rm RenameMap
	if err := json.Unmarshal(data, &rm); err != nil {
		return RenameMap{}, fmt.Errorf("rename: unmarshal: %w", err)
	}
	return rm, nil
}

// ApplyRenames returns a new map with keys renamed according to the rules.
// The original map is not mutated. If a "from" key is not present the rule
// is silently skipped. If the "to" key already exists it is overwritten.
func ApplyRenames(secrets map[string]string, rm RenameMap) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, rule := range rm.Rules {
		if rule.From == "" || rule.To == "" {
			continue
		}
		val, ok := out[rule.From]
		if !ok {
			continue
		}
		delete(out, rule.From)
		out[rule.To] = val
	}
	return out
}
