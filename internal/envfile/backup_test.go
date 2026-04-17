package envfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBackup_CreatesBackupFile(t *testing.T) {
	dir := t.TempDir()
	original := filepath.Join(dir, ".env")

	if err := os.WriteFile(original, []byte("FOO=bar\n"), 0600); err != nil {
		t.Fatal(err)
	}

	backupPath, err := Backup(original)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backupPath == "" {
		t.Fatal("expected non-empty backup path")
	}
	if !strings.HasSuffix(backupPath, ".bak") {
		t.Errorf("backup path should end with .bak, got %s", backupPath)
	}

	data, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("cannot read backup: %v", err)
	}
	if string(data) != "FOO=bar\n" {
		t.Errorf("backup content mismatch: %q", data)
	}
}

func TestBackup_NonExistentFile(t *testing.T) {
	dir := t.TempDir()
	path, err := Backup(filepath.Join(dir, ".env"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "" {
		t.Errorf("expected empty path for non-existent file, got %s", path)
	}
}

func TestRemoveBackup_RemovesFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "backup.bak")
	if err := os.WriteFile(f, []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := RemoveBackup(f); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		t.Error("expected file to be removed")
	}
}

func TestRemoveBackup_EmptyPath(t *testing.T) {
	if err := RemoveBackup(""); err != nil {
		t.Fatalf("unexpected error on empty path: %v", err)
	}
}
