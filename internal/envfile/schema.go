package envfile

import (
	"fmt"
	"regexp"
)

// SchemaRule defines a validation rule for a specific env key.
type SchemaRule struct {
	Required bool
	Pattern  string // optional regex pattern for value
}

// Schema maps env key names to their rules.
type Schema map[string]SchemaRule

// ValidateSchema checks a secrets map against a schema.
// It returns a list of violations found.
func ValidateSchema(secrets map[string]string, schema Schema) []error {
	var errs []error

	for key, rule := range schema {
		val, exists := secrets[key]
		if rule.Required && !exists {
			errs = append(errs, fmt.Errorf("required key %q is missing", key))
			continue
		}
		if exists && rule.Pattern != "" {
			matched, err := regexp.MatchString(rule.Pattern, val)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid pattern for key %q: %w", key, err))
				continue
			}
			if !matched {
				errs = append(errs, fmt.Errorf("key %q value does not match pattern %q", key, rule.Pattern))
			}
		}
	}

	return errs
}
