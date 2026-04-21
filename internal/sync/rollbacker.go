package sync

import (
	"fmt"
	"path/filepath"

	"github.com/your-org/vaultpull/internal/envfile"
)

// RollbackOptions configures the rollback operation.
type RollbackOptions struct {
	// HistoryPath is the path to the history JSON file.
	HistoryPath string
	// TargetPath is the .env file to restore into.
	TargetPath string
	// Index is the history entry index to roll back to.
	Index int
	// DryRun prints what would be restored without writing.
	DryRun bool
}

// RunRollback performs a rollback of the target .env file to a prior history state.
func RunRollback(opts RollbackOptions) error {
	if opts.HistoryPath == "" {
		opts.HistoryPath = filepath.Join(
			envfile.RollbackDir(opts.TargetPath),
			"history.json",
		)
	}

	points, err := envfile.ListRollbackPoints(opts.HistoryPath)
	if err != nil {
		return fmt.Errorf("rollback: %w", err)
	}

	if opts.DryRun {
		fmt.Println(envfile.FormatRollbackList(points))
		if opts.Index >= 0 && opts.Index < len(points) {
			fmt.Printf("\n[dry-run] would restore index %d (%s)\n",
				opts.Index, points[opts.Index].Timestamp)
		}
		return nil
	}

	if err := envfile.Rollback(opts.HistoryPath, opts.TargetPath, opts.Index); err != nil {
		return fmt.Errorf("rollback: %w", err)
	}
	fmt.Printf("rolled back %s to history index %d\n", opts.TargetPath, opts.Index)
	return nil
}
