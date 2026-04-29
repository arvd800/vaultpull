package envfile

import (
	"strings"
	"unicode"
)

// SanitizeOptions controls how sanitization is applied.
type SanitizeOptions struct {
	// TrimSpace removes leading/trailing whitespace from values.
	TrimSpace bool
	// RemoveNonPrintable strips non-printable characters from values.
	RemoveNonPrintable bool
	// NormalizeKeys uppercases all keys and replaces hyphens/spaces with underscores.
	NormalizeKeys bool
	// DropEmpty removes entries with empty values after sanitization.
	DropEmpty bool
}

// DefaultSanitizeOptions returns a sensible default configuration.
func DefaultSanitizeOptions() SanitizeOptions {
	return SanitizeOptions{
		TrimSpace:          true,
		RemoveNonPrintable: true,
		NormalizeKeys:      false,
		DropEmpty:          false,
	}
}

// Sanitize applies the given options to a copy of secrets and returns the
// cleaned map. The original map is never mutated.
func Sanitize(secrets map[string]string, opts SanitizeOptions) map[string]string {
	out := make(map[string]string, len(secrets))

	for k, v := range secrets {
		key := k
		val := v

		if opts.TrimSpace {
			val = strings.TrimSpace(val)
		}

		if opts.RemoveNonPrintable {
			val = removePrintable(val)
		}

		if opts.NormalizeKeys {
			key = normalizeKey(k)
		}

		if opts.DropEmpty && val == "" {
			continue
		}

		out[key] = val
	}

	return out
}

func removePrintable(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsPrint(r) || r == '\t' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func normalizeKey(k string) string {
	k = strings.ToUpper(k)
	k = strings.ReplaceAll(k, "-", "_")
	k = strings.ReplaceAll(k, " ", "_")
	return k
}
