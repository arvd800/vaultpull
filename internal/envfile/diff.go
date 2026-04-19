package envfile

// DiffResult holds the changes between two env maps.
type DiffResult struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string]string // key -> new value
}

// Diff compares an existing env map against an incoming one and returns
// what was added, removed, or changed.
func Diff(existing, incoming map[string]string) DiffResult {
	result := DiffResult{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string]string),
	}

	for k, newVal := range incoming {
		oldVal, ok := existing[k]
		if !ok {
			result.Added[k] = newVal
		} else if oldVal != newVal {
			result.Changed[k] = newVal
		}
	}

	for k := range existing {
		if _, ok := incoming[k]; !ok {
			result.Removed[k] = existing[k]
		}
	}

	return result
}

// HasChanges returns true if the DiffResult contains any differences.
func (d DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

// Summary returns a brief human-readable string describing the number of
// additions, removals, and changes in the DiffResult.
func (d DiffResult) Summary() string {
	if !d.HasChanges() {
		return "no changes"
	}
	return fmt.Sprintf("%d added, %d removed, %d changed",
		len(d.Added), len(d.Removed), len(d.Changed))
}
