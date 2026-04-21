package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/your-org/vaultpull/internal/sync"
)

func runRollback(args []string) error {
	fs := flag.NewFlagSet("rollback", flag.ContinueOnError)

	target := fs.String("env", ".env", "path to the .env file to restore")
	history := fs.String("history", "", "path to history file (default: <env-dir>/.vaultpull/history.json)")
	index := fs.Int("index", -1, "history index to roll back to (required unless --list)")
	list := fs.Bool("list", false, "list available rollback points and exit")
	dryRun := fs.Bool("dry-run", false, "print what would be restored without writing")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *list || *dryRun {
		return sync.RunRollback(sync.RollbackOptions{
			HistoryPath: *history,
			TargetPath:  *target,
			Index:       *index,
			DryRun:      true,
		})
	}

	if *index < 0 {
		fmt.Fprintln(os.Stderr, "error: --index is required (use --list to see available points)")
		fs.Usage()
		return fmt.Errorf("missing --index")
	}

	return sync.RunRollback(sync.RollbackOptions{
		HistoryPath: *history,
		TargetPath:  *target,
		Index:       *index,
		DryRun:      false,
	})
}
