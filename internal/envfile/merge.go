package envfile

// Merge combines existing and incoming secret maps.
// Incoming values take precedence over existing ones.
// Keys present only in existing are preserved.
// The original maps are not mutated.
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
