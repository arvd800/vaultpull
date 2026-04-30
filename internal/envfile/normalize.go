package envfile

import (
	"strings"
)

// NormalizeOptions controls how keys and values are normalized.
type NormalizeOptions struct {
	// UpperKeys converts all keys to UPPER_CASE.
	UpperKeys bool
	// ReplaceHyphens replaces hyphens in keys with underscores.
	ReplaceHyphens bool
	// TrimValues trims leading and trailing whitespace from values.
	TrimValues bool
	// CollapseUnderscores collapses consecutive underscores into one.
	CollapseUnderscores bool
}

// DefaultNormalizeOptions returns an opinionated default normalization config.
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		UpperKeys:           true,
		ReplaceHyphens:      true,
		TrimValues:          true,
		CollapseUnderscores: false,
	}
}

// Normalize applies the given NormalizeOptions to a secrets map and returns
// a new map with normalized keys and/or values. The original map is not mutated.
func Normalize(secrets map[string]string, opts NormalizeOptions) (map[string]string, error) {
	if secrets == nil {
		return map[string]string{}, nil
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		nk := normalizeKeyN(k, opts)
		nv := normalizeValue(v, opts)
		out[nk] = nv
	}
	return out, nil
}

func normalizeKeyN(key string, opts NormalizeOptions) string {
	if opts.ReplaceHyphens {
		key = strings.ReplaceAll(key, "-", "_")
	}
	if opts.CollapseUnderscores {
		for strings.Contains(key, "__") {
			key = strings.ReplaceAll(key, "__", "_")
		}
	}
	if opts.UpperKeys {
		key = strings.ToUpper(key)
	}
	return key
}

func normalizeValue(val string, opts NormalizeOptions) string {
	if opts.TrimValues {
		val = strings.TrimSpace(val)
	}
	return val
}
