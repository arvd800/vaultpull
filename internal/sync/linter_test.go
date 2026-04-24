package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunLint_NoIssues(t *testing.T) {
	secrets := map[string]string{
		"MY_KEY": "value",
	}
	var buf bytes.Buffer
	err := RunLint(secrets, LintConfig{Output: &buf, FailOnWarnings: true})
	if err != nil {
		t.Errorf("expected no error for clean secrets, got: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got: %q", buf.String())
	}
}

func TestRunLint_WritesFindings(t *testing.T) {
	secrets := map[string]string{
		"bad_key": "",
	}
	var buf bytes.Buffer
	_ = RunLint(secrets, LintConfig{Output: &buf})
	out := buf.String()
	if !strings.Contains(out, "Lint warnings:") {
		t.Errorf("expected header in output, got: %q", out)
	}
}

func TestRunLint_FailOnWarnings_ReturnsError(t *testing.T) {
	secrets := map[string]string{
		"bad_key": "",
	}
	var buf bytes.Buffer
	err := RunLint(secrets, LintConfig{Output: &buf, FailOnWarnings: true})
	if err == nil {
		t.Error("expected error when FailOnWarnings is true and findings exist")
	}
}

func TestRunLint_NoFailOnWarnings_NoError(t *testing.T) {
	secrets := map[string]string{
		"bad_key": "",
	}
	var buf bytes.Buffer
	err := RunLint(secrets, LintConfig{Output: &buf, FailOnWarnings: false})
	if err != nil {
		t.Errorf("expected no error when FailOnWarnings is false, got: %v", err)
	}
}

func TestRunLint_DefaultOutput_DoesNotPanic(t *testing.T) {
	secrets := map[string]string{"CLEAN": "value"}
	// nil output should default to stderr without panicking
	_ = RunLint(secrets, LintConfig{})
}
