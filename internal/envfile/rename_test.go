package envfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/envfile"
)

func TestSaveAndLoadRenames_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "renames.json")

	rm := envfile.RenameMap{
		Rules: []envfile.RenameRule{
			{From: "OLD_KEY", To: "NEW_KEY"},
		},
	}
	if err := envfile.SaveRenames(path, rm); err != nil {
		t.Fatalf("SaveRenames: %v", err)
	}
	loaded, err := envfile.LoadRenames(path)
	if err != nil {
		t.Fatalf("LoadRenames: %v", err)
	}
	if len(loaded.Rules) != 1 || loaded.Rules[0].From != "OLD_KEY" || loaded.Rules[0].To != "NEW_KEY" {
		t.Errorf("unexpected rules: %+v", loaded.Rules)
	}
}

func TestLoadRenames_NonExistent(t *testing.T) {
	rm, err := envfile.LoadRenames("/nonexistent/renames.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(rm.Rules) != 0 {
		t.Errorf("expected empty rules")
	}
}

func TestSaveRenames_EmptyPath(t *testing.T) {
	if err := envfile.SaveRenames("", envfile.RenameMap{}); err != nil {
		t.Errorf("expected no error for empty path, got %v", err)
	}
}

func TestSaveRenames_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "renames.json")
	if err := envfile.SaveRenames(path, envfile.RenameMap{}); err != nil {
		t.Fatalf("SaveRenames: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}

func TestApplyRenames_RenamesKey(t *testing.T) {
	secrets := map[string]string{"OLD_KEY": "value1", "KEEP": "value2"}
	rm := envfile.RenameMap{Rules: []envfile.RenameRule{{From: "OLD_KEY", To: "NEW_KEY"}}}
	out := envfile.ApplyRenames(secrets, rm)
	if _, ok := out["OLD_KEY"]; ok {
		t.Error("OLD_KEY should have been removed")
	}
	if out["NEW_KEY"] != "value1" {
		t.Errorf("expected NEW_KEY=value1, got %q", out["NEW_KEY"])
	}
	if out["KEEP"] != "value2" {
		t.Errorf("expected KEEP=value2, got %q", out["KEEP"])
	}
}

func TestApplyRenames_SkipsMissingFrom(t *testing.T) {
	secrets := map[string]string{"EXISTING": "v"}
	rm := envfile.RenameMap{Rules: []envfile.RenameRule{{From: "MISSING", To: "NEW"}}}
	out := envfile.ApplyRenames(secrets, rm)
	if _, ok := out["NEW"]; ok {
		t.Error("NEW should not exist when FROM key is absent")
	}
	if out["EXISTING"] != "v" {
		t.Error("EXISTING should be preserved")
	}
}

func TestApplyRenames_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	rm := envfile.RenameMap{Rules: []envfile.RenameRule{{From: "A", To: "B"}}}
	envfile.ApplyRenames(secrets, rm)
	if _, ok := secrets["A"]; !ok {
		t.Error("original map should not be mutated")
	}
}
