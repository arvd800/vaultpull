package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultpull/internal/envfile"
)

func runLock(args []string) error {
	fs := flag.NewFlagSet("lock", flag.ContinueOnError)
	lockFile := fs.String("lock-file", ".vaultpull.locks.json", "path to lock file")
	reason := fs.String("reason", "", "reason for locking the key")
	unlock := fs.Bool("unlock", false, "unlock the key instead of locking")

	if err := fs.Parse(args); err != nil {
		return err
	}
	keys := fs.Args()
	if len(keys) == 0 {
		return fmt.Errorf("usage: vaultpull lock [--unlock] [--reason=...] KEY [KEY...]")
	}

	locks, err := envfile.LoadLocks(*lockFile)
	if err != nil {
		return fmt.Errorf("loading lock file: %w", err)
	}

	for _, key := range keys {
		if *unlock {
			locks = envfile.UnlockKey(locks, key)
			fmt.Fprintf(os.Stdout, "unlocked %q\n", key)
		} else {
			locks = envfile.LockKey(locks, key, *reason)
			fmt.Fprintf(os.Stdout, "locked %q\n", key)
		}
	}

	if err := envfile.SaveLocks(*lockFile, locks); err != nil {
		return fmt.Errorf("saving lock file: %w", err)
	}
	return nil
}
