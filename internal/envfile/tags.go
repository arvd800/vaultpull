package envfile

import (
	"encoding/json"
	"os"
	"time"
)

// TagRecord holds metadata tags for a secret map.
type TagRecord struct {
	Source    string            `json:"source"`
	FetchedAt time.Time         `json:"fetched_at"`
	Tags      map[string]string `json:"tags"`
}

// SaveTags writes a TagRecord to a JSON file.
func SaveTags(path string, record TagRecord) error {
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadTags reads a TagRecord from a JSON file.
func LoadTags(path string) (TagRecord, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return TagRecord{Tags: map[string]string{}}, nil
		}
		return TagRecord{}, err
	}
	var record TagRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return TagRecord{}, err
	}
	if record.Tags == nil {
		record.Tags = map[string]string{}
	}
	return record, nil
}

// MergeTags merges additional tags into an existing TagRecord.
func MergeTags(record *TagRecord, extra map[string]string) {
	for k, v := range extra {
		record.Tags[k] = v
	}
}
