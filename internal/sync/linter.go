package sync

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultpull/internal/envfile"
)

// LintConfig controls lint behaviour during sync.
type LintConfig struct {
	// FailOnWarnings causes RunLint to return an error if any findings exist.
	FailOnWarnings bool
	// Output is the writer for lint output; defaults to os.Stderr.
	Output io.Writer
}

// RunLint lints the provided secrets map and writes findings to the configured
// output. If FailOnWarnings is set and findings exist, an error is returned.
func RunLint(secrets map[string]string, cfg LintConfig) error {
	out := cfg.Output
	if out == nil {
		out = os.Stderr
	}

	results := envfile.Lint(secrets)
	if len(results) == 0 {
		return nil
	}

	fmt.Fprintln(out, "Lint warnings:")
	for _, r := range results {
		fmt.Fprintf(out, "  %s\n", r.String())
	}

	if cfg.FailOnWarnings {
		return fmt.Errorf("lint: %d issue(s) found", len(results))
	}
	return nil
}
