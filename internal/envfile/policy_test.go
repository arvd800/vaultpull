package envfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/envfile"
)

func TestSaveAndLoadPolicy_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.json")
	p := envfile.Policy{
		Rules: []envfile.PolicyRule{
			{Key: "DB_PASSWORD", Required: true},
			{Key: "DEBUG", Pattern: `^(true|false)$`},
		},
	}
	if err := envfile.SavePolicy(path, p); err != nil {
		t.Fatalf("SavePolicy: %v", err)
	}
	loaded, err := envfile.LoadPolicy(path)
	if err != nil {
		t.Fatalf("LoadPolicy: %v", err)
	}
	if len(loaded.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(loaded.Rules))
	}
}

func TestLoadPolicy_NonExistent(t *testing.T) {
	p, err := envfile.LoadPolicy("/nonexistent/policy.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Rules) != 0 {
		t.Errorf("expected empty policy")
	}
}

func TestSavePolicy_EmptyPath(t *testing.T) {
	err := envfile.SavePolicy("", envfile.Policy{})
	if err != nil {
		t.Errorf("expected no error for empty path, got %v", err)
	}
}

func TestSavePolicy_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.json")
	if err := envfile.SavePolicy(path, envfile.Policy{}); err != nil {
		t.Fatalf("SavePolicy: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
}

func TestEnforcePolicy_RequiredMissing(t *testing.T) {
	p := envfile.Policy{
		Rules: []envfile.PolicyRule{{Key: "DB_PASSWORD", Required: true}},
	}
	violations := envfile.EnforcePolicy(map[string]string{}, p)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "DB_PASSWORD" {
		t.Errorf("unexpected key: %s", violations[0].Key)
	}
}

func TestEnforcePolicy_DeniedKey(t *testing.T) {
	p := envfile.Policy{
		Rules: []envfile.PolicyRule{{Key: "SECRET", Deny: true}},
	}
	violations := envfile.EnforcePolicy(map[string]string{"SECRET": "val"}, p)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestEnforcePolicy_PatternMismatch(t *testing.T) {
	p := envfile.Policy{
		Rules: []envfile.PolicyRule{{Key: "DEBUG", Pattern: `^(true|false)$`}},
	}
	violations := envfile.EnforcePolicy(map[string]string{"DEBUG": "yes"}, p)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestEnforcePolicy_NoViolations(t *testing.T) {
	p := envfile.Policy{
		Rules: []envfile.PolicyRule{
			{Key: "DEBUG", Pattern: `^(true|false)$`, Required: true},
		},
	}
	violations := envfile.EnforcePolicy(map[string]string{"DEBUG": "true"}, p)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}
