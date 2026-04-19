package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// HistoryEntry records a single sync event.
type HistoryEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	Added     []string          `json:"added"`
	Removed   []string          `json:"removed"`
	Changed   []string          `json:"changed"`
	Snapshot  map[string]string `json:"snapshot"`
}

// HistoryLog holds multiple entries.
type HistoryLog struct {
	Entries []HistoryEntry `json:"entries"`
}

// AppendHistory loads an existing history file, appends the entry, and saves it.
func AppendHistory(path string, entry HistoryEntry) error {
	if path == "" {
		return nil
	}
	log, _ := LoadHistory(path)
	log.Entries = append(log.Entries, entry)
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return fmt.Errorf("history marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadHistory reads a history log from disk.
func LoadHistory(path string) (HistoryLog, error) {
	var log HistoryLog
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return log, nil
		}
		return log, fmt.Errorf("history read: %w", err)
	}
	if err := json.Unmarshal(data, &log); err != nil {
		return log, fmt.Errorf("history unmarshal: %w", err)
	}
	return log, nil
}

// BuildHistoryEntry constructs a HistoryEntry from a Diff and current secrets.
func BuildHistoryEntry(source string, d DiffResult, current map[string]string) HistoryEntry {
	snap := make(map[string]string, len(current))
	for k, v := range current {
		snap[k] = v
	}
	return HistoryEntry{
		Timestamp: time.Now().UTC(),
		Source:    source,
		Added:     keys(d.Added),
		Removed:   keys(d.Removed),
		Changed:   keys(d.Changed),
		Snapshot:  snap,
	}
}

func keys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
