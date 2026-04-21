package envfile

import (
	"os"
	"testing"
)

func TestResolve_BasicInterpolation(t *testing.T) {
	secrets := map[string]string{
		"BASE_URL": "https://example.com",
		"API_URL":  "${BASE_URL}/api",
	}
	out, err := Resolve(secrets, ResolveOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["API_URL"]; got != "https://example.com/api" {
		t.Errorf("API_URL = %q, want %q", got, "https://example.com/api")
	}
}

func TestResolve_NoReferences(t *testing.T) {
	secrets := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}
	out, err := Resolve(secrets, ResolveOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("values should be unchanged: %v", out)
	}
}

func TestResolve_MissingKey_Error(t *testing.T) {
	secrets := map[string]string{
		"URL": "${MISSING}/path",
	}
	_, err := Resolve(secrets, ResolveOptions{})
	if err == nil {
		t.Fatal("expected error for undefined variable, got nil")
	}
}

func TestResolve_AllowMissing(t *testing.T) {
	secrets := map[string]string{
		"URL": "${MISSING}/path",
	}
	out, err := Resolve(secrets, ResolveOptions{AllowMissing: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["URL"]; got != "${MISSING}/path" {
		t.Errorf("URL = %q, want original value preserved", got)
	}
}

func TestResolve_FallbackToEnv(t *testing.T) {
	t.Setenv("MY_HOST", "localhost")
	secrets := map[string]string{
		"DSN": "postgres://${MY_HOST}/db",
	}
	out, err := Resolve(secrets, ResolveOptions{FallbackToEnv: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["DSN"]; got != "postgres://localhost/db" {
		t.Errorf("DSN = %q, want %q", got, "postgres://localhost/db")
	}
}

func TestResolve_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{
		"A": "hello",
		"B": "${A} world",
	}
	orig := map[string]string{"A": "hello", "B": "${A} world"}
	_, _ = Resolve(secrets, ResolveOptions{})
	for k, v := range orig {
		if secrets[k] != v {
			t.Errorf("input mutated: secrets[%q] = %q, want %q", k, secrets[k], v)
		}
	}
	_ = os.Getenv // suppress unused import
}
