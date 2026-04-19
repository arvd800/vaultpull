package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/envfile"
)

// runChain resolves a layered chain of env files and writes the merged result.
// Usage: vaultpull chain --layers base.env,staging.env --out merged.env
func runChain(cmd *cobra.Command, args []string) error {
	layersFlag, _ := cmd.Flags().GetString("layers")
	outPath, _ := cmd.Flags().GetString("out")

	if layersFlag == "" {
		return fmt.Errorf("--layers is required (comma-separated list of .env files)")
	}

	paths := strings.Split(layersFlag, ",")
	chain := envfile.NewChain()

	for _, p := range paths {
		p = strings.TrimSpace(p)
		secrets, err := envfile.Read(p)
		if err != nil {
			return fmt.Errorf("reading layer %q: %w", p, err)
		}
		chain.Add(p, secrets)
	}

	resolved, err := chain.Resolve()
	if err != nil {
		return fmt.Errorf("resolving chain: %w", err)
	}

	if outPath == "" {
		for _, line := range envfile.Format(resolved) {
			fmt.Println(line)
		}
		return nil
	}

	if err := envfile.Write(outPath, resolved); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	fmt.Fprintf(os.Stderr, "merged %d layers -> %s (%d keys)\n", len(paths), outPath, len(resolved))
	return nil
}

func chainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain",
		Short: "Merge layered .env files in order (later layers override earlier)",
		RunE:  runChain,
	}
	cmd.Flags().String("layers", "", "comma-separated list of .env files to merge in order")
	cmd.Flags().String("out", "", "output .env file (default: stdout)")
	return cmd
}
