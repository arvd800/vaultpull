package envfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/envfile"
)

func TestSaveAndLoadInherit_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "inherit.json")

	m := envfile.InheritMap{
		Parent: "secret/base",
		Keys: map[string]string{
			"DB_HOST": "DATABASE_HOST",
			"API_KEY": "",
		},
	}

	if err := envfile.SaveInherit(path, m); err != nil {
		t.Fatalf("SaveInherit: %v", err)
	}

	loaded, err := envfile.LoadInherit(path)
	if err != nil {
		t.Fatalf("LoadInherit: %v", err)
	}
	if loaded.Parent != m.Parent {
		t.Errorf("Parent: got %q want %q", loaded.Parent, m.Parent)
	}
	if loaded.Keys["DB_HOST"] != "DATABASE_HOST" {
		t.Errorf("DB_HOST mapping incorrect: %v", loaded.Keys)
	}
}

func TestLoadInherit_NonExistent(t *testing.T) {
	m, err := envfile.LoadInherit("/tmp/no-such-inherit-file.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(m.Keys) != 0 {
		t.Errorf("expected empty keys, got %v", m.Keys)
	}
}

func TestSaveInherit_EmptyPath(t *testing.T) {
	if err := envfile.SaveInherit("", envfile.InheritMap{}); err != nil {
		t.Fatalf("expected no error on empty path, got %v", err)
	}
}

func TestSaveInherit_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "inherit.json")
	m := envfile.InheritMap{Parent: "x", Keys: map[string]string{"A": "B"}}
	if err := envfile.SaveInherit(path, m); err != nil {
		t.Fatalf("SaveInherit: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
}

func TestApplyInherit_FillsMissingKeys(t *testing.T) {
	parent := map[string]string{
		"DATABASE_HOST": "db.prod.internal",
		"API_KEY":       "secret-key-123",
	}
	child := map[string]string{
		"APP_NAME": "myapp",
	}
	m := envfile.InheritMap{
		Parent: "secret/base",
		Keys: map[string]string{
			"DB_HOST": "DATABASE_HOST",
			"API_KEY": "",
		},
	}

	out := envfile.ApplyInherit(child, parent, m)

	if out["DB_HOST"] != "db.prod.internal" {
		t.Errorf("DB_HOST: got %q", out["DB_HOST"])
	}
	if out["API_KEY"] != "secret-key-123" {
		t.Errorf("API_KEY: got %q", out["API_KEY"])
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should be preserved: got %q", out["APP_NAME"])
	}
}

func TestApplyInherit_DoesNotOverwriteExisting(t *testing.T) {
	parent := map[string]string{"API_KEY": "parent-key"}
	child := map[string]string{"API_KEY": "child-key"}
	m := envfile.InheritMap{Keys: map[string]string{"API_KEY": ""}}

	out := envfile.ApplyInherit(child, parent, m)
	if out["API_KEY"] != "child-key" {
		t.Errorf("child value should not be overwritten, got %q", out["API_KEY"])
	}
}

func TestApplyInherit_DoesNotMutateInput(t *testing.T) {
	parent := map[string]string{"X": "1"}
	child := map[string]string{"Y": "2"}
	m := envfile.InheritMap{Keys: map[string]string{"X": ""}}

	_ = envfile.ApplyInherit(child, parent, m)
	if _, ok := child["X"]; ok {
		t.Error("ApplyInherit must not mutate child map")
	}
}
