package envfile_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/envfile"
)

func TestSaveAndLoadTags_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tags.json")

	record := envfile.TagRecord{
		Source:    "secret/myapp",
		FetchedAt: time.Now().UTC().Truncate(time.Second),
		Tags:      map[string]string{"env": "production", "team": "platform"},
	}

	if err := envfile.SaveTags(path, record); err != nil {
		t.Fatalf("SaveTags: %v", err)
	}

	loaded, err := envfile.LoadTags(path)
	if err != nil {
		t.Fatalf("LoadTags: %v", err)
	}

	if loaded.Source != record.Source {
		t.Errorf("Source: got %q, want %q", loaded.Source, record.Source)
	}
	if loaded.Tags["env"] != "production" {
		t.Errorf("tag env: got %q", loaded.Tags["env"])
	}
}

func TestLoadTags_NonExistent(t *testing.T) {
	record, err := envfile.LoadTags("/nonexistent/tags.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if record.Tags == nil {
		t.Error("expected non-nil Tags map")
	}
}

func TestSaveTags_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tags.json")

	record := envfile.TagRecord{Source: "s", FetchedAt: time.Now(), Tags: map[string]string{}}
	if err := envfile.SaveTags(path, record); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
}

func TestMergeTags_AddsKeys(t *testing.T) {
	record := &envfile.TagRecord{Tags: map[string]string{"a": "1"}}
	envfile.MergeTags(record, map[string]string{"b": "2", "a": "overwritten"})

	if record.Tags["b"] != "2" {
		t.Errorf("expected b=2, got %q", record.Tags["b"])
	}
	if record.Tags["a"] != "overwritten" {
		t.Errorf("expected a=overwritten, got %q", record.Tags["a"])
	}
}
