package envfile

import "strings"

// FilterOptions controls which keys are included or excluded.
type FilterOptions struct {
	// IncludePrefix, if non-empty, only keeps keys with this prefix.
	IncludePrefix string
	// ExcludeKeys is a set of exact key names to drop.
	ExcludeKeys []string
	// StripPrefix removes the prefix from key names after filtering.
	StripPrefix bool
}

// Filter returns a new map containing only the entries that pass the filter
// options. The original map is never mutated.
func Filter(secrets map[string]string, opts FilterOptions) map[string]string {
	exclude := make(map[string]struct{}, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		exclude[k] = struct{}{}
	}

	result := make(map[string]string)
	for k, v := range secrets {
		if _, skip := exclude[k]; skip {
			continue
		}
		if opts.IncludePrefix != "" && !strings.HasPrefix(k, opts.IncludePrefix) {
			continue
		}
		outKey := k
		if opts.StripPrefix && opts.IncludePrefix != "" {
			outKey = strings.TrimPrefix(k, opts.IncludePrefix)
			if outKey == "" {
				continue
			}
		}
		result[outKey] = v
	}
	return result
}
