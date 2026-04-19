package envfile

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveAndLoadLineage_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lineage.json")

	record := LineageRecord{
		SyncedAt:   time.Now().UTC().Truncate(time.Second),
		VaultAddr:  "https://vault.example.com",
		SecretPath: "secret/myapp",
		Keys:       []string{"DB_HOST", "DB_PASS"},
		Meta:       map[string]string{"env": "production"},
	}

	if err := SaveLineage(path, record); err != nil {
		t.Fatalf("SaveLineage: %v", err)
	}

	loaded, err := LoadLineage(path)
	if err != nil {
		t.Fatalf("LoadLineage: %v", err)
	}

	if loaded.VaultAddr != record.VaultAddr {
		t.Errorf("VaultAddr mismatch: got %q want %q", loaded.VaultAddr, record.VaultAddr)
	}
	if loaded.SecretPath != record.SecretPath {
		t.Errorf("SecretPath mismatch: got %q want %q", loaded.SecretPath, record.SecretPath)
	}
	if len(loaded.Keys) != len(record.Keys) {
		t.Errorf("Keys length mismatch: got %d want %d", len(loaded.Keys), len(record.Keys))
	}
	if loaded.Meta["env"] != "production" {
		t.Errorf("Meta mismatch: got %v", loaded.Meta)
	}
}

func TestLoadLineage_NonExistent(t *testing.T) {
	record, err := LoadLineage("/nonexistent/lineage.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if record.VaultAddr != "" {
		t.Errorf("expected empty record, got %+v", record)
	}
}

func TestSaveLineage_EmptyPath(t *testing.T) {
	err := SaveLineage("", LineageRecord{})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestSaveLineage_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lineage.json")

	if err := SaveLineage(path, LineageRecord{VaultAddr: "https://vault"}); err != nil {
		t.Fatalf("SaveLineage: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected perm 0600, got %v", info.Mode().Perm())
	}
}

func TestBuildLineage_PopulatesKeys(t *testing.T) {
	secrets := map[string]string{"API_KEY": "abc", "DB_URL": "postgres://"}
	rec := BuildLineage("https://vault", "secret/app", secrets, nil)

	if len(rec.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(rec.Keys))
	}
	if rec.VaultAddr != "https://vault" {
		t.Errorf("unexpected VaultAddr: %q", rec.VaultAddr)
	}
	if rec.SyncedAt.IsZero() {
		t.Error("SyncedAt should not be zero")
	}
}
