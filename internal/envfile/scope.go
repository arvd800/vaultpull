package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// Scope represents a named set of allowed keys for a given environment scope.
type Scope struct {
	Name string   `json:"name"`
	Keys []string `json:"keys"`
}

// ScopeMap maps scope names to their key sets.
type ScopeMap map[string]Scope

// SaveScopes writes the scope map to a JSON file.
func SaveScopes(path string, scopes ScopeMap) error {
	if path == "" {
		return nil
	}
	data, err := json.MarshalIndent(scopes, "", "  ")
	if err != nil {
		return fmt.Errorf("scope: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadScopes reads the scope map from a JSON file.
// Returns an empty map if the file does not exist.
func LoadScopes(path string) (ScopeMap, error) {
	scopes := make(ScopeMap)
	if path == "" {
		return scopes, nil
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return scopes, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scope: read: %w", err)
	}
	if err := json.Unmarshal(data, &scopes); err != nil {
		return nil, fmt.Errorf("scope: unmarshal: %w", err)
	}
	return scopes, nil
}

// ApplyScope filters secrets to only those keys listed in the named scope.
// If the scope name is not found, an error is returned.
func ApplyScope(secrets map[string]string, scopes ScopeMap, name string) (map[string]string, error) {
	scope, ok := scopes[name]
	if !ok {
		return nil, fmt.Errorf("scope: %q not found", name)
	}
	result := make(map[string]string, len(scope.Keys))
	for _, key := range scope.Keys {
		if val, exists := secrets[key]; exists {
			result[key] = val
		}
	}
	return result, nil
}

// ListScopes returns sorted scope names from the map.
func ListScopes(scopes ScopeMap) []string {
	names := make([]string, 0, len(scopes))
	for name := range scopes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
