package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultpull/internal/envfile"
)

// runMergeStrategy merges a JSON secrets file into an existing .env file using
// the specified strategy and writes the result back out.
func runMergeStrategy(cmd *cobra.Command, args []string) error {
	envPath, _ := cmd.Flags().GetString("env")
	inputPath, _ := cmd.Flags().GetString("input")
	strategy, _ := cmd.Flags().GetString("strategy")
	showConflicts, _ := cmd.Flags().GetBool("show-conflicts")

	if envPath == "" || inputPath == "" {
		return fmt.Errorf("--env and --input are required")
	}

	existing, err := envfile.Read(envPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading env file: %w", err)
	}
	if existing == nil {
		existing = map[string]string{}
	}

	raw, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}
	var incoming map[string]string
	if err := json.Unmarshal(raw, &incoming); err != nil {
		return fmt.Errorf("parsing input JSON: %w", err)
	}

	var conflicts []string
	opts := envfile.MergeOptions{
		Strategy:     envfile.MergeStrategy(strategy),
		ConflictKeys: &conflicts,
	}

	merged, err := envfile.MergeWithStrategy(existing, incoming, opts)
	if err != nil {
		return err
	}

	if showConflicts && len(conflicts) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "conflicts (%d): %v\n", len(conflicts), conflicts)
	}

	if err := envfile.Write(envPath, merged); err != nil {
		return fmt.Errorf("writing env file: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "merged %d keys into %s (strategy=%s)\n", len(merged), envPath, strategy)
	return nil
}

func mergeStrategyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge-strategy",
		Short: "Merge secrets into a .env file using a conflict resolution strategy",
		RunE:  runMergeStrategy,
	}
	cmd.Flags().String("env", "", "path to the target .env file")
	cmd.Flags().String("input", "", "path to JSON file containing incoming secrets")
	cmd.Flags().String("strategy", "vault-wins", "merge strategy: vault-wins | local-wins | prompt")
	cmd.Flags().Bool("show-conflicts", false, "print conflicting keys to stdout")
	return cmd
}
