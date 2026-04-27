package sync

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultpull/internal/envfile"
)

// ApplyRenames loads rename rules from rulesPath and applies them to secrets.
// If rulesPath is empty the secrets are returned unchanged.
// Renamed keys are reported to w (os.Stdout if nil).
func ApplyRenames(secrets map[string]string, rulesPath string, w io.Writer) (map[string]string, error) {
	if rulesPath == "" {
		return secrets, nil
	}
	if w == nil {
		w = os.Stdout
	}

	rm, err := envfile.LoadRenames(rulesPath)
	if err != nil {
		return nil, fmt.Errorf("renamer: load rules: %w", err)
	}
	if len(rm.Rules) == 0 {
		return secrets, nil
	}

	out := envfile.ApplyRenames(secrets, rm)

	for _, rule := range rm.Rules {
		if _, hadFrom := secrets[rule.From]; hadFrom {
			fmt.Fprintf(w, "rename: %s -> %s\n", rule.From, rule.To)
		}
	}
	return out, nil
}

// AddRenameRule appends a rename rule to the file at rulesPath.
func AddRenameRule(rulesPath, from, to string) error {
	if rulesPath == "" {
		return fmt.Errorf("renamer: rules path is required")
	}
	if from == "" || to == "" {
		return fmt.Errorf("renamer: from and to must not be empty")
	}
	rm, err := envfile.LoadRenames(rulesPath)
	if err != nil {
		return err
	}
	rm.Rules = append(rm.Rules, envfile.RenameRule{From: from, To: to})
	return envfile.SaveRenames(rulesPath, rm)
}
