package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// TTLRecord stores metadata about when secrets were last synced and when they expire.
type TTLRecord struct {
	SyncedAt  time.Time     `json:"synced_at"`
	ExpiresAt time.Time     `json:"expires_at"`
	TTL       time.Duration `json:"ttl"`
}

// Expired returns true if the TTL has passed.
func (r TTLRecord) Expired() bool {
	return time.Now().After(r.ExpiresAt)
}

// SaveTTL writes a TTLRecord to the given path as JSON.
func SaveTTL(path string, ttl time.Duration) error {
	now := time.Now().UTC()
	rec := TTLRecord{
		SyncedAt:  now,
		ExpiresAt: now.Add(ttl),
		TTL:       ttl,
	}
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal ttl: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadTTL reads a TTLRecord from the given path.
func LoadTTL(path string) (TTLRecord, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return TTLRecord{}, fmt.Errorf("read ttl file: %w", err)
	}
	var rec TTLRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return TTLRecord{}, fmt.Errorf("unmarshal ttl: %w", err)
	}
	return rec, nil
}

// CheckTTL returns an error if the secrets at path are expired.
func CheckTTL(path string) error {
	rec, err := LoadTTL(path)
	if err != nil {
		return err
	}
	if rec.Expired() {
		return fmt.Errorf("secrets expired at %s (TTL: %s)", rec.ExpiresAt.Format(time.RFC3339), rec.TTL)
	}
	return nil
}
