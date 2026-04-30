package sync

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/envfile"
)

func writeCondenseConfig(t *testing.T, dir string, cfg envfile.CondenseConfig) string {
	t.Helper()
	path := filepath.Join(dir, "condense.json")
	data, _ := json.MarshalIndent(cfg, "", "  ")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("write condense config: %v", err)
	}
	return path
}

func TestApplyCondense_NoPath_PassesThrough(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	out, err := ApplyCondense(secrets, "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" {
		t.Errorf("expected passthrough, got %v", out)
	}
}

func TestApplyCondense_AppliesRule(t *testing.T) {
	dir := t.TempDir()
	cfg := envfile.CondenseConfig{
		Rules: []envfile.CondenseRule{
			{OutputKey: "HOST_PORT", SourceKeys: []string{"HOST", "PORT"}, Separator: ":"},
		},
	}
	path := writeCondenseConfig(t, dir, cfg)
	secrets := map[string]string{"HOST": "db", "PORT": "5432"}
	var buf bytes.Buffer
	out, err := ApplyCondense(secrets, path, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOST_PORT"] != "db:5432" {
		t.Errorf("expected 'db:5432', got %q", out["HOST_PORT"])
	}
	if buf.String() == "" {
		t.Error("expected log output")
	}
}

func TestApplyCondense_EmptyRules_LogsWarning(t *testing.T) {
	dir := t.TempDir()
	path := writeCondenseConfig(t, dir, envfile.CondenseConfig{})
	secrets := map[string]string{"X": "y"}
	var buf bytes.Buffer
	out, err := ApplyCondense(secrets, path, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["X"] != "y" {
		t.Errorf("expected passthrough")
	}
	if buf.String() == "" {
		t.Error("expected warning in output")
	}
}

func TestApplyCondense_MissingSourceKey_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	cfg := envfile.CondenseConfig{
		Rules: []envfile.CondenseRule{
			{OutputKey: "OUT", SourceKeys: []string{"MISSING"}, Separator: ""},
		},
	}
	path := writeCondenseConfig(t, dir, cfg)
	_, err := ApplyCondense(map[string]string{}, path, nil)
	if err == nil {
		t.Fatal("expected error for missing source key")
	}
}
