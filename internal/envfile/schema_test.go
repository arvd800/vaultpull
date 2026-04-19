package envfile

import (
	"os"
	"testing"
)

func TestValidateSchema_AllPresent(t *testing.T) {
	schema := Schema{
		"DATABASE_URL": {Required: true},
		"PORT":         {Required: true, Pattern: `^[0-9]+$`},
	}
	secrets := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"PORT":         "5432",
	}
	errs := ValidateSchema(secrets, schema)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidateSchema_MissingRequired(t *testing.T) {
	schema := Schema{
		"DATABASE_URL": {Required: true},
	}
	errs := ValidateSchema(map[string]string{}, schema)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestValidateSchema_PatternMismatch(t *testing.T) {
	schema := Schema{
		"PORT": {Required: false, Pattern: `^[0-9]+$`},
	}
	secrets := map[string]string{"PORT": "not-a-number"}
	errs := ValidateSchema(secrets, schema)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestValidateSchema_EmptySchema(t *testing.T) {
	errs := ValidateSchema(map[string]string{"FOO": "bar"}, Schema{})
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestLoadSchema_RoundTrip(t *testing.T) {
	content := `rules:
  DATABASE_URL:
    required: true
  PORT:
    required: false
    pattern: '^[0-9]+$'
`
	f, err := os.CreateTemp(t.TempDir(), "schema*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()

	schema, err := LoadSchema(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := schema["DATABASE_URL"]; !ok {
		t.Error("expected DATABASE_URL in schema")
	}
	if schema["PORT"].Pattern != `^[0-9]+$` {
		t.Errorf("unexpected pattern: %q", schema["PORT"].Pattern)
	}
}
