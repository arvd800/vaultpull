package envfile

import (
	"strings"
	"testing"
)

func TestFormat_SortsKeys(t *testing.T) {
	secrets := map[string]string{
		"ZEBRA": "last",
		"ALPHA": "first",
		"MIDDLE": "mid",
	}
	out := Format(secrets)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "ALPHA=") {
		t.Errorf("expected first line to start with ALPHA=, got %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZEBRA=") {
		t.Errorf("expected last line to start with ZEBRA=, got %s", lines[2])
	}
}

func TestFormat_QuotesValuesWithSpaces(t *testing.T) {
	secrets := map[string]string{
		"KEY": "hello world",
	}
	out := Format(secrets)
	if !strings.Contains(out, `KEY="hello world"`) {
		t.Errorf("expected quoted value, got: %s", out)
	}
}

func TestFormat_EmptyMap(t *testing.T) {
	out := Format(map[string]string{})
	if out != "" {
		t.Errorf("expected empty string, got: %q", out)
	}
}

func TestParseLine_Basic(t *testing.T) {
	k, v, ok := ParseLine("FOO=bar")
	if !ok || k != "FOO" || v != "bar" {
		t.Errorf("unexpected result: %q %q %v", k, v, ok)
	}
}

func TestParseLine_QuotedValue(t *testing.T) {
	k, v, ok := ParseLine(`BAR="hello world"`)
	if !ok || k != "BAR" || v != "hello world" {
		t.Errorf("unexpected result: %q %q %v", k, v, ok)
	}
}

func TestParseLine_Comment(t *testing.T) {
	_, _, ok := ParseLine("# this is a comment")
	if ok {
		t.Error("expected comment line to return ok=false")
	}
}

func TestParseLine_Empty(t *testing.T) {
	_, _, ok := ParseLine("")
	if ok {
		t.Error("expected empty line to return ok=false")
	}
}

func TestParseLine_NoEquals(t *testing.T) {
	_, _, ok := ParseLine("NOEQUALS")
	if ok {
		t.Error("expected line without '=' to return ok=false")
	}
}
