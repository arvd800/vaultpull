package sync_test

import (
	"bytes"
	"testing"

	"github.com/yourusername/vaultpull/internal/envfile"
	"github.com/yourusername/vaultpull/internal/sync"
)

func TestApplyPreset_NoPresetName(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	out, err := sync.ApplyPreset("", t.TempDir(), secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" {
		t.Errorf("expected A=1, got %q", out["A"])
	}
}

func TestApplyPreset_AppliesOverrides(t *testing.T) {
	dir := t.TempDir()
	presets := []envfile.Preset{
		{Name: "ci", Values: map[string]string{"ENV": "ci", "CI": "true"}},
	}
	if err := envfile.SavePresets(dir, presets); err != nil {
		t.Fatalf("SavePresets: %v", err)
	}
	secrets := map[string]string{"ENV": "dev", "DB": "localhost"}
	var buf bytes.Buffer
	out, err := sync.ApplyPreset("ci", dir, secrets, &buf)
	if err != nil {
		t.Fatalf("ApplyPreset: %v", err)
	}
	if out["ENV"] != "ci" {
		t.Errorf("expected ENV=ci, got %q", out["ENV"])
	}
	if out["CI"] != "true" {
		t.Errorf("expected CI=true, got %q", out["CI"])
	}
	if out["DB"] != "localhost" {
		t.Errorf("expected DB preserved")
	}
	if buf.Len() == 0 {
		t.Error("expected log output")
	}
}

func TestApplyPreset_PresetNotFound_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	_, err := sync.ApplyPreset("ghost", dir, map[string]string{}, nil)
	if err == nil {
		t.Error("expected error for missing preset")
	}
}

func TestApplyPreset_NilOutput_DoesNotPanic(t *testing.T) {
	dir := t.TempDir()
	presets := []envfile.Preset{
		{Name: "x", Values: map[string]string{"K": "V"}},
	}
	if err := envfile.SavePresets(dir, presets); err != nil {
		t.Fatalf("SavePresets: %v", err)
	}
	_, err := sync.ApplyPreset("x", dir, map[string]string{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
