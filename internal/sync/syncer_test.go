package sync

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type mockVault struct {
	secrets map[string]string
	err     error
}

func (m *mockVault) ReadSecret(_ string) (map[string]string, error) {
	return m.secrets, m.err
}

func TestSyncer_Run(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	v := &mockVault{secrets: map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}}
	s := New(v, "secret/app", envPath, false)

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty env file")
	}
}

func TestSyncer_Run_WithBackup(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	if err := os.WriteFile(envPath, []byte("OLD=value\n"), 0600); err != nil {
		t.Fatal(err)
	}

	v := &mockVault{secrets: map[string]string{"NEW_KEY": "new_val"}}
	s := New(v, "secret/app", envPath, true)

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := os.ReadDir(dir)
	backupFound := false
	for _, e := range entries {
		if e.Name() != ".env" {
			backupFound = true
		}
	}
	if !backupFound {
		t.Error("expected backup file to exist")
	}
}

func TestSyncer_Run_BadPath(t *testing.T) {
	v := &mockVault{err: errors.New("forbidden")}
	s := New(v, "secret/missing", "/tmp/noop.env", false)

	if err := s.Run(); err == nil {
		t.Error("expected error for bad vault path")
	}
}
