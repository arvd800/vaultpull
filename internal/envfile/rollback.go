package envfile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RollbackEntry represents a single rollback point derived from history.
type RollbackEntry struct {
	Index     int
	Timestamp string
	Secrets   map[string]string
}

// ListRollbackPoints returns all history entries available for rollback.
func ListRollbackPoints(historyPath string) ([]RollbackEntry, error) {
	entries, err := LoadHistory(historyPath)
	if err != nil {
		return nil, fmt.Errorf("rollback: load history: %w", err)
	}
	points := make([]RollbackEntry, len(entries))
	for i, e := range entries {
		points[i] = RollbackEntry{
			Index:     i,
			Timestamp: e.Timestamp,
			Secrets:   e.Secrets,
		}
	}
	return points, nil
}

// Rollback restores secrets from the given history index to the target .env file.
// It writes a backup of the current file before overwriting.
func Rollback(historyPath, targetPath string, index int) error {
	points, err := ListRollbackPoints(historyPath)
	if err != nil {
		return err
	}
	if index < 0 || index >= len(points) {
		return fmt.Errorf("rollback: index %d out of range (0-%d)", index, len(points)-1)
	}
	_, err = Backup(targetPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("rollback: backup current file: %w", err)
	}
	secrets := points[index].Secrets
	if err := Write(targetPath, secrets); err != nil {
		return fmt.Errorf("rollback: write secrets: %w", err)
	}
	return nil
}

// FormatRollbackList returns a human-readable summary of rollback points.
func FormatRollbackList(points []RollbackEntry) string {
	if len(points) == 0 {
		return "no rollback points available"
	}
	var sb strings.Builder
	for _, p := range points {
		keys := make([]string, 0, len(p.Secrets))
		for k := range p.Secrets {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		fmt.Fprintf(&sb, "[%d] %s (%d keys: %s)\n",
			p.Index, p.Timestamp, len(keys), strings.Join(keys, ", "))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// RollbackDir returns the default directory for rollback/history files.
func RollbackDir(base string) string {
	return filepath.Join(filepath.Dir(base), ".vaultpull")
}
