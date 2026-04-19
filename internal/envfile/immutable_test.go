package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadImmutable_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "immutable.json")

	rec := ImmutableRecord{Keys: map[string]bool{"API_KEY": true, "DB_PASS": true}}
	if err := SaveImmutable(path, rec); err != nil {
		t.Fatalf("SaveImmutable: %v", err)
	}

	loaded, err := LoadImmutable(path)
	if err != nil {
		t.Fatalf("LoadImmutable: %v", err)
	}
	if !loaded.Keys["API_KEY"] || !loaded.Keys["DB_PASS"] {
		t.Errorf("expected keys to be present, got %v", loaded.Keys)
	}
}

func TestLoadImmutable_NonExistent(t *testing.T) {
	rec, err := LoadImmutable("/nonexistent/path/immutable.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(rec.Keys) != 0 {
		t.Errorf("expected empty keys, got %v", rec.Keys)
	}
}

func TestSaveImmutable_EmptyPath(t *testing.T) {
	if err := SaveImmutable("", ImmutableRecord{}); err != nil {
		t.Errorf("expected no error for empty path, got %v", err)
	}
}

func TestMarkImmutable_AddsKey(t *testing.T) {
	rec := ImmutableRecord{Keys: map[string]bool{}}
	rec = MarkImmutable(rec, "SECRET")
	if !rec.Keys["SECRET"] {
		t.Error("expected SECRET to be immutable")
	}
}

func TestApplyImmutable_RestoresExistingValue(t *testing.T) {
	existing := map[string]string{"API_KEY": "old", "OTHER": "x"}
	incoming := map[string]string{"API_KEY": "new", "OTHER": "y"}
	rec := ImmutableRecord{Keys: map[string]bool{"API_KEY": true}}

	out := ApplyImmutable(existing, incoming, rec)
	if out["API_KEY"] != "old" {
		t.Errorf("expected old value, got %q", out["API_KEY"])
	}
	if out["OTHER"] != "y" {
		t.Errorf("expected new value for OTHER, got %q", out["OTHER"])
	}
}

func TestApplyImmutable_DropsKeyNotInExisting(t *testing.T) {
	existing := map[string]string{}
	incoming := map[string]string{"API_KEY": "new"}
	rec := ImmutableRecord{Keys: map[string]bool{"API_KEY": true}}

	out := ApplyImmutable(existing, incoming, rec)
	if _, ok := out["API_KEY"]; ok {
		t.Error("expected API_KEY to be dropped")
	}
}

func TestSaveImmutable_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "immutable.json")
	rec := ImmutableRecord{Keys: map[string]bool{"X": true}}
	if err := SaveImmutable(path, rec); err != nil {
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
