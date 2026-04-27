package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ExpiryRecord holds expiry metadata for a single secret key.
type ExpiryRecord struct {
	Key       string    `json:"key"`
	ExpiresAt time.Time `json:"expires_at"`
	Note      string    `json:"note,omitempty"`
}

// ExpiryMap maps secret keys to their expiry records.
type ExpiryMap map[string]ExpiryRecord

// SaveExpiry persists expiry records to the given path as JSON.
// A no-op if path is empty.
func SaveExpiry(path string, expiries ExpiryMap) error {
	if path == "" {
		return nil
	}
	data, err := json.MarshalIndent(expiries, "", "  ")
	if err != nil {
		return fmt.Errorf("expire: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadExpiry reads expiry records from the given path.
// Returns an empty map if the file does not exist.
func LoadExpiry(path string) (ExpiryMap, error) {
	if path == "" {
		return ExpiryMap{}, nil
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return ExpiryMap{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("expire: read: %w", err)
	}
	var m ExpiryMap
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("expire: unmarshal: %w", err)
	}
	return m, nil
}

// SetExpiry adds or updates an expiry record for the given key.
func SetExpiry(expiries ExpiryMap, key string, ttl time.Duration, note string) ExpiryMap {
	out := make(ExpiryMap, len(expiries)+1)
	for k, v := range expiries {
		out[k] = v
	}
	out[key] = ExpiryRecord{
		Key:       key,
		ExpiresAt: time.Now().UTC().Add(ttl),
		Note:      note,
	}
	return out
}

// CheckExpiry returns the keys that have expired relative to now.
func CheckExpiry(secrets map[string]string, expiries ExpiryMap) []string {
	now := time.Now().UTC()
	var expired []string
	for key := range secrets {
		rec, ok := expiries[key]
		if !ok {
			continue
		}
		if now.After(rec.ExpiresAt) {
			expired = append(expired, key)
		}
	}
	return expired
}
