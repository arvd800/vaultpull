package envfile

import "fmt"

// MergeStrategy defines how conflicts are resolved when merging secrets.
type MergeStrategy string

const (
	// StrategyVaultWins overwrites existing local values with vault values.
	StrategyVaultWins MergeStrategy = "vault-wins"
	// StrategyLocalWins keeps existing local values, only adds new keys.
	StrategyLocalWins MergeStrategy = "local-wins"
	// StrategyPromptOnConflict is a marker; callers must handle interactively.
	StrategyPromptOnConflict MergeStrategy = "prompt"
)

// MergeOptions configures how MergeWithStrategy behaves.
type MergeOptions struct {
	Strategy MergeStrategy
	// ConflictKeys receives the list of keys that had conflicts (any strategy).
	ConflictKeys *[]string
}

// MergeWithStrategy merges incoming secrets into existing using the given strategy.
// existing and incoming are not mutated; a new map is returned.
func MergeWithStrategy(existing, incoming map[string]string, opts MergeOptions) (map[string]string, error) {
	if opts.Strategy == "" {
		opts.Strategy = StrategyVaultWins
	}

	switch opts.Strategy {
	case StrategyVaultWins, StrategyLocalWins, StrategyPromptOnConflict:
		// valid
	default:
		return nil, fmt.Errorf("unknown merge strategy: %q", opts.Strategy)
	}

	result := make(map[string]string, len(existing))
	for k, v := range existing {
		result[k] = v
	}

	var conflicts []string

	for k, incomingVal := range incoming {
		existingVal, exists := result[k]
		if exists && existingVal != incomingVal {
			conflicts = append(conflicts, k)
		}

		switch opts.Strategy {
		case StrategyVaultWins:
			result[k] = incomingVal
		case StrategyLocalWins:
			if !exists {
				result[k] = incomingVal
			}
		case StrategyPromptOnConflict:
			// caller handles conflicts; default to vault-wins for non-conflicts
			if !exists {
				result[k] = incomingVal
			}
		}
	}

	if opts.ConflictKeys != nil {
		*opts.ConflictKeys = conflicts
	}

	return result, nil
}
