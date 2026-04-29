package envfile

import (
	"testing"
)

func TestSanitize_TrimSpace(t *testing.T) {
	secrets := map[string]string{
		"KEY": "  hello world  ",
	}
	opts := DefaultSanitizeOptions()
	out := Sanitize(secrets, opts)
	if got := out["KEY"]; got != "hello world" {
		t.Errorf("expected 'hello world', got %q", got)
	}
}

func TestSanitize_RemoveNonPrintable(t *testing.T) {
	secrets := map[string]string{
		"KEY": "val\x00ue\x1f",
	}
	opts := DefaultSanitizeOptions()
	out := Sanitize(secrets, opts)
	if got := out["KEY"]; got != "value" {
		t.Errorf("expected 'value', got %q", got)
	}
}

func TestSanitize_NormalizeKeys(t *testing.T) {
	secrets := map[string]string{
		"my-key":    "v1",
		"other key": "v2",
	}
	opts := SanitizeOptions{NormalizeKeys: true}
	out := Sanitize(secrets, opts)
	if _, ok := out["MY_KEY"]; !ok {
		t.Error("expected MY_KEY to be present")
	}
	if _, ok := out["OTHER_KEY"]; !ok {
		t.Error("expected OTHER_KEY to be present")
	}
}

func TestSanitize_DropEmpty(t *testing.T) {
	secrets := map[string]string{
		"PRESENT": "value",
		"EMPTY":   "",
		"SPACES":  "   ",
	}
	opts := SanitizeOptions{TrimSpace: true, DropEmpty: true}
	out := Sanitize(secrets, opts)
	if _, ok := out["EMPTY"]; ok {
		t.Error("expected EMPTY to be dropped")
	}
	if _, ok := out["SPACES"]; ok {
		t.Error("expected SPACES to be dropped after trimming")
	}
	if _, ok := out["PRESENT"]; !ok {
		t.Error("expected PRESENT to be retained")
	}
}

func TestSanitize_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{
		"KEY": "  spaced  ",
	}
	opts := DefaultSanitizeOptions()
	Sanitize(secrets, opts)
	if secrets["KEY"] != "  spaced  " {
		t.Error("input map was mutated")
	}
}

func TestSanitize_EmptyMap(t *testing.T) {
	out := Sanitize(map[string]string{}, DefaultSanitizeOptions())
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}

func TestDefaultSanitizeOptions(t *testing.T) {
	opts := DefaultSanitizeOptions()
	if !opts.TrimSpace {
		t.Error("expected TrimSpace to be true by default")
	}
	if !opts.RemoveNonPrintable {
		t.Error("expected RemoveNonPrintable to be true by default")
	}
	if opts.NormalizeKeys {
		t.Error("expected NormalizeKeys to be false by default")
	}
	if opts.DropEmpty {
		t.Error("expected DropEmpty to be false by default")
	}
}
