package envfile

import (
	"fmt"
	"strings"
)

// LintRule describes a single lint check.
type LintRule struct {
	Code    string
	Message string
}

// LintResult holds a single lint finding.
type LintResult struct {
	Key  string
	Rule LintRule
}

func (r LintResult) String() string {
	return fmt.Sprintf("[%s] %s: %s", r.Rule.Code, r.Key, r.Rule.Message)
}

var builtinRules = []struct {
	code    string
	message string
	check   func(k, v string) bool
}{
	{
		code:    "EMPTY_VALUE",
		message: "key has an empty value",
		check:   func(k, v string) bool { return strings.TrimSpace(v) == "" },
	},
	{
		code:    "LOWERCASE_KEY",
		message: "key contains lowercase letters (prefer UPPER_SNAKE_CASE)",
		check:   func(k, v string) bool { return k != strings.ToUpper(k) },
	},
	{
		code:    "LEADING_TRAILING_SPACE",
		message: "value has leading or trailing whitespace",
		check:   func(k, v string) bool { return v != strings.TrimSpace(v) && strings.TrimSpace(v) != "" },
	},
	{
		code:    "DOUBLE_UNDERSCORE",
		message: "key contains consecutive underscores",
		check:   func(k, v string) bool { return strings.Contains(k, "__") },
	},
}

// Lint runs all built-in lint rules against the provided secrets map.
// It returns a slice of LintResult for every finding.
func Lint(secrets map[string]string) []LintResult {
	var results []LintResult
	for k, v := range secrets {
		for _, rule := range builtinRules {
			if rule.check(k, v) {
				results = append(results, LintResult{
					Key: k,
					Rule: LintRule{Code: rule.code, Message: rule.message},
				})
			}
		}
	}
	return results
}

// FormatLintResults returns a human-readable summary of lint findings.
func FormatLintResults(results []LintResult) string {
	if len(results) == 0 {
		return "No lint issues found."
	}
	var sb strings.Builder
	for _, r := range results {
		sb.WriteString(r.String())
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}
