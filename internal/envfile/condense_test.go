package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCondense_JoinsKeys(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	cfg := CondenseConfig{
		Rules: []CondenseRule{
			{OutputKey: "DB_ADDR", SourceKeys: []string{"DB_HOST", "DB_PORT"}, Separator: ":"},
		},
	}
	out, err := Condense(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_ADDR"] != "localhost:5432" {
		t.Errorf("expected 'localhost:5432', got %q", out["DB_ADDR"])
	}
}

func TestCondense_DropSources(t *testing.T) {
	secrets := map[string]string{
		"FIRST": "hello",
		"SECOND": "world",
	}
	cfg := CondenseConfig{
		Rules: []CondenseRule{
			{OutputKey: "GREETING", SourceKeys: []string{"FIRST", "SECOND"}, Separator: " ", DropSources: true},
		},
	}
	out, err := Condense(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["GREETING"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", out["GREETING"])
	}
	if _, ok := out["FIRST"]; ok {
		t.Error("expected FIRST to be dropped")
	}
	if _, ok := out["SECOND"]; ok {
		t.Error("expected SECOND to be dropped")
	}
}

func TestCondense_MissingSourceKey_ReturnsError(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost"}
	cfg := CondenseConfig{
		Rules: []CondenseRule{
			{OutputKey: "DB_ADDR", SourceKeys: []string{"DB_HOST", "DB_PORT"}, Separator: ":"},
		},
	}
	_, err := Condense(secrets, cfg)
	if err == nil {
		t.Fatal("expected error for missing source key")
	}
}

func TestCondense_NilInput(t *testing.T) {
	cfg := CondenseConfig{}
	out, err := Condense(nil, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestCondense_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	cfg := CondenseConfig{
		Rules: []CondenseRule{
			{OutputKey: "AB", SourceKeys: []string{"A", "B"}, Separator: "-", DropSources: true},
		},
	}
	_, err := Condense(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := secrets["A"]; !ok {
		t.Error("input map was mutated: key A missing")
	}
}

func TestSaveAndLoadCondenseConfig_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "condense.json")
	cfg := CondenseConfig{
		Rules: []CondenseRule{
			{OutputKey: "OUT", SourceKeys: []string{"X", "Y"}, Separator: "_", DropSources: false},
		},
	}
	if err := SaveCondenseConfig(path, cfg); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadCondenseConfig(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Rules) != 1 || loaded.Rules[0].OutputKey != "OUT" {
		t.Errorf("unexpected loaded config: %+v", loaded)
	}
}

func TestLoadCondenseConfig_NonExistent(t *testing.T) {
	cfg, err := LoadCondenseConfig("/nonexistent/condense.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Rules) != 0 {
		t.Errorf("expected empty config")
	}
}

func TestSaveCondenseConfig_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "condense.json")
	if err := SaveCondenseConfig(path, CondenseConfig{}); err != nil {
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

func TestListCondenseOutputKeys(t *testing.T) {
	cfg := CondenseConfig{
		Rules: []CondenseRule{
			{OutputKey: "Z_KEY"},
			{OutputKey: "A_KEY"},
			{OutputKey: "M_KEY"},
		},
	}
	keys := ListCondenseOutputKeys(cfg)
	if len(keys) != 3 || keys[0] != "A_KEY" || keys[2] != "Z_KEY" {
		t.Errorf("unexpected keys: %v", keys)
	}
}
