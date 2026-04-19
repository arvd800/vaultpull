package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadPins_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")

	pins := PinMap{}
	pins = PinKey(pins, "DB_PASS", "supersecret", "alice")
	pins = PinKey(pins, "API_KEY", "abc123", "")

	if err := SavePins(path, pins); err != nil {
		t.Fatalf("SavePins: %v", err)
	}

	loaded, err := LoadPins(path)
	if err != nil {
		t.Fatalf("LoadPins: %v", err)
	}
	if loaded["DB_PASS"].Value != "supersecret" {
		t.Errorf("expected supersecret, got %s", loaded["DB_PASS"].Value)
	}
	if loaded["DB_PASS"].PinnedBy != "alice" {
		t.Errorf("expected alice, got %s", loaded["DB_PASS"].PinnedBy)
	}
	if loaded["API_KEY"].Value != "abc123" {
		t.Errorf("expected abc123, got %s", loaded["API_KEY"].Value)
	}
}

func TestLoadPins_NonExistent(t *testing.T) {
	pins, err := LoadPins("/nonexistent/pins.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(pins) != 0 {
		t.Errorf("expected empty map")
	}
}

func TestSavePins_EmptyPath(t *testing.T) {
	err := SavePins("", PinMap{})
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestUnpinKey_RemovesKey(t *testing.T) {
	pins := PinMap{}
	pins = PinKey(pins, "FOO", "bar", "")
	pins = PinKey(pins, "BAZ", "qux", "")
	pins = UnpinKey(pins, "FOO")
	if _, ok := pins["FOO"]; ok {
		t.Error("FOO should have been unpinned")
	}
	if pins["BAZ"].Value != "qux" {
		t.Error("BAZ should remain")
	}
}

func TestApplyPins_OverridesValues(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "original", "HOST": "localhost"}
	pins := PinMap{}
	pins = PinKey(pins, "DB_PASS", "pinned_value", "")

	result := ApplyPins(secrets, pins)
	if result["DB_PASS"] != "pinned_value" {
		t.Errorf("expected pinned_value, got %s", result["DB_PASS"])
	}
	if result["HOST"] != "localhost" {
		t.Errorf("expected localhost, got %s", result["HOST"])
	}
	if secrets["DB_PASS"] != "original" {
		t.Error("ApplyPins must not mutate input")
	}
}

func TestSavePins_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")
	if err := SavePins(path, PinMap{}); err != nil {
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
