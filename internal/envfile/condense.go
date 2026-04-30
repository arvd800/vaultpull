package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// CondenseRule defines how to combine multiple keys into a single key.
type CondenseRule struct {
	OutputKey string   `json:"output_key"`
	SourceKeys []string `json:"source_keys"`
	Separator  string   `json:"separator"`
	DropSources bool    `json:"drop_sources"`
}

// CondenseConfig holds a list of condense rules.
type CondenseConfig struct {
	Rules []CondenseRule `json:"rules"`
}

// SaveCondenseConfig writes condense rules to a JSON file.
func SaveCondenseConfig(path string, cfg CondenseConfig) error {
	if path == "" {
		return nil
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("condense: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadCondenseConfig reads condense rules from a JSON file.
func LoadCondenseConfig(path string) (CondenseConfig, error) {
	if path == "" {
		return CondenseConfig{}, nil
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return CondenseConfig{}, nil
	}
	if err != nil {
		return CondenseConfig{}, fmt.Errorf("condense: read: %w", err)
	}
	var cfg CondenseConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return CondenseConfig{}, fmt.Errorf("condense: unmarshal: %w", err)
	}
	return cfg, nil
}

// Condense applies condense rules to a secrets map, joining source key values
// into a single output key. If DropSources is true, source keys are removed.
func Condense(secrets map[string]string, cfg CondenseConfig) (map[string]string, error) {
	if secrets == nil {
		return map[string]string{}, nil
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, rule := range cfg.Rules {
		if rule.OutputKey == "" || len(rule.SourceKeys) == 0 {
			continue
		}
		sep := rule.Separator
		parts := make([]string, 0, len(rule.SourceKeys))
		for _, src := range rule.SourceKeys {
			val, ok := out[src]
			if !ok {
				return nil, fmt.Errorf("condense: source key %q not found for rule %q", src, rule.OutputKey)
			}
			parts = append(parts, val)
		}
		combined := ""
		for i, p := range parts {
			if i > 0 {
				combined += sep
			}
			combined += p
		}
		out[rule.OutputKey] = combined
		if rule.DropSources {
			for _, src := range rule.SourceKeys {
				delete(out, src)
			}
		}
	}
	return out, nil
}

// ListCondenseOutputKeys returns the sorted list of output keys defined in cfg.
func ListCondenseOutputKeys(cfg CondenseConfig) []string {
	keys := make([]string, 0, len(cfg.Rules))
	for _, r := range cfg.Rules {
		if r.OutputKey != "" {
			keys = append(keys, r.OutputKey)
		}
	}
	sort.Strings(keys)
	return keys
}
