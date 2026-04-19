package envfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/envfile"
)

func TestSaveAndLoadNamespaces_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "namespaces.json")

	ns := envfile.NamespaceMap{
		"prod": {Name: "prod", Prefix: "PROD_", Tags: map[string]string{"env": "production"}},
		"dev":  {Name: "dev", Prefix: "DEV_"},
	}

	if err := envfile.SaveNamespaces(path, ns); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := envfile.LoadNamespaces(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded["prod"].Prefix != "PROD_" {
		t.Errorf("expected PROD_ prefix, got %q", loaded["prod"].Prefix)
	}
	if loaded["dev"].Name != "dev" {
		t.Errorf("expected dev name, got %q", loaded["dev"].Name)
	}
}

func TestLoadNamespaces_NonExistent(t *testing.T) {
	ns, err := envfile.LoadNamespaces("/nonexistent/path.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ns) != 0 {
		t.Errorf("expected empty map")
	}
}

func TestSaveNamespaces_EmptyPath(t *testing.T) {
	if err := envfile.SaveNamespaces("", envfile.NamespaceMap{}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSaveNamespaces_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ns.json")
	_ = envfile.SaveNamespaces(path, envfile.NamespaceMap{"x": {Name: "x"}})
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}

func TestApplyNamespace_FiltersAndStrips(t *testing.T) {
	secrets := map[string]string{
		"PROD_DB_HOST": "db.prod",
		"PROD_API_KEY": "key123",
		"DEV_DB_HOST":  "localhost",
	}
	ns := envfile.Namespace{Name: "prod", Prefix: "PROD_"}
	out := envfile.ApplyNamespace(secrets, ns)
	if out["DB_HOST"] != "db.prod" {
		t.Errorf("expected db.prod, got %q", out["DB_HOST"])
	}
	if _, ok := out["DEV_DB_HOST"]; ok {
		t.Error("DEV key should be excluded")
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestApplyNamespace_EmptyPrefix(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	ns := envfile.Namespace{Name: "all", Prefix: ""}
	out := envfile.ApplyNamespace(secrets, ns)
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}
