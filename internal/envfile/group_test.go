package envfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/envfile"
)

func TestSaveAndLoadGroups_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "groups.json")

	groups := envfile.GroupMap{
		"backend": {"DB_HOST", "DB_PORT"},
		"frontend": {"API_URL"},
	}

	if err := envfile.SaveGroups(path, groups); err != nil {
		t.Fatalf("SaveGroups: %v", err)
	}

	loaded, err := envfile.LoadGroups(path)
	if err != nil {
		t.Fatalf("LoadGroups: %v", err)
	}

	if len(loaded["backend"]) != 2 {
		t.Errorf("expected 2 backend keys, got %d", len(loaded["backend"]))
	}
	if len(loaded["frontend"]) != 1 {
		t.Errorf("expected 1 frontend key, got %d", len(loaded["frontend"]))
	}
}

func TestLoadGroups_NonExistent(t *testing.T) {
	groups, err := envfile.LoadGroups("/nonexistent/groups.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(groups) != 0 {
		t.Errorf("expected empty map, got %v", groups)
	}
}

func TestSaveGroups_EmptyPath(t *testing.T) {
	if err := envfile.SaveGroups("", envfile.GroupMap{"g": {"K"}}); err != nil {
		t.Errorf("expected no error for empty path, got %v", err)
	}
}

func TestAddGroup_DoesNotMutateInput(t *testing.T) {
	orig := envfile.GroupMap{"a": {"X"}}
	out := envfile.AddGroup(orig, "b", []string{"Y", "Z"})
	if _, ok := orig["b"]; ok {
		t.Error("original map was mutated")
	}
	if len(out["b"]) != 2 {
		t.Errorf("expected 2 keys in new group, got %d", len(out["b"]))
	}
}

func TestApplyGroup_FiltersSecrets(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"API_URL": "https://example.com",
	}
	groups := envfile.GroupMap{"backend": {"DB_HOST", "DB_PORT"}}

	out, err := envfile.ApplyGroup(secrets, groups, "backend")
	if err != nil {
		t.Fatalf("ApplyGroup: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["API_URL"]; ok {
		t.Error("API_URL should not be in output")
	}
}

func TestApplyGroup_UnknownGroup(t *testing.T) {
	groups := envfile.GroupMap{"backend": {"DB_HOST"}}
	_, err := envfile.ApplyGroup(map[string]string{}, groups, "unknown")
	if err == nil {
		t.Error("expected error for unknown group")
	}
}

func TestSaveGroups_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "groups.json")
	if err := envfile.SaveGroups(path, envfile.GroupMap{"g": {"K"}}); err != nil {
		t.Fatalf("SaveGroups: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
