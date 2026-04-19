package envfile

import (
	"fmt"
	"maps"
)

// PromoteOptions controls how secrets are promoted between environments.
type PromoteOptions struct {
	// Keys to promote; if empty, all keys are promoted.
	Keys []string
	// SkipExisting skips keys already present in the destination.
	SkipExisting bool
	// DryRun returns the result without writing anything.
	DryRun bool
}

// PromoteResult describes what changed during a promotion.
type PromoteResult struct {
	Promoted []string
	Skipped  []string
}

// Promote copies secrets from src into dst according to opts.
// It returns the merged map and a result describing what happened.
func Promote(src, dst map[string]string, opts PromoteOptions) (map[string]string, PromoteResult, error) {
	if src == nil {
		return nil, PromoteResult{}, fmt.Errorf("promote: src must not be nil")
	}

	keys := opts.Keys
	if len(keys) == 0 {
		for k := range src {
			keys = append(keys, k)
		}
	}

	out := make(map[string]string)
	if dst != nil {
		maps.Copy(out, dst)
	}

	var result PromoteResult
	for _, k := range keys {
		v, ok := src[k]
		if !ok {
			continue
		}
		if opts.SkipExisting {
			if _, exists := out[k]; exists {
				result.Skipped = append(result.Skipped, k)
				continue
			}
		}
		out[k] = v
		result.Promoted = append(result.Promoted, k)
	}

	return out, result, nil
}
