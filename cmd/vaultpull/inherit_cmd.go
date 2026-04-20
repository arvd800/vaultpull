package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/your-org/vaultpull/internal/envfile"
)

// runInherit implements the `vaultpull inherit` sub-command.
// Usage:
//
//	vaultpull inherit --inherit=.vaultpull.inherit.json --parent=.env.base --child=.env
func runInherit(args []string) error {
	var inheritFile, parentFile, childFile string
	var showOnly bool

	for _, a := range args {
		switch {
		case len(a) > 10 && a[:10] == "--inherit=":
			inheritFile = a[10:]
		case len(a) > 9 && a[:9] == "--parent=":
			parentFile = a[9:]
		case len(a) > 8 && a[:8] == "--child=":
			childFile = a[8:]
		case a == "--show":
			showOnly = true
		}
	}

	if inheritFile == "" {
		return fmt.Errorf("inherit: --inherit path is required")
	}

	m, err := envfile.LoadInherit(inheritFile)
	if err != nil {
		return fmt.Errorf("inherit: load map: %w", err)
	}

	if showOnly {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(m)
	}

	if parentFile == "" || childFile == "" {
		return fmt.Errorf("inherit: --parent and --child are required")
	}

	parent, err := envfile.Read(parentFile)
	if err != nil {
		return fmt.Errorf("inherit: read parent %q: %w", parentFile, err)
	}

	child, err := envfile.Read(childFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("inherit: read child %q: %w", childFile, err)
	}
	if child == nil {
		child = make(map[string]string)
	}

	result := envfile.ApplyInherit(child, parent, m)

	if err := envfile.Write(childFile, result); err != nil {
		return fmt.Errorf("inherit: write child %q: %w", childFile, err)
	}

	fmt.Fprintf(os.Stdout, "inherit: applied %d key(s) from %q to %q\n",
		len(result)-len(child), parentFile, childFile)
	return nil
}
