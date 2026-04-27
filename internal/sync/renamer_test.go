package sync_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/envfile"
	syncer "github.com/yourusername/vaultpull/internal/sync"
)

func writeRenameRules(t *testing.T, dir string, rm envfile.RenameMap) string {
	t.Helper()
	path := filepath.Join(dir, "renames.json")
	data, _ := json.Marshal(rm)
	os.WriteFile(path, data, 0600)
	return path
}

func TestApplyRenames_NoPath_PassesThrough(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	out, err := syncer.ApplyRenames(secrets, "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "val" {
		t.Errorf("expected KEY=val")
	}
}

func TestApplyRenames_AppliesRule(t *testing.T) {
	dir := t.TempDir()
	rm := envfile.RenameMap{Rules: []envfile.RenameRule{{From: "OLD", To: "NEW"}}}
	path := writeRenameRules(t, dir, rm)

	secrets := map[string]string{"OLD": "secret", "OTHER": "x"}
	var buf bytes.Buffer
	out, err := syncer.ApplyRenames(secrets, path, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW"] != "secret" {
		t.Errorf("expected NEW=secret")
	}
	if _, ok := out["OLD"]; ok {
		t.Error("OLD should be removed")
	}
	if buf.Len() == 0 {
		t.Error("expected log output")
	}
}

func TestApplyRenames_EmptyRules_NoChange(t *testing.T) {
	dir := t.TempDir()
	path := writeRenameRules(t, dir, envfile.RenameMap{})
	secrets := map[string]string{"A": "1"}
	out, err := syncer.ApplyRenames(secrets, path, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" {
		t.Error("expected A=1 unchanged")
	}
}

func TestAddRenameRule_AppendsRule(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "renames.json")
	if err := syncer.AddRenameRule(path, "FOO", "BAR"); err != nil {
		t.Fatalf("AddRenameRule: %v", err)
	}
	rm, err := envfile.LoadRenames(path)
	if err != nil {
		t.Fatalf("LoadRenames: %v", err)
	}
	if len(rm.Rules) != 1 || rm.Rules[0].From != "FOO" || rm.Rules[0].To != "BAR" {
		t.Errorf("unexpected rules: %+v", rm.Rules)
	}
}

func TestAddRenameRule_EmptyPath_ReturnsError(t *testing.T) {
	if err := syncer.AddRenameRule("", "A", "B"); err == nil {
		t.Error("expected error for empty path")
	}
}
