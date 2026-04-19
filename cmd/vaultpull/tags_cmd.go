package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yourusername/vaultpull/internal/envfile"
)

// runShowTags prints the tags file for a given .env output path.
func runShowTags(tagsPath string, asJSON bool) error {
	record, err := envfile.LoadTags(tagsPath)
	if err != nil {
		return fmt.Errorf("load tags: %w", err)
	}

	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(record)
	}

	fmt.Printf("Source:     %s\n", record.Source)
	fmt.Printf("FetchedAt:  %s\n", record.FetchedAt.Format("2006-01-02 15:04:05 UTC"))
	if len(record.Tags) == 0 {
		fmt.Println("Tags:       (none)")
		return nil
	}
	fmt.Println("Tags:")
	for k, v := range record.Tags {
		fmt.Printf("  %s = %s\n", k, v)
	}
	return nil
}
