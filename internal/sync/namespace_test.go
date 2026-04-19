package sync_test

import (
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/envfile"
	syncer "github.com/yourusername/vaultpull/internal/sync"
)

func TestNamespacedSync_FiltersSecrets(t *testing.T) {
	incoming := map[string]string{
		"PROD_HOST": "prod.example.com",
		"PROD_KEY":  "abc",
		"DEV_HOST":  "localhost",
	}
	existing := map[string]string{"OLD_KEY": "old"}
	ns := envfile.Namespace{Name: "prod", Prefix: "PROD_"}

	result, err := syncer.NamespacedSync(incoming, existing, ns, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["HOST"] != "prod.example.com" {
		t.Errorf("expected prod.example.com, got %q", result["HOST"])
	}
	if _, ok := result["DEV_HOST"]; ok {
		t.Error("DEV_HOST should be excluded")
	}
	if result["OLD_KEY"] != "old" {
		t.Error("existing key should be preserved")
	}
}

func TestNamespacedSync_EmptyNamespace_PassesThrough(t *testing.T) {
	incoming := map[string]string{"A": "1", "B": "2"}
	ns := envfile.Namespace{}

	result, err := syncer.NamespacedSync(incoming, map[string]string{}, ns, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestNamespacedSync_PersistsNamespace(t *testing.T) {
	dir := t.TempDir()
	nsPath := filepath.Join(dir, "ns.json")

	ns := envfile.Namespace{Name: "staging", Prefix: "STG_"}
	incoming := map[string]string{"STG_DB": "db.stg"}

	_, err := syncer.NamespacedSync(incoming, map[string]string{}, ns, nsPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, err := envfile.LoadNamespaces(nsPath)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded["staging"].Prefix != "STG_" {
		t.Errorf("expected STG_ prefix in persisted namespace")
	}
}
