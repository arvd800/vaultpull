package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadLocks_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "locks.json")

	locks := LockFile{}
	locks = LockKey(locks, "DB_PASSWORD", "manual hold")
	locks = LockKey(locks, "API_SECRET", "")

	if err := SaveLocks(path, locks); err != nil {
		t.Fatalf("SaveLocks: %v", err)
	}

	loaded, err := LoadLocks(path)
	if err != nil {
		t.Fatalf("LoadLocks: %v", err)
	}
	if len(loaded) != 2 {
		t.Errorf("expected 2 locks, got %d", len(loaded))
	}
	if loaded["DB_PASSWORD"].Reason != "manual hold" {
		t.Errorf("expected reason 'manual hold', got %q", loaded["DB_PASSWORD"].Reason)
	}
}

func TestLoadLocks_NonExistent(t *testing.T) {
	locks, err := LoadLocks("/tmp/does-not-exist-locks.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(locks) != 0 {
		t.Errorf("expected empty locks")
	}
}

func TestSaveLocks_EmptyPath(t *testing.T) {
	err := SaveLocks("", LockFile{})
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestLockKey_DoesNotMutateInput(t *testing.T) {
	orig := LockFile{}
	updated := LockKey(orig, "SECRET", "reason")
	if IsLocked(orig, "SECRET") {
		t.Error("original LockFile was mutated")
	}
	if !IsLocked(updated, "SECRET") {
		t.Error("updated LockFile missing key")
	}
}

func TestUnlockKey_RemovesKey(t *testing.T) {
	locks := LockKey(LockFile{}, "TOKEN", "")
	locks = UnlockKey(locks, "TOKEN")
	if IsLocked(locks, "TOKEN") {
		t.Error("key should have been unlocked")
	}
}

func TestSaveLocks_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "locks.json")
	if err := SaveLocks(path, LockFile{}); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
