package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// GroupMap maps group names to sets of keys.
type GroupMap map[string][]string

// SaveGroups persists a GroupMap to a JSON file at path.
func SaveGroups(path string, groups GroupMap) error {
	if path == "" {
		return nil
	}
	data, err := json.MarshalIndent(groups, "", "  ")
	if err != nil {
		return fmt.Errorf("group: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadGroups reads a GroupMap from a JSON file at path.
// Returns an empty GroupMap if the file does not exist.
func LoadGroups(path string) (GroupMap, error) {
	groups := make(GroupMap)
	if path == "" {
		return groups, nil
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return groups, nil
	}
	if err != nil {
		return nil, fmt.Errorf("group: read: %w", err)
	}
	if err := json.Unmarshal(data, &groups); err != nil {
		return nil, fmt.Errorf("group: unmarshal: %w", err)
	}
	return groups, nil
}

// AddGroup adds or replaces a named group with the given keys.
func AddGroup(groups GroupMap, name string, keys []string) GroupMap {
	out := make(GroupMap, len(groups)+1)
	for k, v := range groups {
		out[k] = v
	}
	copy := make([]string, len(keys))
	for i, k := range keys {
		copy[i] = k
	}
	out[name] = copy
	return out
}

// ApplyGroup filters secrets to only those keys belonging to the named group.
// Returns an error if the group does not exist.
func ApplyGroup(secrets map[string]string, groups GroupMap, name string) (map[string]string, error) {
	keys, ok := groups[name]
	if !ok {
		available := make([]string, 0, len(groups))
		for k := range groups {
			available = append(available, k)
		}
		sort.Strings(available)
		return nil, fmt.Errorf("group %q not found; available: %v", name, available)
	}
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		if v, exists := secrets[k]; exists {
			out[k] = v
		}
	}
	return out, nil
}
