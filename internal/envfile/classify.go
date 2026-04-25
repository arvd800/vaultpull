package envfile

import (
	"encoding/json"
	"os"
	"sort"
)

// Classification holds a label and the set of keys assigned to it.
type Classification struct {
	Label string   `json:"label"`
	Keys  []string `json:"keys"`
}

// ClassificationMap maps label names to their Classification entries.
type ClassificationMap map[string]Classification

// SaveClassifications writes the classification map to the given path as JSON.
func SaveClassifications(path string, cm ClassificationMap) error {
	if path == "" {
		return nil
	}
	b, err := json.MarshalIndent(cm, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0600)
}

// LoadClassifications reads a classification map from the given path.
// Returns an empty map if the file does not exist.
func LoadClassifications(path string) (ClassificationMap, error) {
	cm := make(ClassificationMap)
	if path == "" {
		return cm, nil
	}
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cm, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &cm); err != nil {
		return nil, err
	}
	return cm, nil
}

// Classify returns a new map containing only the keys whose values belong to
// the requested label according to the classification map. If the label is not
// found the returned map is empty.
func Classify(secrets map[string]string, cm ClassificationMap, label string) map[string]string {
	out := make(map[string]string)
	entry, ok := cm[label]
	if !ok {
		return out
	}
	for _, k := range entry.Keys {
		if v, exists := secrets[k]; exists {
			out[k] = v
		}
	}
	return out
}

// ListLabels returns all label names in the classification map, sorted.
func ListLabels(cm ClassificationMap) []string {
	labels := make([]string, 0, len(cm))
	for l := range cm {
		labels = append(labels, l)
	}
	sort.Strings(labels)
	return labels
}
