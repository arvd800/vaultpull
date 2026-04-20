package envfile

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms a secret value.
type TransformFunc func(key, value string) (string, error)

// TransformRule defines a named transformation to apply to matching keys.
type TransformRule struct {
	KeyPrefix string
	Transform TransformFunc
}

// Transform applies a set of TransformRules to a secrets map.
// Rules are applied in order; the first matching prefix wins.
// Keys with no matching rule are passed through unchanged.
func Transform(secrets map[string]string, rules []TransformRule) (map[string]string, error) {
	if len(secrets) == 0 {
		return map[string]string{}, nil
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		transformed := false
		for _, rule := range rules {
			if rule.KeyPrefix == "" || strings.HasPrefix(k, rule.KeyPrefix) {
				newVal, err := rule.Transform(k, v)
				if err != nil {
					return nil, fmt.Errorf("transform rule %q failed on key %q: %w", rule.KeyPrefix, k, err)
				}
				out[k] = newVal
				transformed = true
				break
			}
		}
		if !transformed {
			out[k] = v
		}
	}
	return out, nil
}

// UpperCase is a built-in TransformFunc that upper-cases the value.
func UpperCase(_, value string) (string, error) {
	return strings.ToUpper(value), nil
}

// LowerCase is a built-in TransformFunc that lower-cases the value.
func LowerCase(_, value string) (string, error) {
	return strings.ToLower(value), nil
}

// TrimSpace is a built-in TransformFunc that trims whitespace from the value.
func TrimSpace(_, value string) (string, error) {
	return strings.TrimSpace(value), nil
}
