package envfile

import (
	"fmt"
	"os"
	"time"
)

// Backup creates a timestamped backup of the given file if it exists.
// Returns the backup path, or empty string if the original did not exist.
func Backup(path string) (string, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("stat %s: %w", path, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", path, err)
	}

	ts := time.Now().UTC().Format("20060102T150405Z")
	backupPath := fmt.Sprintf("%s.%s.bak", path, ts)

	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return "", fmt.Errorf("write backup %s: %w", backupPath, err)
	}

	return backupPath, nil
}

// RemoveBackup deletes a backup file. Ignores not-found errors.
func RemoveBackup(path string) error {
	if path == "" {
		return nil
	}
	err := os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
