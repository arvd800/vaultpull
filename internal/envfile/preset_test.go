package envfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/envfile"
)

func TestSaveAndLoadPresets_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	presets := []envfile.Preset{
		{Name: "dev", Values: map[string]string{"ENV": "development", "DEBUG": "true"}},
		{Name: "prod", Values: map[string]string{"ENV": "production", "DEBUG": "false"}},
	}
	if err := envfile.SavePresets(dir, presets); err != nil {
		t.Fatalf("SavePresets: %v", err)
	}
	got, err := envfile.LoadPresets(dir)
	if err != nil {
		t.Fatalf("LoadPresets: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 presets, got %d", len(got))
	}
	if got[0].Name != "dev" || got[0].Values["ENV"] != "development" {
		t.Errorf("unexpected preset[0]: %+v", got[0])
	}
}

func TestLoadPresets_NonExistent(t *testing.T) {
	dir := t.TempDir()
	presets, err := envfile.LoadPresets(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(presets) != 0 {
		t.Errorf("expected empty slice, got %d", len(presets))
	}
}

func TestSavePresets_EmptyPath(t *testing.T) {
	if err := envfile.SavePresets("", nil); err != nil {
		t.Errorf("expected no error for empty path, got %v", err)
	}
}

func TestSavePresets_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	if err := envfile.SavePresets(dir, []envfile.Preset{{Name: "test", Values: map[string]string{"K": "V"}}}}); err != nil {
		t.Fatalf("SavePresets: %v", err)
	}
	info, err := os.Stat(filepath.Join(dir, ".vaultpull.presets.json"))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}

func TestApplyPreset_OverridesValues(t *testing.T) {
	presets := []envfile.Preset{
		{Name: "staging", Values: map[string]string{"ENV": "staging", "NEW_KEY": "hello"}},
	}
	secrets := map[string]string{"ENV": "dev", "DB_URL": "postgres://localhost"}
	out, err := envfile.ApplyPreset("staging", presets, secrets)
	if err != nil {
		t.Fatalf("ApplyPreset: %v", err)
	}
	if out["ENV"] != "staging" {
		t.Errorf("expected ENV=staging, got %q", out["ENV"])
	}
	if out["DB_URL"] != "postgres://localhost" {
		t.Errorf("expected DB_URL preserved, got %q", out["DB_URL"])
	}
	if out["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", out["NEW_KEY"])
	}
}

func TestApplyPreset_NotFound(t *testing.T) {
	_, err := envfile.ApplyPreset("missing", []envfile.Preset{}, map[string]string{})
	if err == nil {
		t.Error("expected error for missing preset")
	}
}

func TestApplyPreset_DoesNotMutateInput(t *testing.T) {
	presets := []envfile.Preset{
		{Name: "p", Values: map[string]string{"A": "1"}},
	}
	secrets := map[string]string{"B": "2"}
	_, err := envfile.ApplyPreset("p", presets, secrets)
	if err != nil {
		t.Fatalf("ApplyPreset: %v", err)
	}
	if _, ok := secrets["A"]; ok {
		t.Error("ApplyPreset mutated input secrets")
	}
}
