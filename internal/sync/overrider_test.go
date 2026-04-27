package sync_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/envfile"
	syncer "github.com/your-org/vaultpull/internal/sync"
)

func writeOverrides(t *testing.T, dir string, sets []envfile.OverrideSet) string {
	t.Helper()
	path := filepath.Join(dir, "overrides.json")
	data, _ := json.Marshal(sets)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("write overrides: %v", err)
	}
	return path
}

func TestApplyOverrides_NoPath_PassesThrough(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	out, err := syncer.ApplyOverrides(secrets, "", "prod", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" {
		t.Error("expected passthrough")
	}
}

func TestApplyOverrides_AppliesSet(t *testing.T) {
	dir := t.TempDir()
	sets := []envfile.OverrideSet{
		{Name: "prod", Overrides: []envfile.Override{
			{Key: "LOG_LEVEL", Value: "error"},
		}},
	}
	path := writeOverrides(t, dir, sets)
	var buf bytes.Buffer
	secrets := map[string]string{"LOG_LEVEL": "debug"}
	out, err := syncer.ApplyOverrides(secrets, path, "prod", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["LOG_LEVEL"] != "error" {
		t.Errorf("expected error, got %q", out["LOG_LEVEL"])
	}
	if buf.Len() == 0 {
		t.Error("expected log output")
	}
}

func TestApplyOverrides_SetNotFound_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	sets := []envfile.OverrideSet{
		{Name: "dev", Overrides: []envfile.Override{{Key: "X", Value: "y"}}},
	}
	path := writeOverrides(t, dir, sets)
	_, err := syncer.ApplyOverrides(map[string]string{}, path, "ghost", nil)
	if err == nil {
		t.Fatal("expected error for missing set")
	}
}

func TestApplyOverrides_EmptyFile_PassesThrough(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "overrides.json")
	if err := os.WriteFile(path, []byte("[]"), 0600); err != nil {
		t.Fatal(err)
	}
	secrets := map[string]string{"K": "v"}
	out, err := syncer.ApplyOverrides(secrets, path, "any", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["K"] != "v" {
		t.Error("expected passthrough on empty sets")
	}
}
