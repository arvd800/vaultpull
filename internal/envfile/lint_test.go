package envfile

import (
	"strings"
	"testing"
)

func TestLint_EmptyValue(t *testing.T) {
	secrets := map[string]string{
		"MY_KEY": "",
	}
	results := Lint(secrets)
	if !containsCode(results, "EMPTY_VALUE") {
		t.Error("expected EMPTY_VALUE lint finding")
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	secrets := map[string]string{
		"my_key": "value",
	}
	results := Lint(secrets)
	if !containsCode(results, "LOWERCASE_KEY") {
		t.Error("expected LOWERCASE_KEY lint finding")
	}
}

func TestLint_LeadingTrailingSpace(t *testing.T) {
	secrets := map[string]string{
		"MY_KEY": "  value  ",
	}
	results := Lint(secrets)
	if !containsCode(results, "LEADING_TRAILING_SPACE") {
		t.Error("expected LEADING_TRAILING_SPACE lint finding")
	}
}

func TestLint_DoubleUnderscore(t *testing.T) {
	secrets := map[string]string{
		"MY__KEY": "value",
	}
	results := Lint(secrets)
	if !containsCode(results, "DOUBLE_UNDERSCORE") {
		t.Error("expected DOUBLE_UNDERSCORE lint finding")
	}
}

func TestLint_CleanMap(t *testing.T) {
	secrets := map[string]string{
		"MY_KEY":    "value",
		"OTHER_KEY": "123",
	}
	results := Lint(secrets)
	if len(results) != 0 {
		t.Errorf("expected no lint findings, got %d", len(results))
	}
}

func TestFormatLintResults_NoFindings(t *testing.T) {
	out := FormatLintResults(nil)
	if out != "No lint issues found." {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatLintResults_WithFindings(t *testing.T) {
	results := []LintResult{
		{Key: "bad_key", Rule: LintRule{Code: "LOWERCASE_KEY", Message: "key contains lowercase letters"}},
	}
	out := FormatLintResults(results)
	if !strings.Contains(out, "LOWERCASE_KEY") {
		t.Errorf("expected LOWERCASE_KEY in output, got: %q", out)
	}
	if !strings.Contains(out, "bad_key") {
		t.Errorf("expected key name in output, got: %q", out)
	}
}

func TestLint_EmptyMap(t *testing.T) {
	results := Lint(map[string]string{})
	if len(results) != 0 {
		t.Errorf("expected no results for empty map, got %d", len(results))
	}
}

func containsCode(results []LintResult, code string) bool {
	for _, r := range results {
		if r.Rule.Code == code {
			return true
		}
	}
	return false
}
