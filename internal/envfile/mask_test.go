package envfile

import (
	"testing"
)

func TestMask_DefaultRule_HidesAll(t *testing.T) {
	secrets := map[string]string{
		"API_KEY": "supersecret",
	}
	results := Mask(secrets, nil)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Masked != "***********" {
		t.Errorf("expected all stars, got %q", results[0].Masked)
	}
}

func TestMask_ShowFirst(t *testing.T) {
	secrets := map[string]string{"TOKEN": "abcdef"}
	rules := map[string]MaskRule{
		"TOKEN": {ShowFirst: 2, Char: "*"},
	}
	results := Mask(secrets, rules)
	if results[0].Masked != "ab****" {
		t.Errorf("unexpected masked value: %q", results[0].Masked)
	}
}

func TestMask_ShowLast(t *testing.T) {
	secrets := map[string]string{"TOKEN": "abcdef"}
	rules := map[string]MaskRule{
		"TOKEN": {ShowLast: 2, Char: "#"},
	}
	results := Mask(secrets, rules)
	if results[0].Masked != "####ef" {
		t.Errorf("unexpected masked value: %q", results[0].Masked)
	}
}

func TestMask_ShowFirstAndLast(t *testing.T) {
	secrets := map[string]string{"SECRET": "hello_world"}
	rules := map[string]MaskRule{
		"SECRET": {ShowFirst: 2, ShowLast: 2, Char: "-"},
	}
	results := Mask(secrets, rules)
	// "he-------ld" — 11 chars, 2+2 shown, 7 masked
	if results[0].Masked != "he-------ld" {
		t.Errorf("unexpected masked value: %q", results[0].Masked)
	}
}

func TestMask_ShortValue_AllMasked(t *testing.T) {
	secrets := map[string]string{"K": "ab"}
	rules := map[string]MaskRule{
		"K": {ShowFirst: 3, ShowLast: 3, Char: "*"},
	}
	results := Mask(secrets, rules)
	// show >= len, so mask everything
	if results[0].Masked != "**" {
		t.Errorf("expected all masked for short value, got %q", results[0].Masked)
	}
}

func TestMask_DefaultCharUsedWhenEmpty(t *testing.T) {
	secrets := map[string]string{"X": "123"}
	rules := map[string]MaskRule{
		"X": {ShowFirst: 0, ShowLast: 0, Char: ""},
	}
	results := Mask(secrets, rules)
	if results[0].Masked != "***" {
		t.Errorf("expected default '*' char, got %q", results[0].Masked)
	}
}

func TestMask_SortedOutput(t *testing.T) {
	secrets := map[string]string{
		"Z_KEY": "zzz",
		"A_KEY": "aaa",
	}
	results := Mask(secrets, nil)
	if results[0].Key != "A_KEY" || results[1].Key != "Z_KEY" {
		t.Errorf("expected sorted output, got %v, %v", results[0].Key, results[1].Key)
	}
}
