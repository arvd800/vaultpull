package envfile

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// ChecksumRecord holds a SHA-256 digest of a secrets map along with metadata.
type ChecksumRecord struct {
	Digest    string    `json:"digest"`
	KeyCount  int       `json:"key_count"`
	CreatedAt time.Time `json:"created_at"`
}

// ComputeChecksum returns a deterministic SHA-256 hex digest of the given
// secrets map. Keys are sorted before hashing to ensure stability.
func ComputeChecksum(secrets map[string]string) string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s\n", k, secrets[k])
	}
	return hex.EncodeToString(h.Sum(nil))
}

// SaveChecksum writes a ChecksumRecord for the given secrets map to path.
func SaveChecksum(path string, secrets map[string]string) error {
	if path == "" {
		return nil
	}
	rec := ChecksumRecord{
		Digest:    ComputeChecksum(secrets),
		KeyCount:  len(secrets),
		CreatedAt: time.Now().UTC(),
	}
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("checksum: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("checksum: write %s: %w", path, err)
	}
	return nil
}

// LoadChecksum reads a ChecksumRecord from path.
// Returns a zero-value record and no error when the file does not exist.
func LoadChecksum(path string) (ChecksumRecord, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return ChecksumRecord{}, nil
	}
	if err != nil {
		return ChecksumRecord{}, fmt.Errorf("checksum: read %s: %w", path, err)
	}
	var rec ChecksumRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return ChecksumRecord{}, fmt.Errorf("checksum: unmarshal: %w", err)
	}
	return rec, nil
}

// VerifyChecksum returns true when the digest stored at path matches the
// computed digest of secrets. It returns false (not an error) when no
// checksum file exists yet.
func VerifyChecksum(path string, secrets map[string]string) (bool, error) {
	rec, err := LoadChecksum(path)
	if err != nil {
		return false, err
	}
	if rec.Digest == "" {
		return false, nil
	}
	return rec.Digest == ComputeChecksum(secrets), nil
}
