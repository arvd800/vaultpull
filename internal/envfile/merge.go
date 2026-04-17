package envfile

// MergeStrategy controls how existing keys are handled during a merge.
type MergeStrategy int

const (
	// Overwrite replaces existing keys with values from incoming.
	Overwrite MergeStrategy = iota
	// Preserve keeps existing keys and only adds new ones.
	Preserve
)

// Merge combines existing env vars with incoming ones from Vault.
// The strategy determines behaviour when a key exists in both maps.
func Merge(existing, incoming map[string]string, strategy MergeStrategy) map[string]string {
	result := make(map[string]string, len(existing))
	for k, v := range existing {
		result[k] = v
	}
	for k, v := range incoming {
		if _, exists := result[k]; exists && strategy == Preserve {
			continue
		}
		result[k] = v
	}
	return result
}
