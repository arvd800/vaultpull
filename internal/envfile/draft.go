package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Draft represents an uncommitted set of secret changes staged for review.
type Draft struct {
	ID        string            `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	Secrets   map[string]string `json:"secrets"`
	Message   string            `json:"message,omitempty"`
}

// SaveDraft persists a Draft to disk as JSON.
func SaveDraft(path string, d Draft) error {
	if path == "" {
		return nil
	}
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return fmt.Errorf("draft: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadDraft reads a Draft from disk. Returns an empty Draft if the file does
// not exist.
func LoadDraft(path string) (Draft, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Draft{}, nil
	}
	if err != nil {
		return Draft{}, fmt.Errorf("draft: read: %w", err)
	}
	var d Draft
	if err := json.Unmarshal(data, &d); err != nil {
		return Draft{}, fmt.Errorf("draft: unmarshal: %w", err)
	}
	return d, nil
}

// NewDraft creates a new Draft with a unique ID derived from the current time.
func NewDraft(secrets map[string]string, message string) Draft {
	now := time.Now().UTC()
	return Draft{
		ID:        fmt.Sprintf("draft-%d", now.UnixNano()),
		CreatedAt: now,
		Secrets:   copyMap(secrets),
		Message:   message,
	}
}

// DiscardDraft removes the draft file from disk.
func DiscardDraft(path string) error {
	if path == "" {
		return nil
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("draft: discard: %w", err)
	}
	return nil
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
