package envfile

// Merge combines existing env vars with incoming ones from Vault.
// Incoming keys are added or updated; keys present only in existing are preserved.
// Neither input map is mutated.
func Merge(existing, incoming map[string]string) map[string]string {
	result := make(map[string]string, len(existing)+len(incoming))

	for k, v := range existing {
		result[k] = v
	}
	for k, v := range incoming {
		result[k] = v
	}

	return result
}
