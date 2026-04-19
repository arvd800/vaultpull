package sync

import (
	"fmt"

	"github.com/yourusername/vaultpull/internal/envfile"
)

// ApplyLocks filters out locked keys from incoming secrets before writing,
// preserving their existing values from the current env file.
func ApplyLocks(
	incoming map[string]string,
	existing map[string]string,
	locks envfile.LockFile,
) (merged map[string]string, skipped []string) {
	merged = make(map[string]string, len(incoming))
	for k, v := range incoming {
		merged[k] = v
	}
	for key := range locks {
		if _, ok := merged[key]; !ok {
			continue
		}
		if prev, hasPrev := existing[key]; hasPrev {
			merged[key] = prev
		} else {
			delete(merged, key)
		}
		skipped = append(skipped, key)
	}
	return merged, skipped
}

// LogLockedKeys prints a notice for each skipped locked key.
func LogLockedKeys(skipped []string, locks envfile.LockFile) {
	for _, key := range skipped {
		rec := locks[key]
		msg := fmt.Sprintf("[lock] skipping %q (locked", key)
		if rec.Reason != "" {
			msg += fmt.Sprintf(", reason: %s", rec.Reason)
		}
		msg += ")"
		fmt.Println(msg)
	}
}
