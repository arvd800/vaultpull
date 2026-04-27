package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTakeSnapshot_CopiesSecrets(t *testing.T) {
	orig := map[string]string{"FOO": "bar", "BAZ": "qux"}
	snap := TakeSnapshot("vault/secret/app", orig)

	if snap.Source != "vault/secret/app" {
		t.Errorf("unexpected source: %s", snap.Source)
	}
	if snap.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
	// mutate original — snapshot must not change
	orig["FOO"] = "mutated"
	if snap.Secrets["FOO"] != "bar" {
		t.Error("snapshot should not reflect mutation of source map")
	}
}

func TestSaveAndLoadSnapshot_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	origSnap := TakeSnapshot("vault/secret/db", secrets)

	if err := SaveSnapshot(path, origSnap); err != nil {
		t.Fatalf("SaveSnapshot error: %v", err)
	}

	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot error: %v", err)
	}
	if loaded.Source != origSnap.Source {
		t.Errorf("source mismatch: got %s want %s", loaded.Source, origSnap.Source)
	}
	if len(loaded.Secrets) != len(secrets) {
		t.Errorf("secrets length mismatch: got %d want %d", len(loaded.Secrets), len(secrets))
	}
	for k, v := range secrets {
		if loaded.Secrets[k] != v {
			t.Errorf("key %s: got %q want %q", k, loaded.Secrets[k], v)
		}
	}
}

func TestSaveAndLoadSnapshot_TimestampPreserved(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	origSnap := TakeSnapshot("vault/secret/app", map[string]string{"KEY": "val"})
	if err := SaveSnapshot(path, origSnap); err != nil {
		t.Fatalf("SaveSnapshot error: %v", err)
	}

	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot error: %v", err)
	}
	if !loaded.Timestamp.Equal(origSnap.Timestamp) {
		t.Errorf("timestamp mismatch: got %v want %v", loaded.Timestamp, origSnap.Timestamp)
	}
}

func TestSaveSnapshot_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	snap := TakeSnapshot("vault/secret/app", map[string]string{"KEY": "val"})
	if err := SaveSnapshot(path, snap); err != nil {
		t.Fatalf("SaveSnapshot error: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat error: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected perm 0600, got %o", perm)
	}
}

func TestLoadSnapshot_NonExistent(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/snap.json")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}
