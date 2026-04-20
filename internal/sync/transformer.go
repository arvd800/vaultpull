package sync

import (
	"fmt"

	"github.com/your-org/vaultpull/internal/envfile"
)

// TransformConfig holds the configuration for secret transformations.
type TransformConfig struct {
	// Rules is an ordered list of transform rules to apply.
	Rules []envfile.TransformRule
}

// ApplyTransforms runs the configured transformation rules against the
// incoming secrets map and returns the transformed result.
// If cfg is nil or has no rules, the original map is returned unchanged.
func ApplyTransforms(secrets map[string]string, cfg *TransformConfig) (map[string]string, error) {
	if cfg == nil || len(cfg.Rules) == 0 {
		return secrets, nil
	}

	transformed, err := envfile.Transform(secrets, cfg.Rules)
	if err != nil {
		return nil, fmt.Errorf("applyTransforms: %w", err)
	}
	return transformed, nil
}
