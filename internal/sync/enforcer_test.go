package sync_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/envfile"
	"github.com/yourusername/vaultpull/internal/sync"
)

func writePolicy(t *testing.T, dir string, p envfile.Policy) string {
	t.Helper()
	path := filepath.Join(dir, "policy.json")
	data, _ := json.MarshalIndent(p, "", "  ")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("writePolicy: %v", err)
	}
	return path
}

func TestApplyPolicyEnforcement_NoPath(t *testing.T) {
	warn, err := sync.ApplyPolicyEnforcement(map[string]string{"K": "v"}, sync.EnforceConfig{})
	if err != nil || warn != "" {
		t.Errorf("expected no-op, got warn=%q err=%v", warn, err)
	}
}

func TestApplyPolicyEnforcement_NoViolations(t *testing.T) {
	dir := t.TempDir()
	path := writePolicy(t, dir, envfile.Policy{
		Rules: []envfile.PolicyRule{{Key: "API_KEY", Required: true}},
	})
	warn, err := sync.ApplyPolicyEnforcement(
		map[string]string{"API_KEY": "abc"},
		sync.EnforceConfig{PolicyPath: path},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if warn != "" {
		t.Errorf("expected no warning, got %q", warn)
	}
}

func TestApplyPolicyEnforcement_WarnsOnViolation(t *testing.T) {
	dir := t.TempDir()
	path := writePolicy(t, dir, envfile.Policy{
		Rules: []envfile.PolicyRule{{Key: "API_KEY", Required: true}},
	})
	warn, err := sync.ApplyPolicyEnforcement(
		map[string]string{},
		sync.EnforceConfig{PolicyPath: path, FailOnViolation: false},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(warn, "API_KEY") {
		t.Errorf("expected warning to mention API_KEY, got %q", warn)
	}
}

func TestApplyPolicyEnforcement_FailOnViolation(t *testing.T) {
	dir := t.TempDir()
	path := writePolicy(t, dir, envfile.Policy{
		Rules: []envfile.PolicyRule{{Key: "DB_PASS", Required: true}},
	})
	_, err := sync.ApplyPolicyEnforcement(
		map[string]string{},
		sync.EnforceConfig{PolicyPath: path, FailOnViolation: true},
	)
	if err == nil {
		t.Fatal("expected error due to FailOnViolation, got nil")
	}
	if !strings.Contains(err.Error(), "DB_PASS") {
		t.Errorf("expected error to mention DB_PASS, got %v", err)
	}
}
