package sync

import (
	"fmt"
	"strings"

	"github.com/yourusername/vaultpull/internal/envfile"
)

// EnforceConfig holds configuration for policy enforcement during sync.
type EnforceConfig struct {
	// PolicyPath is the path to the policy JSON file. Empty means no enforcement.
	PolicyPath string
	// FailOnViolation causes Run to return an error if any violations are found.
	FailOnViolation bool
}

// ApplyPolicyEnforcement loads a policy and checks secrets against it.
// It returns a formatted warning string and an error if FailOnViolation is set.
func ApplyPolicyEnforcement(secrets map[string]string, cfg EnforceConfig) (string, error) {
	if cfg.PolicyPath == "" {
		return "", nil
	}
	policy, err := envfile.LoadPolicy(cfg.PolicyPath)
	if err != nil {
		return "", fmt.Errorf("load policy: %w", err)
	}
	violations := envfile.EnforcePolicy(secrets, policy)
	if len(violations) == 0 {
		return "", nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("policy enforcement: %d violation(s):\n", len(violations)))
	for _, v := range violations {
		sb.WriteString(fmt.Sprintf("  - %s\n", v.Error()))
	}
	warning := strings.TrimRight(sb.String(), "\n")
	if cfg.FailOnViolation {
		return warning, fmt.Errorf("%s", warning)
	}
	return warning, nil
}
