package envfile

import "testing"

func TestRedact_MasksSensitiveKeys(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_TOKEN":   "abc123",
		"APP_NAME":    "myapp",
	}
	got := Redact(secrets, nil)
	if got["DB_PASSWORD"] != "***" {
		t.Errorf("expected DB_PASSWORD masked, got %q", got["DB_PASSWORD"])
	}
	if got["API_TOKEN"] != "***" {
		t.Errorf("expected API_TOKEN masked, got %q", got["API_TOKEN"])
	}
	if got["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME unmasked, got %q", got["APP_NAME"])
	}
}

func TestRedact_CustomMask(t *testing.T) {
	secrets := map[string]string{"SECRET_KEY": "value"}
	got := Redact(secrets, &RedactOptions{Mask: "REDACTED"})
	if got["SECRET_KEY"] != "REDACTED" {
		t.Errorf("expected REDACTED, got %q", got["SECRET_KEY"])
	}
}

func TestRedact_CustomSubstrings(t *testing.T) {
	secrets := map[string]string{
		"DB_PASS":  "secret",
		"MY_TOKEN": "tok",
		"REGION":   "us-east-1",
	}
	got := Redact(secrets, &RedactOptions{SensitiveSubstrings: []string{"pass"}})
	if got["DB_PASS"] != "***" {
		t.Errorf("expected DB_PASS masked")
	}
	if got["MY_TOKEN"] != "tok" {
		t.Errorf("expected MY_TOKEN unmasked")
	}
}

func TestRedact_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"API_KEY": "original"}
	Redact(secrets, nil)
	if secrets["API_KEY"] != "original" {
		t.Error("input map was mutated")
	}
}

func TestRedact_EmptyMap(t *testing.T) {
	got := Redact(map[string]string{}, nil)
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestRedact_CaseInsensitiveKey(t *testing.T) {
	secrets := map[string]string{"db_password": "s3cr3t"}
	got := Redact(secrets, nil)
	if got["db_password"] != "***" {
		t.Errorf("expected lowercase key to be masked")
	}
}
