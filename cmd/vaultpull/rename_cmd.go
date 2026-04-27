package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultpull/internal/envfile"
	syncer "github.com/yourusername/vaultpull/internal/sync"
)

// runRename handles the "rename" sub-command.
//
//	Usage:
//	  vaultpull rename add   --rules <path> --from OLD --to NEW
//	  vaultpull rename apply --rules <path> --env <envfile>
func runRename(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("rename: expected sub-command: add | apply")
	}
	switch args[0] {
	case "add":
		return renameAdd(args[1:])
	case "apply":
		return renameApply(args[1:])
	default:
		return fmt.Errorf("rename: unknown sub-command %q", args[0])
	}
}

func renameAdd(args []string) error {
	fs := flag.NewFlagSet("rename add", flag.ContinueOnError)
	rulesPath := fs.String("rules", ".vaultpull.renames.json", "path to rename rules file")
	from := fs.String("from", "", "original key name")
	to := fs.String("to", "", "new key name")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *from == "" || *to == "" {
		return fmt.Errorf("rename add: --from and --to are required")
	}
	if err := syncer.AddRenameRule(*rulesPath, *from, *to); err != nil {
		return err
	}
	fmt.Printf("rename rule added: %s -> %s\n", *from, *to)
	return nil
}

func renameApply(args []string) error {
	fs := flag.NewFlagSet("rename apply", flag.ContinueOnError)
	rulesPath := fs.String("rules", ".vaultpull.renames.json", "path to rename rules file")
	envPath := fs.String("env", ".env", "path to .env file")
	if err := fs.Parse(args); err != nil {
		return err
	}

	secrets, err := envfile.Read(*envPath)
	if err != nil {
		return fmt.Errorf("rename apply: read env: %w", err)
	}

	renamed, err := syncer.ApplyRenames(secrets, *rulesPath, os.Stdout)
	if err != nil {
		return err
	}

	if err := envfile.Write(*envPath, renamed); err != nil {
		return fmt.Errorf("rename apply: write env: %w", err)
	}
	fmt.Println("rename apply: done")
	return nil
}
