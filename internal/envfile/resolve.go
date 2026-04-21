package envfile

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ResolveOptions controls how variable interpolation is performed.
type ResolveOptions struct {
	// AllowMissing skips unresolved references instead of returning an error.
	AllowMissing bool
	// FallbackToEnv checks the process environment for missing keys.
	FallbackToEnv bool
}

var interpolationRe = regexp.MustCompile(`\$\{([^}]+)\}`)

// Resolve performs variable interpolation on values within the secrets map.
// References of the form ${KEY} are replaced with the value of KEY from the
// same map (or from the OS environment when FallbackToEnv is set).
func Resolve(secrets map[string]string, opts ResolveOptions) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		resolved, err := interpolate(v, secrets, opts)
		if err != nil {
			return nil, fmt.Errorf("resolve %q: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}

func interpolate(value string, secrets map[string]string, opts ResolveOptions) (string, error) {
	var resolveErr error
	result := interpolationRe.ReplaceAllStringFunc(value, func(match string) string {
		if resolveErr != nil {
			return match
		}
		key := strings.TrimSpace(match[2 : len(match)-1])
		if v, ok := secrets[key]; ok {
			return v
		}
		if opts.FallbackToEnv {
			if v, ok := os.LookupEnv(key); ok {
				return v
			}
		}
		if opts.AllowMissing {
			return match
		}
		resolveErr = fmt.Errorf("undefined variable %q", key)
		return match
	})
	if resolveErr != nil {
		return "", resolveErr
	}
	return result, nil
}
