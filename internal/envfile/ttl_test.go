package envfile

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveAndLoadTTL_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".ttl.json")

	ttl := 10 * time.Minute
	if err := SaveTTL(path, ttl); err != nil {
		t.Fatalf("SaveTTL: %v", err)
	}

	rec, err := LoadTTL(path)
	if err != nil {
		t.Fatalf("LoadTTL: %v", err)
	}

	if rec.TTL != ttl {
		t.Errorf("expected TTL %s, got %s", ttl, rec.TTL)
	}
	if rec.Expired() {
		t.Error("expected record to not be expired")
	}
}

func TestTTLRecord_Expired(t *testing.T) {
	rec := TTLRecord{
		SyncedAt:  time.Now().Add(-2 * time.Hour),
		ExpiresAt: time.Now().Add(-1 * time.Hour),
		TTL:       1 * time.Hour,
	}
	if !rec.Expired() {
		t.Error("expected record to be expired")
	}
}

func TestCheckTTL_NotExpired(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".ttl.json")

	if err := SaveTTL(path, 1*time.Hour); err != nil {
		t.Fatalf("SaveTTL: %v", err)
	}
	if err := CheckTTL(path); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestCheckTTL_Expired(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".ttl.json")

	if err := SaveTTL(path, -1*time.Minute); err != nil {
		t.Fatalf("SaveTTL: %v", err)
	}
	if err := CheckTTL(path); err == nil {
		t.Error("expected error for expired TTL")
	}
}

func TestLoadTTL_NonExistent(t *testing.T) {
	_, err := LoadTTL("/nonexistent/.ttl.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSaveTTL_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".ttl.json")

	if err := SaveTTL(path, 5*time.Minute); err != nil {
		t.Fatalf("SaveTTL: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
