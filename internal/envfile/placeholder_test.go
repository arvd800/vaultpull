package envfile

import (
	"testing"
)

func TestResolvePlaceholders_Basic(t *testing.T) {
	dst := map[string]string{
		"DB_PASS": "PLACEHOLDER:REAL_DB_PASS",
		"APP_KEY": "static-value",
	}
	src := map[string]string{
		"REAL_DB_PASS": "s3cr3t",
	}
	out, err := ResolvePlaceholders(dst, src, PlaceholderConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASS"] != "s3cr3t" {
		t.Errorf("expected s3cr3t, got %q", out["DB_PASS"])
	}
	if out["APP_KEY"] != "static-value" {
		t.Errorf("expected static-value, got %q", out["APP_KEY"])
	}
}

func TestResolvePlaceholders_CustomPrefix(t *testing.T) {
	dst := map[string]string{"TOKEN": "FILL:MY_TOKEN"}
	src := map[string]string{"MY_TOKEN": "abc123"}
	out, err := ResolvePlaceholders(dst, src, PlaceholderConfig{Prefix: "FILL:"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TOKEN"] != "abc123" {
		t.Errorf("expected abc123, got %q", out["TOKEN"])
	}
}

func TestResolvePlaceholders_FailOnUnresolved(t *testing.T) {
	dst := map[string]string{"SECRET": "PLACEHOLDER:MISSING_KEY"}
	src := map[string]string{}
	_, err := ResolvePlaceholders(dst, src, PlaceholderConfig{FailOnUnresolved: true})
	if err == nil {
		t.Fatal("expected error for unresolved placeholder, got nil")
	}
}

func TestResolvePlaceholders_UnresolvedSilent(t *testing.T) {
	dst := map[string]string{"SECRET": "PLACEHOLDER:MISSING_KEY"}
	src := map[string]string{}
	out, err := ResolvePlaceholders(dst, src, PlaceholderConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// value should remain unchanged
	if out["SECRET"] != "PLACEHOLDER:MISSING_KEY" {
		t.Errorf("expected original placeholder, got %q", out["SECRET"])
	}
}

func TestResolvePlaceholders_DoesNotMutateInput(t *testing.T) {
	dst := map[string]string{"X": "PLACEHOLDER:Y"}
	src := map[string]string{"Y": "resolved"}
	_, _ = ResolvePlaceholders(dst, src, PlaceholderConfig{})
	if dst["X"] != "PLACEHOLDER:Y" {
		t.Error("ResolvePlaceholders mutated the input dst map")
	}
}

func TestListPlaceholders_ReturnsKeys(t *testing.T) {
	m := map[string]string{
		"A": "PLACEHOLDER:SRC_A",
		"B": "real-value",
		"C": "PLACEHOLDER:SRC_C",
	}
	keys := ListPlaceholders(m, "")
	if len(keys) != 2 {
		t.Fatalf("expected 2 placeholder keys, got %d", len(keys))
	}
	if keys[0] != "A" || keys[1] != "C" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestListPlaceholders_EmptyMap(t *testing.T) {
	keys := ListPlaceholders(map[string]string{}, "")
	if len(keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(keys))
	}
}
