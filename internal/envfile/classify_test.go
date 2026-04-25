package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadClassifications_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "classify.json")

	cm := ClassificationMap{
		"public": {Label: "public", Keys: []string{"APP_NAME", "APP_VERSION"}},
		"secret": {Label: "secret", Keys: []string{"DB_PASSWORD", "API_KEY"}},
	}

	if err := SaveClassifications(path, cm); err != nil {
		t.Fatalf("SaveClassifications: %v", err)
	}

	loaded, err := LoadClassifications(path)
	if err != nil {
		t.Fatalf("LoadClassifications: %v", err)
	}

	if len(loaded) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(loaded))
	}
	if len(loaded["public"].Keys) != 2 {
		t.Errorf("expected 2 public keys, got %d", len(loaded["public"].Keys))
	}
}

func TestLoadClassifications_NonExistent(t *testing.T) {
	cm, err := LoadClassifications("/tmp/no-such-classify-file.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cm) != 0 {
		t.Errorf("expected empty map, got %d entries", len(cm))
	}
}

func TestSaveClassifications_EmptyPath(t *testing.T) {
	if err := SaveClassifications("", ClassificationMap{}); err != nil {
		t.Errorf("expected nil error for empty path, got %v", err)
	}
}

func TestSaveClassifications_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "classify.json")

	if err := SaveClassifications(path, ClassificationMap{}); err != nil {
		t.Fatalf("SaveClassifications: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected perm 0600, got %v", info.Mode().Perm())
	}
}

func TestClassify_FiltersToLabel(t *testing.T) {
	secrets := map[string]string{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "key123",
	}
	cm := ClassificationMap{
		"secret": {Label: "secret", Keys: []string{"DB_PASSWORD", "API_KEY"}},
	}

	out := Classify(secrets, cm, "secret")
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["APP_NAME"]; ok {
		t.Error("APP_NAME should not be in secret classification")
	}
}

func TestClassify_UnknownLabel_Empty(t *testing.T) {
	out := Classify(map[string]string{"K": "v"}, ClassificationMap{}, "missing")
	if len(out) != 0 {
		t.Errorf("expected empty map for unknown label")
	}
}

func TestListLabels_Sorted(t *testing.T) {
	cm := ClassificationMap{
		"zebra":  {Label: "zebra"},
		"alpha":  {Label: "alpha"},
		"middle": {Label: "middle"},
	}
	labels := ListLabels(cm)
	if labels[0] != "alpha" || labels[1] != "middle" || labels[2] != "zebra" {
		t.Errorf("unexpected order: %v", labels)
	}
}
