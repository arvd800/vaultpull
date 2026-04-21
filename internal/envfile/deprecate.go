package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// DeprecationRecord holds metadata about a deprecated key.
type DeprecationRecord struct {
	Key        string    `json:"key"`
	Reason     string    `json:"reason"`
	ReplacedBy string    `json:"replaced_by,omitempty"`
	DeprecatedAt time.Time `json:"deprecated_at"`
}

// DeprecationMap maps key names to their deprecation records.
type DeprecationMap map[string]DeprecationRecord

// SaveDeprecations writes deprecation metadata to a JSON file.
func SaveDeprecations(path string, dm DeprecationMap) error {
	if path == "" {
		return nil
	}
	data, err := json.MarshalIndent(dm, "", "  ")
	if err != nil {
		return fmt.Errorf("deprecate: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadDeprecations reads deprecation metadata from a JSON file.
// Returns an empty map if the file does not exist.
func LoadDeprecations(path string) (DeprecationMap, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return DeprecationMap{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("deprecate: read: %w", err)
	}
	var dm DeprecationMap
	if err := json.Unmarshal(data, &dm); err != nil {
		return nil, fmt.Errorf("deprecate: unmarshal: %w", err)
	}
	return dm, nil
}

// MarkDeprecated adds or updates a deprecation entry for the given key.
func MarkDeprecated(dm DeprecationMap, key, reason, replacedBy string) DeprecationMap {
	out := make(DeprecationMap, len(dm)+1)
	for k, v := range dm {
		out[k] = v
	}
	out[key] = DeprecationRecord{
		Key:          key,
		Reason:       reason,
		ReplacedBy:   replacedBy,
		DeprecatedAt: time.Now().UTC(),
	}
	return out
}

// CheckDeprecations returns a list of warnings for any secrets keys that
// appear in the deprecation map.
func CheckDeprecations(secrets map[string]string, dm DeprecationMap) []string {
	var warnings []string
	for key := range secrets {
		if rec, ok := dm[key]; ok {
			msg := fmt.Sprintf("key %q is deprecated: %s", key, rec.Reason)
			if rec.ReplacedBy != "" {
				msg += fmt.Sprintf(" (use %q instead)", rec.ReplacedBy)
			}
			warnings = append(warnings, msg)
		}
	}
	return warnings
}
