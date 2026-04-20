package envfile

import (
	"encoding/json"
	"fmt"
	"os"
)

// InheritMap represents a parent→child inheritance relationship between env files.
type InheritMap struct {
	Parent string            `json:"parent"`
	Keys   map[string]string `json:"keys"` // child key → parent key (or same if empty)
}

// SaveInherit writes an InheritMap to disk as JSON.
func SaveInherit(path string, m InheritMap) error {
	if path == "" {
		return nil
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("inherit: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadInherit reads an InheritMap from disk. Returns empty map if file absent.
func LoadInherit(path string) (InheritMap, error) {
	m := InheritMap{Keys: make(map[string]string)}
	if path == "" {
		return m, nil
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return m, nil
	}
	if err != nil {
		return m, fmt.Errorf("inherit: read: %w", err)
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return m, fmt.Errorf("inherit: unmarshal: %w", err)
	}
	if m.Keys == nil {
		m.Keys = make(map[string]string)
	}
	return m, nil
}

// ApplyInherit fills missing keys in child from parent secrets according to the
// InheritMap. Parent key mapping: if Keys[childKey] is non-empty it is used as
// the parent key, otherwise childKey itself is looked up in parent.
// Existing child values are never overwritten.
func ApplyInherit(child, parent map[string]string, m InheritMap) map[string]string {
	out := make(map[string]string, len(child))
	for k, v := range child {
		out[k] = v
	}
	for childKey, parentKey := range m.Keys {
		if _, exists := out[childKey]; exists {
			continue
		}
		lookup := parentKey
		if lookup == "" {
			lookup = childKey
		}
		if val, ok := parent[lookup]; ok {
			out[childKey] = val
		}
	}
	return out
}
