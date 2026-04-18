package envfile

import "strings"

// RedactOptions controls which keys are redacted.
type RedactOptions struct {
	// SensitiveSubstrings are case-insensitive substrings; any key containing
	// one will have its value replaced with the mask.
	SensitiveSubstrings []string
	// Mask is the replacement string. Defaults to "***".
	Mask string
}

var defaultSensitive = []string{"password", "secret", "token", "key", "apikey", "api_key", "passwd", "credential"}

// Redact returns a copy of secrets with sensitive values masked.
func Redact(secrets map[string]string, opts *RedactOptions) map[string]string {
	substrs := defaultSensitive
	mask := "***"
	if opts != nil {
		if len(opts.SensitiveSubstrings) > 0 {
			substrs = opts.SensitiveSubstrings
		}
		if opts.Mask != "" {
			mask = opts.Mask
		}
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if isSensitive(k, substrs) {
			out[k] = mask
		} else {
			out[k] = v
		}
	}
	return out
}

func isSensitive(key string, substrs []string) bool {
	lower := strings.ToLower(key)
	for _, s := range substrs {
		if strings.Contains(lower, strings.ToLower(s)) {
			return true
		}
	}
	return false
}
