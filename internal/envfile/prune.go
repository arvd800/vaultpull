package envfile

import (
	"fmt"
	"sort"
)

// PruneResult holds the outcome of a prune operation.
type PruneResult struct {
	Removed []string
	Retained map[string]string
}

// PruneOptions controls how pruning behaves.
type PruneOptions struct {
	// DryRun reports what would be removed without modifying anything.
	DryRun bool
	// KeepKeys is an explicit allowlist; if non-empty, any key not in this
	// set is considered a candidate for removal.
	KeepKeys []string
	// RemoveKeys is an explicit denylist of keys to remove.
	RemoveKeys []string
}

// Prune removes keys from secrets according to opts and returns a PruneResult.
// If both KeepKeys and RemoveKeys are empty the input is returned unchanged.
func Prune(secrets map[string]string, opts PruneOptions) (PruneResult, error) {
	if secrets == nil {
		return PruneResult{Retained: map[string]string{}}, nil
	}

	keepSet := make(map[string]bool, len(opts.KeepKeys))
	for _, k := range opts.KeepKeys {
		keepSet[k] = true
	}

	removeSet := make(map[string]bool, len(opts.RemoveKeys))
	for _, k := range opts.RemoveKeys {
		removeSet[k] = true
	}

	if len(keepSet) > 0 && len(removeSet) > 0 {
		return PruneResult{}, fmt.Errorf("prune: KeepKeys and RemoveKeys are mutually exclusive")
	}

	retained := make(map[string]string, len(secrets))
	var removed []string

	for k, v := range secrets {
		shouldRemove := false
		if len(keepSet) > 0 && !keepSet[k] {
			shouldRemove = true
		}
		if removeSet[k] {
			shouldRemove = true
		}
		if shouldRemove {
			removed = append(removed, k)
		} else {
			retained[k] = v
		}
	}

	sort.Strings(removed)

	if opts.DryRun {
		// Return original secrets as retained, only report what would be removed.
		retainedCopy := make(map[string]string, len(secrets))
		for k, v := range secrets {
			retainedCopy[k] = v
		}
		return PruneResult{Removed: removed, Retained: retainedCopy}, nil
	}

	return PruneResult{Removed: removed, Retained: retained}, nil
}
