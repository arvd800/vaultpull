package envfile

import (
	"fmt"
	"strings"
)

// MaskRule defines how a specific key's value should be masked in output.
type MaskRule struct {
	// ShowFirst is the number of characters to reveal at the start.
	ShowFirst int
	// ShowLast is the number of characters to reveal at the end.
	ShowLast int
	// Char is the masking character (defaults to '*').
	Char string
}

// MaskResult holds the original key and its masked value.
type MaskResult struct {
	Key    string
	Masked string
}

// DefaultMaskRule is used when no specific rule is provided.
var DefaultMaskRule = MaskRule{ShowFirst: 0, ShowLast: 0, Char: "*"}

// Mask applies masking rules to a map of secrets and returns masked results.
// Keys not present in rules use DefaultMaskRule.
func Mask(secrets map[string]string, rules map[string]MaskRule) []MaskResult {
	results := make([]MaskResult, 0, len(secrets))
	for _, k := range sorted(secrets) {
		v := secrets[k]
		rule, ok := rules[k]
		if !ok {
			rule = DefaultMaskRule
		}
		results = append(results, MaskResult{
			Key:    k,
			Masked: applyMask(v, rule),
		})
	}
	return results
}

// applyMask masks a single value according to the given rule.
func applyMask(value string, rule MaskRule) string {
	char := rule.Char
	if char == "" {
		char = "*"
	}
	n := len(value)
	show := rule.ShowFirst + rule.ShowLast
	if show >= n {
		// Reveal nothing if the value is too short to mask meaningfully.
		return strings.Repeat(char, n)
	}
	prefix := value[:rule.ShowFirst]
	suffix := ""
	if rule.ShowLast > 0 {
		suffix = value[n-rule.ShowLast:]
	}
	midLen := n - rule.ShowFirst - rule.ShowLast
	return fmt.Sprintf("%s%s%s", prefix, strings.Repeat(char, midLen), suffix)
}
