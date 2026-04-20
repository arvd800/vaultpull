package envfile

import (
	"errors"
	"testing"
)

func TestTransform_NoRules(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := Transform(secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("expected passthrough, got %v", out)
	}
}

func TestTransform_UpperCasePrefix(t *testing.T) {
	secrets := map[string]string{"APP_NAME": "hello", "DB_HOST": "localhost"}
	rules := []TransformRule{
		{KeyPrefix: "APP_", Transform: UpperCase},
	}
	out, err := Transform(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_NAME"] != "HELLO" {
		t.Errorf("expected HELLO, got %q", out["APP_NAME"])
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected passthrough for DB_HOST, got %q", out["DB_HOST"])
	}
}

func TestTransform_TrimSpace(t *testing.T) {
	secrets := map[string]string{"TOKEN": "  abc123  "}
	rules := []TransformRule{
		{KeyPrefix: "", Transform: TrimSpace},
	}
	out, err := Transform(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TOKEN"] != "abc123" {
		t.Errorf("expected trimmed value, got %q", out["TOKEN"])
	}
}

func TestTransform_ErrorPropagates(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	failFn := func(k, v string) (string, error) {
		return "", errors.New("boom")
	}
	rules := []TransformRule{{KeyPrefix: "KEY", Transform: failFn}}
	_, err := Transform(secrets, rules)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestTransform_EmptySecrets(t *testing.T) {
	out, err := Transform(map[string]string{}, []TransformRule{
		{KeyPrefix: "", Transform: UpperCase},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestTransform_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	_, err := Transform(secrets, []TransformRule{{KeyPrefix: "", Transform: UpperCase}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["FOO"] != "bar" {
		t.Errorf("input was mutated")
	}
}

func TestTransform_FirstMatchingRuleWins(t *testing.T) {
	secrets := map[string]string{"APP_KEY": "hello"}
	rules := []TransformRule{
		{KeyPrefix: "APP_", Transform: UpperCase},
		{KeyPrefix: "APP_", Transform: LowerCase},
	}
	out, err := Transform(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_KEY"] != "HELLO" {
		t.Errorf("expected first rule to win, got %q", out["APP_KEY"])
	}
}
