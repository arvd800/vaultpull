package envfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/envfile"
)

func TestSaveAndLoadScopes_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "scopes.json")

	scopes := envfile.ScopeMap{
		"production": {Name: "production", Keys: []string{"DB_URL", "API_KEY"}},
		"staging": {Name: "staging", Keys: []string{"DB_URL"}},
	}

	if err := envfile.SaveScopes(path, scopes); err != nil {
		t.Fatalf("SaveScopes: %v", err)
	}

	loaded, err := envfile.LoadScopes(path)
	if err != nil {
		t.Fatalf("LoadScopes: %v", err)
	}

	if len(loaded) != 2 {
		t.Errorf("expected 2 scopes, got %d", len(loaded))
	}
	if loaded["production"].Keys[1] != "API_KEY" {
		t.Errorf("unexpected key: %v", loaded["production"].Keys)
	}
}

func TestLoadScopes_NonExistent(t *testing.T) {
	scopes, err := envfile.LoadScopes("/nonexistent/scopes.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(scopes) != 0 {
		t.Errorf("expected empty map")
	}
}

func TestSaveScopes_EmptyPath(t *testing.T) {
	if err := envfile.SaveScopes("", envfile.ScopeMap{}); err != nil {
		t.Errorf("expected no error for empty path, got %v", err)
	}
}

func TestApplyScope_FiltersKeys(t *testing.T) {
	scopes := envfile.ScopeMap{
		"prod": {Name: "prod", Keys: []string{"DB_URL", "API_KEY"}},
	}
	secrets := map[string]string{
		"DB_URL":  "postgres://localhost",
		"API_KEY": "secret",
		"DEBUG":   "true",
	}

	result, err := envfile.ApplyScope(secrets, scopes, "prod")
	if err != nil {
		t.Fatalf("ApplyScope: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["DEBUG"]; ok {
		t.Error("DEBUG should not be in scoped result")
	}
}

func TestApplyScope_UnknownScope(t *testing.T) {
	scopes := envfile.ScopeMap{}
	_, err := envfile.ApplyScope(map[string]string{}, scopes, "missing")
	if err == nil {
		t.Error("expected error for unknown scope")
	}
}

func TestApplyScope_EmptyScopeName(t *testing.T) {
	scopes := envfile.ScopeMap{
		"prod": {Name: "prod", Keys: []string{"DB_URL"}},
	}
	_, err := envfile.ApplyScope(map[string]string{"DB_URL": "val"}, scopes, "")
	if err == nil {
		t.Error("expected error for empty scope name")
	}
}

func TestSaveScopes_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "scopes.json")
	if err := envfile.SaveScopes(path, envfile.ScopeMap{}); err != nil {
		t.Fatalf("SaveScopes: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}

func TestListScopes_ReturnsSorted(t *testing.T) {
	scopes := envfile.ScopeMap{
		"staging":    {Name: "staging"},
		"production": {Name: "production"},
		"dev":        {Name: "dev"},
	}
	names := envfile.ListScopes(scopes)
	if names[0] != "dev" || names[1] != "production" || names[2] != "staging" {
		t.Errorf("unexpected order: %v", names)
	}
}
