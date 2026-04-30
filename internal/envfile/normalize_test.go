package envfile

import (
	"testing"
)

func TestNormalize_UpperKeys(t *testing.T) {
	input := map[string]string{"db_host": "localhost", "api_key": "secret"}
	opts := NormalizeOptions{UpperKeys: true}
	out, err := Normalize(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
	if out["API_KEY"] != "secret" {
		t.Errorf("expected API_KEY=secret, got %q", out["API_KEY"])
	}
}

func TestNormalize_ReplaceHyphens(t *testing.T) {
	input := map[string]string{"my-key": "value"}
	opts := NormalizeOptions{ReplaceHyphens: true}
	out, err := Normalize(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["my_key"]; !ok {
		t.Errorf("expected key my_key to exist, got keys: %v", out)
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	input := map[string]string{"KEY": "  hello world  "}
	opts := NormalizeOptions{TrimValues: true}
	out, err := Normalize(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "hello world" {
		t.Errorf("expected trimmed value, got %q", out["KEY"])
	}
}

func TestNormalize_CollapseUnderscores(t *testing.T) {
	input := map[string]string{"MY__KEY": "val"}
	opts := NormalizeOptions{CollapseUnderscores: true}
	out, err := Normalize(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["MY_KEY"]; !ok {
		t.Errorf("expected MY_KEY, got: %v", out)
	}
}

func TestNormalize_NilInput(t *testing.T) {
	out, err := Normalize(nil, DefaultNormalizeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestNormalize_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"my-key": "  val  "}
	opts := DefaultNormalizeOptions()
	_, err := Normalize(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input["my-key"] != "  val  " {
		t.Errorf("input was mutated: %v", input)
	}
}

func TestNormalize_DefaultOptions(t *testing.T) {
	input := map[string]string{"my-db-host": "  localhost  "}
	out, err := Normalize(input, DefaultNormalizeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["MY_DB_HOST"] != "localhost" {
		t.Errorf("expected MY_DB_HOST=localhost, got %v", out)
	}
}
