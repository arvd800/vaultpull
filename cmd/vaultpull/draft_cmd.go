package main

import (
	"fmt"
	"os"

	"github.com/yourusername/vaultpull/internal/envfile"
)

// runDraft handles the `vaultpull draft` sub-commands: save, show, discard.
func runDraft(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("draft: subcommand required (save|show|discard)")
	}

	draftPath := os.Getenv("VAULTPULL_DRAFT_PATH")
	if draftPath == "" {
		draftPath = ".vaultpull.draft.json"
	}

	switch args[0] {
	case "save":
		return runDraftSave(draftPath, args[1:])
	case "show":
		return runDraftShow(draftPath)
	case "discard":
		return runDraftDiscard(draftPath)
	default:
		return fmt.Errorf("draft: unknown subcommand %q", args[0])
	}
}

func runDraftSave(path string, args []string) error {
	message := ""
	if len(args) > 0 {
		message = args[0]
	}

	existing, err := envfile.Read(".env")
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("draft save: read .env: %w", err)
	}

	d := envfile.NewDraft(existing, message)
	if err := envfile.SaveDraft(path, d); err != nil {
		return fmt.Errorf("draft save: %w", err)
	}
	fmt.Printf("Draft saved: %s\n", d.ID)
	return nil
}

func runDraftShow(path string) error {
	d, err := envfile.LoadDraft(path)
	if err != nil {
		return fmt.Errorf("draft show: %w", err)
	}
	if d.ID == "" {
		fmt.Println("No draft found.")
		return nil
	}
	fmt.Printf("Draft ID : %s\n", d.ID)
	fmt.Printf("Created  : %s\n", d.CreatedAt.Format("2006-01-02 15:04:05 UTC"))
	if d.Message != "" {
		fmt.Printf("Message  : %s\n", d.Message)
	}
	fmt.Printf("Keys     : %d\n", len(d.Secrets))
	for _, line := range envfile.Format(d.Secrets) {
		fmt.Println(" ", line)
	}
	return nil
}

func runDraftDiscard(path string) error {
	if err := envfile.DiscardDraft(path); err != nil {
		return fmt.Errorf("draft discard: %w", err)
	}
	fmt.Println("Draft discarded.")
	return nil
}
