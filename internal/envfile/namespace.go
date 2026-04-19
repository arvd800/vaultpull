package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Namespace represents a named grouping for secrets.
type Namespace struct {
	Name   string            `json:"name"`
	Prefix string            `json:"prefix"`
	Tags   map[string]string `json:"tags,omitempty"`
}

// NamespaceMap maps namespace names to their definitions.
type NamespaceMap map[string]Namespace

// SaveNamespaces writes namespace definitions to a JSON file.
func SaveNamespaces(path string, ns NamespaceMap) error {
	if path == "" {
		return nil
	}
	data, err := json.MarshalIndent(ns, "", "  ")
	if err != nil {
		return fmt.Errorf("namespace marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadNamespaces reads namespace definitions from a JSON file.
func LoadNamespaces(path string) (NamespaceMap, error) {
	if path == "" {
		return NamespaceMap{}, nil
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return NamespaceMap{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("namespace read: %w", err)
	}
	var ns NamespaceMap
	if err := json.Unmarshal(data, &ns); err != nil {
		return nil, fmt.Errorf("namespace unmarshal: %w", err)
	}
	return ns, nil
}

// ApplyNamespace filters and strips the namespace prefix from a secrets map.
func ApplyNamespace(secrets map[string]string, ns Namespace) map[string]string {
	out := make(map[string]string)
	for k, v := range secrets {
		if ns.Prefix == "" || strings.HasPrefix(k, ns.Prefix) {
			newKey := strings.TrimPrefix(k, ns.Prefix)
			if newKey == "" {
				continue
			}
			out[newKey] = v
		}
	}
	return out
}
