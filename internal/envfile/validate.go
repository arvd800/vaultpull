package envfile

import (
	"fmt"
	"regexp"
)

var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ValidationError holds all invalid keys found during validation.
type ValidationError struct {
	InvalidKeys []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid env keys: %v", e.InvalidKeys)
}

// Validate checks that all keys in the provided map are valid environment
// variable names (letters, digits, underscores; must not start with a digit).
// Returns a *ValidationError if any keys are invalid, nil otherwise.
func Validate(secrets map[string]string) error {
	var invalid []string
	for k := range secrets {
		if !validKeyRe.MatchString(k) {
			invalid = append(invalid, k)
		}
	}
	if len(invalid) > 0 {
		return &ValidationError{InvalidKeys: invalid}
	}
	return nil
}
