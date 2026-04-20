package sync

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envfile"
)

func TestApplyTransforms_NilConfig(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	out, err := ApplyTransforms(secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" {
		t.Errorf("expected passthrough, got %q", out["FOO"])
	}
}

func TestApplyTransforms_EmptyRules(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	out, err := ApplyTransforms(secrets, &TransformConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" {
		t.Errorf("expected passthrough, got %q", out["FOO"])
	}
}

func TestApplyTransforms_AppliesRule(t *testing.T) {
	secrets := map[string]string{"DB_URL": "  postgres://localhost  "}
	cfg := &TransformConfig{
		Rules: []envfile.TransformRule{
			{KeyPrefix: "DB_", Transform: envfile.TrimSpace},
		},
	}
	out, err := ApplyTransforms(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_URL"] != "postgres://localhost" {
		t.Errorf("expected trimmed value, got %q", out["DB_URL"])
	}
}

func TestApplyTransforms_ErrorPropagates(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	cfg := &TransformConfig{
		Rules: []envfile.TransformRule{
			{
				KeyPrefix: "KEY",
				Transform: func(k, v string) (string, error) {
					return "", fmt.Errorf("intentional failure")
				},
			},
		},
	}
	_, err := ApplyTransforms(secrets, cfg)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
