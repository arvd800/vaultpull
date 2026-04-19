package envfile

import (
	"encoding/json"
	"errors"
	"os"
)

// ImmutableRecord stores keys that cannot be overwritten during sync.
type ImmutableRecord struct {
	Keys map[string]bool `json:"keys"`
}

// SaveImmutable persists the immutable key set to path.
func SaveImmutable(path string, rec ImmutableRecord) error {
	if path == "" {
		return nil
	}
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadImmutable reads the immutable key set from path.
// Returns an empty record if the file does not exist.
func LoadImmutable(path string) (ImmutableRecord, error) {
	rec := ImmutableRecord{Keys: map[string]bool{}}
	if path == "" {
		return rec, nil
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return rec, nil
	}
	if err != nil {
		return rec, err
	}
	if err := json.Unmarshal(data, &rec); err != nil {
		return rec, err
	}
	if rec.Keys == nil {
		rec.Keys = map[string]bool{}
	}
	return rec, nil
}

// MarkImmutable adds key to the record.
func MarkImmutable(rec ImmutableRecord, key string) ImmutableRecord {
	out := ImmutableRecord{Keys: map[string]bool{}}
	for k, v := range rec.Keys {
		out.Keys[k] = v
	}
	out.Keys[key] = true
	return out
}

// ApplyImmutable returns incoming with any immutable keys restored from existing.
func ApplyImmutable(existing, incoming map[string]string, rec ImmutableRecord) map[string]string {
	out := map[string]string{}
	for k, v := range incoming {
		out[k] = v
	}
	for key := range rec.Keys {
		if val, ok := existing[key]; ok {
			out[key] = val
		} else {
			delete(out, key)
		}
	}
	return out
}
