package sync

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/envfile"
)

func TestApplyLocks_SkipsLockedKey(t *testing.T) {
	incoming := map[string]string{"DB_PASS": "new", "API_KEY": "fresh"}
	existing := map[string]string{"DB_PASS": "old"}
	locks := envfile.LockKey(envfile.LockFile{}, "DB_PASS", "hold")

	merged, skipped := ApplyLocks(incoming, existing, locks)

	if merged["DB_PASS"] != "old" {
		t.Errorf("expected old value preserved, got %q", merged["DB_PASS"])
	}
	if merged["API_KEY"] != "fresh" {
		t.Errorf("expected fresh API_KEY, got %q", merged["API_KEY"])
	}
	if len(skipped) != 1 || skipped[0] != "DB_PASS" {
		t.Errorf("unexpected skipped: %v", skipped)
	}
}

func TestApplyLocks_LockedKeyNotInExisting(t *testing.T) {
	incoming := map[string]string{"SECRET": "value"}
	existing := map[string]string{}
	locks := envfile.LockKey(envfile.LockFile{}, "SECRET", "")

	merged, skipped := ApplyLocks(incoming, existing, locks)

	if _, ok := merged["SECRET"]; ok {
		t.Error("locked key with no existing value should be removed")
	}
	if len(skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(skipped))
	}
}

func TestApplyLocks_NoLocks(t *testing.T) {
	incoming := map[string]string{"A": "1", "B": "2"}
	merged, skipped := ApplyLocks(incoming, map[string]string{}, envfile.LockFile{})
	if len(skipped) != 0 {
		t.Error("expected no skipped keys")
	}
	if merged["A"] != "1" || merged["B"] != "2" {
		t.Error("all keys should pass through")
	}
}

func TestApplyLocks_DoesNotMutateIncoming(t *testing.T) {
	incoming := map[string]string{"TOKEN": "new"}
	existing := map[string]string{"TOKEN": "old"}
	locks := envfile.LockKey(envfile.LockFile{}, "TOKEN", "")

	ApplyLocks(incoming, existing, locks)

	if incoming["TOKEN"] != "new" {
		t.Error("incoming map was mutated")
	}
}
