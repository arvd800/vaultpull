package envfile

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// LockRecord stores metadata about a locked secret key.
type LockRecord struct {
	Key       string    `json:"key"`
	LockedAt  time.Time `json:"locked_at"`
	Reason    string    `json:"reason,omitempty"`
}

// LockFile maps keys to their lock records.
type LockFile map[string]LockRecord

// SaveLocks writes the lock file to disk.
func SaveLocks(path string, locks LockFile) error {
	if path == "" {
		return errors.New("lock file path must not be empty")
	}
	data, err := json.MarshalIndent(locks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadLocks reads the lock file from disk. Returns empty LockFile if not found.
func LoadLocks(path string) (LockFile, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return LockFile{}, nil
	}
	if err != nil {
		return nil, err
	}
	var locks LockFile
	if err := json.Unmarshal(data, &locks); err != nil {
		return nil, err
	}
	return locks, nil
}

// LockKey adds a lock entry for the given key.
func LockKey(locks LockFile, key, reason string) LockFile {
	out := make(LockFile, len(locks)+1)
	for k, v := range locks {
		out[k] = v
	}
	out[key] = LockRecord{Key: key, LockedAt: time.Now().UTC(), Reason: reason}
	return out
}

// UnlockKey removes a lock entry for the given key.
func UnlockKey(locks LockFile, key string) LockFile {
	out := make(LockFile, len(locks))
	for k, v := range locks {
		if k != key {
			out[k] = v
		}
	}
	return out
}

// IsLocked returns true if the key is currently locked.
func IsLocked(locks LockFile, key string) bool {
	_, ok := locks[key]
	return ok
}
