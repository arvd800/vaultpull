package envfile_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/envfile"
)

func TestSaveAndLoadAnnotations_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "annotations.json")

	annotations := envfile.Annotations{
		"DB_PASSWORD": {Note: "rotated quarterly", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		"API_KEY": {Note: "from team vault", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
	}

	if err := envfile.SaveAnnotations(path, annotations); err != nil {
		t.Fatalf("SaveAnnotations: %v", err)
	}

	loaded, err := envfile.LoadAnnotations(path)
	if err != nil {
		t.Fatalf("LoadAnnotations: %v", err)
	}
	if len(loaded) != 2 {
		t.Errorf("expected 2 annotations, got %d", len(loaded))
	}
	if loaded["DB_PASSWORD"].Note != "rotated quarterly" {
		t.Errorf("unexpected note: %s", loaded["DB_PASSWORD"].Note)
	}
}

func TestLoadAnnotations_NonExistent(t *testing.T) {
	annotations, err := envfile.LoadAnnotations("/nonexistent/annotations.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(annotations) != 0 {
		t.Errorf("expected empty map, got %d entries", len(annotations))
	}
}

func TestSaveAnnotations_EmptyPath(t *testing.T) {
	if err := envfile.SaveAnnotations("", envfile.Annotations{"K": {Note: "x"}}); err != nil {
		t.Errorf("expected no error for empty path, got: %v", err)
	}
}

func TestSaveAnnotations_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "annotations.json")

	if err := envfile.SaveAnnotations(path, envfile.Annotations{}); err != nil {
		t.Fatalf("SaveAnnotations: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
}

func TestAnnotate_AddsNewKey(t *testing.T) {
	existing := envfile.Annotations{}
	out := envfile.Annotate(existing, "SECRET_KEY", "important secret")
	if a, ok := out["SECRET_KEY"]; !ok {
		t.Fatal("expected key to be present")
	} else if a.Note != "important secret" {
		t.Errorf("unexpected note: %s", a.Note)
	}
}

func TestAnnotate_UpdatesExistingKey(t *testing.T) {
	created := time.Now().UTC().Add(-time.Hour)
	existing := envfile.Annotations{
		"TOKEN": {Note: "old note", CreatedAt: created, UpdatedAt: created},
	}
	out := envfile.Annotate(existing, "TOKEN", "new note")
	if out["TOKEN"].Note != "new note" {
		t.Errorf("expected updated note")
	}
	if !out["TOKEN"].CreatedAt.Equal(created) {
		t.Errorf("CreatedAt should not change on update")
	}
}

func TestAnnotate_DoesNotMutateInput(t *testing.T) {
	existing := envfile.Annotations{"A": {Note: "original"}}
	_ = envfile.Annotate(existing, "A", "changed")
	if existing["A"].Note != "original" {
		t.Error("input was mutated")
	}
}

func TestRemoveAnnotation_RemovesKey(t *testing.T) {
	existing := envfile.Annotations{
		"A": {Note: "keep"},
		"B": {Note: "remove"},
	}
	out := envfile.RemoveAnnotation(existing, "B")
	if _, ok := out["B"]; ok {
		t.Error("expected key B to be removed")
	}
	if _, ok := out["A"]; !ok {
		t.Error("expected key A to remain")
	}
}
