package sync

import (
	"fmt"

	"github.com/your-org/vaultpull/internal/envfile"
)

// GroupedSync filters the incoming secrets to only those belonging to the
// named group before merging and writing the env file. If groupsPath is empty
// or groupName is empty the function is a no-op pass-through.
func GroupedSync(secrets map[string]string, groupsPath, groupName string) (map[string]string, error) {
	if groupsPath == "" || groupName == "" {
		return secrets, nil
	}

	groups, err := envfile.LoadGroups(groupsPath)
	if err != nil {
		return nil, fmt.Errorf("grouper: load groups: %w", err)
	}

	filtered, err := envfile.ApplyGroup(secrets, groups, groupName)
	if err != nil {
		return nil, fmt.Errorf("grouper: apply group %q: %w", groupName, err)
	}

	return filtered, nil
}

// RegisterGroup adds or replaces a named group in the groups file at groupsPath.
func RegisterGroup(groupsPath, name string, keys []string) error {
	if groupsPath == "" {
		return fmt.Errorf("grouper: groups path is required")
	}
	if name == "" {
		return fmt.Errorf("grouper: group name is required")
	}

	groups, err := envfile.LoadGroups(groupsPath)
	if err != nil {
		return fmt.Errorf("grouper: load groups: %w", err)
	}

	updated := envfile.AddGroup(groups, name, keys)
	if err := envfile.SaveGroups(groupsPath, updated); err != nil {
		return fmt.Errorf("grouper: save groups: %w", err)
	}
	return nil
}
