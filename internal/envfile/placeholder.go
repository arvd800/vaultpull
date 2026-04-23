package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// PlaceholderConfig controls how placeholders are detected and handled.
type PlaceholderConfig struct {
	// Prefix is the sentinel prefix that marks a value as a placeholder.
	// Defaults to "PLACEHOLDER:" if empty.
	Prefix string
	// FailOnUnresolved causes ResolvePlaceholders to return an error if any
	// placeholder remains after substitution.
	FailOnUnresolved bool
}

var identRe = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// ResolvePlaceholders replaces placeholder values in dst with real values
// sourced from src. A placeholder is a value that starts with cfg.Prefix
// followed by a key name that must exist in src.
//
// Example: if cfg.Prefix is "PLACEHOLDER:" and dst contains
// FOO="PLACEHOLDER:BAR", the function looks up "BAR" in src and sets
// FOO to that value.
func ResolvePlaceholders(dst, src map[string]string, cfg PlaceholderConfig) (map[string]string, error) {
	prefix := cfg.Prefix
	if prefix == "" {
		prefix = "PLACEHOLDER:"
	}

	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	for k, v := range out {
		if !strings.HasPrefix(v, prefix) {
			continue
		}
		ref := strings.TrimPrefix(v, prefix)
		ref = strings.TrimSpace(ref)
		if replacement, ok := src[ref]; ok {
			out[k] = replacement
		} else if cfg.FailOnUnresolved {
			return nil, fmt.Errorf("placeholder %q in key %q could not be resolved: key %q not found in source", v, k, ref)
		}
	}
	return out, nil
}

// ListPlaceholders returns the keys in m whose values are placeholders
// according to the given prefix.
func ListPlaceholders(m map[string]string, prefix string) []string {
	if prefix == "" {
		prefix = "PLACEHOLDER:"
	}
	var out []string
	for k, v := range m {
		if strings.HasPrefix(v, prefix) {
			out = append(out, k)
		}
	}
	return sorted(out)
}
