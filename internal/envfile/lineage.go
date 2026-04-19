package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// LineageRecord tracks the origin of a secret sync.
type LineageRecord struct {
	SyncedAt   time.Time         `json:"synced_at"`
	VaultAddr  string            `json:"vault_addr"`
	SecretPath string            `json:"secret_path"`
	Keys       []string          `json:"keys"`
	Meta       map[string]string `json:"meta,omitempty"`
}

// SaveLineage writes a lineage record to path as JSON.
func SaveLineage(path string, record LineageRecord) error {
	if path == "" {
		return fmt.Errorf("lineage path must not be empty")
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal lineage: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadLineage reads a lineage record from path.
func LoadLineage(path string) (LineageRecord, error) {
	var record LineageRecord
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return record, nil
		}
		return record, fmt.Errorf("read lineage: %w", err)
	}
	if err := json.Unmarshal(data, &record); err != nil {
		return record, fmt.Errorf("unmarshal lineage: %w", err)
	}
	return record, nil
}

// BuildLineage constructs a LineageRecord from the given secrets map and metadata.
func BuildLineage(vaultAddr, secretPath string, secrets map[string]string, meta map[string]string) LineageRecord {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	return LineageRecord{
		SyncedAt:   time.Now().UTC(),
		VaultAddr:  vaultAddr,
		SecretPath: secretPath,
		Keys:       keys,
		Meta:       meta,
	}
}
