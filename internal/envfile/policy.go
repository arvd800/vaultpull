package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

// PolicyRule defines a rule applied to a secret key.
type PolicyRule struct {
	Key     string `json:"key"`
	Pattern string `json:"pattern"`
	Required bool   `json:"required"`
	Deny    bool   `json:"deny"`
}

// Policy holds a collection of rules.
type Policy struct {
	Rules []PolicyRule `json:"rules"`
}

// PolicyViolation describes a single policy violation.
type PolicyViolation struct {
	Key     string
	Message string
}

func (v PolicyViolation) Error() string {
	return fmt.Sprintf("policy violation [%s]: %s", v.Key, v.Message)
}

// SavePolicy writes a Policy to a JSON file.
func SavePolicy(path string, p Policy) error {
	if path == "" {
		return nil
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadPolicy reads a Policy from a JSON file.
func LoadPolicy(path string) (Policy, error) {
	var p Policy
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return p, nil
	}
	if err != nil {
		return p, err
	}
	return p, json.Unmarshal(data, &p)
}

// EnforcePolicy checks secrets against the policy and returns all violations.
func EnforcePolicy(secrets map[string]string, p Policy) []PolicyViolation {
	var violations []PolicyViolation
	for _, rule := range p.Rules {
		val, exists := secrets[rule.Key]
		if rule.Deny && exists {
			violations = append(violations, PolicyViolation{Key: rule.Key, Message: "key is denied by policy"})
			continue
		}
		if rule.Required && !exists {
			violations = append(violations, PolicyViolation{Key: rule.Key, Message: "required key is missing"})
			continue
		}
		if rule.Pattern != "" && exists {
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				violations = append(violations, PolicyViolation{Key: rule.Key, Message: fmt.Sprintf("invalid pattern: %v", err)})
				continue
			}
			if !re.MatchString(val) {
				violations = append(violations, PolicyViolation{Key: rule.Key, Message: fmt.Sprintf("value does not match pattern %q", rule.Pattern)})
			}
		}
	}
	return violations
}
