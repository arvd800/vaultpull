package envfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/envfile"
)

func TestSaveAndLoadOverrides_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "overrides.json")

	sets := []envfile.OverrideSet{
		{Name: "prod", Overrides: []envfile.Override{
			{Key: "LOG_LEVEL", Value: "warn", Condition: "always"},
		}},
	}
	if err := envfile.SaveOverrides(path, sets); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := envfile.LoadOverrides(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded) != 1 || loaded[0].Name != "prod" {
		t.Fatalf("unexpected loaded sets: %+v", loaded)
	}
}

func TestLoadOverrides_NonExistent(t *testing.T) {
	sets, err := envfile.LoadOverrides("/nonexistent/overrides.json")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if sets != nil {
		t.Fatal("expected nil sets")
	}
}

func TestSaveOverrides_EmptyPath(t *testing.T) {
	if err := envfile.SaveOverrides("", nil); err != nil {
		t.Fatalf("expected no error on empty path, got %v", err)
	}
}

func TestApplyOverrides_Always(t *testing.T) {
	sets := []envfile.OverrideSet{
		{Name: "staging", Overrides: []envfile.Override{
			{Key: "DB_HOST", Value: "staging-db", Condition: "always"},
		}},
	}
	secrets := map[string]string{"DB_HOST": "prod-db", "APP": "myapp"}
	out, err := envfile.ApplyOverrides(secrets, sets, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "staging-db" {
		t.Errorf("expected staging-db, got %q", out["DB_HOST"])
	}
	if out["APP"] != "myapp" {
		t.Errorf("expected myapp, got %q", out["APP"])
	}
}

func TestApplyOverrides_Missing(t *testing.T) {
	sets := []envfile.OverrideSet{
		{Name: "dev", Overrides: []envfile.Override{
			{Key: "FEATURE_FLAG", Value: "true", Condition: "missing"},
			{Key: "DB_HOST", Value: "localhost", Condition: "missing"},
		}},
	}
	secrets := map[string]string{"DB_HOST": "prod-db"}
	out, err := envfile.ApplyOverrides(secrets, sets, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FEATURE_FLAG"] != "true" {
		t.Errorf("expected true, got %q", out["FEATURE_FLAG"])
	}
	if out["DB_HOST"] != "prod-db" {
		t.Errorf("missing condition should not overwrite existing key")
	}
}

func TestApplyOverrides_SetNotFound(t *testing.T) {
	_, err := envfile.ApplyOverrides(map[string]string{}, nil, "ghost")
	if err == nil {
		t.Fatal("expected error for missing set")
	}
}

func TestApplyOverrides_UnknownCondition(t *testing.T) {
	sets := []envfile.OverrideSet{
		{Name: "x", Overrides: []envfile.Override{
			{Key: "K", Value: "v", Condition: "sometimes"},
		}},
	}
	_, err := envfile.ApplyOverrides(map[string]string{}, sets, "x")
	if err == nil {
		t.Fatal("expected error for unknown condition")
	}
}

func TestSaveOverrides_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "overrides.json")
	if err := envfile.SaveOverrides(path, []envfile.OverrideSet{}); err != nil {
		t.Fatalf("save: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
