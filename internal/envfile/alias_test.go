package envfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/envfile"
)

func TestSaveAndLoadAliases_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "aliases.json")

	aliases := envfile.AliasMap{
		"DB_PASS": "DATABASE_PASSWORD",
		"API": "API_KEY",
	}
	if err := envfile.SaveAliases(path, aliases); err != nil {
		t.Fatalf("SaveAliases: %v", err)
	}
	loaded, err := envfile.LoadAliases(path)
	if err != nil {
		t.Fatalf("LoadAliases: %v", err)
	}
	if loaded["DB_PASS"] != "DATABASE_PASSWORD" {
		t.Errorf("expected DB_PASS -> DATABASE_PASSWORD, got %q", loaded["DB_PASS"])
	}
}

func TestLoadAliases_NonExistent(t *testing.T) {
	aliases, err := envfile.LoadAliases("/nonexistent/aliases.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(aliases) != 0 {
		t.Errorf("expected empty map, got %v", aliases)
	}
}

func TestSaveAliases_EmptyPath(t *testing.T) {
	if err := envfile.SaveAliases("", envfile.AliasMap{"A": "B"}); err != nil {
		t.Errorf("expected no error for empty path, got %v", err)
	}
}

func TestSaveAliases_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "aliases.json")
	if err := envfile.SaveAliases(path, envfile.AliasMap{"X": "Y"}); err != nil {
		t.Fatalf("SaveAliases: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}

func TestApplyAliases_AddsAliasedKeys(t *testing.T) {
	secrets := map[string]string{
		"DATABASE_PASSWORD": "s3cr3t",
		"API_KEY": "abc123",
	}
	aliases := envfile.AliasMap{
		"DB_PASS": "DATABASE_PASSWORD",
		"MISSING": "DOES_NOT_EXIST",
	}
	out := envfile.ApplyAliases(secrets, aliases)
	if out["DB_PASS"] != "s3cr3t" {
		t.Errorf("expected DB_PASS=s3cr3t, got %q", out["DB_PASS"])
	}
	if _, ok := out["MISSING"]; ok {
		t.Error("MISSING should not be present when canonical key absent")
	}
	if out["DATABASE_PASSWORD"] != "s3cr3t" {
		t.Error("original key should still be present")
	}
}

func TestApplyAliases_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	aliases := envfile.AliasMap{"BAZ": "FOO"}
	_ = envfile.ApplyAliases(secrets, aliases)
	if _, ok := secrets["BAZ"]; ok {
		t.Error("ApplyAliases must not mutate input secrets map")
	}
}
