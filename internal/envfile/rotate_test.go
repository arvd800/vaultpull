package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEncrypted(t *testing.T, data map[string]string, pass string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.enc")
	if err := WriteEncryptedFile(path, data, pass); err != nil {
		t.Fatalf("setup: %v", err)
	}
	return path
}

func TestRotate_RoundTrip(t *testing.T) {
	original := map[string]string{"DB_PASS": "secret", "API_KEY": "abc123"}
	path := writeTempEncrypted(t, original, "old-pass")

	result, err := Rotate(path, "old-pass", "new-pass")
	if err != nil {
		t.Fatalf("Rotate: %v", err)
	}
	if len(result.Rotated) != 2 {
		t.Errorf("expected 2 rotated, got %d", len(result.Rotated))
	}

	decrypted, err := ReadDecrypted(path, "new-pass")
	if err != nil {
		t.Fatalf("ReadDecrypted with new pass: %v", err)
	}
	for k, v := range original {
		if decrypted[k] != v {
			t.Errorf("key %s: want %q, got %q", k, v, decrypted[k])
		}
	}
}

func TestRotate_WrongOldPassphrase(t *testing.T) {
	path := writeTempEncrypted(t, map[string]string{"X": "y"}, "correct")
	_, err := Rotate(path, "wrong", "new-pass")
	if err == nil {
		t.Fatal("expected error with wrong old passphrase")
	}
}

func TestRotate_EmptyPassphrase(t *testing.T) {
	path := writeTempEncrypted(t, map[string]string{"X": "y"}, "pass")
	_, err := Rotate(path, "", "new")
	if err == nil {
		t.Fatal("expected error for empty old passphrase")
	}
	_, err = Rotate(path, "pass", "")
	if err == nil {
		t.Fatal("expected error for empty new passphrase")
	}
}

func TestRotate_NonExistentFile(t *testing.T) {
	_, err := Rotate("/nonexistent/.env.enc", "old", "new")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRotateKeys_PartialRotation(t *testing.T) {
	original := map[string]string{"A": "1", "B": "2", "C": "3"}
	path := writeTempEncrypted(t, original, "pass")

	result, err := RotateKeys(path, []string{"A", "B"}, "pass", "new-pass")
	if err != nil {
		t.Fatalf("RotateKeys: %v", err)
	}
	if len(result.Rotated) != 2 {
		t.Errorf("expected 2 rotated, got %d", len(result.Rotated))
	}
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}

	// Verify all keys readable with new pass (full file re-encrypted)
	decrypted, err := ReadDecrypted(path, "new-pass")
	if err != nil {
		t.Fatalf("ReadDecrypted: %v", err)
	}
	if decrypted["C"] != "3" {
		t.Errorf("skipped key C should still be present")
	}
}

func TestRotate_OldPassNoLongerWorks(t *testing.T) {
	path := writeTempEncrypted(t, map[string]string{"K": "v"}, "old")
	if _, err := Rotate(path, "old", "new"); err != nil {
		t.Fatalf("Rotate: %v", err)
	}
	_, err := ReadDecrypted(path, "old")
	if err == nil {
		t.Fatal("old passphrase should no longer work after rotation")
	}
	os.Remove(path)
}
