package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/yourusername/vaultpull/internal/envfile"
)

func runHistory(args []string) error {
	fs := flag.NewFlagSet("history", flag.ContinueOnError)
	historyPath := fs.String("history-file", ".vaultpull.history.json", "Path to history file")
	limit := fs.Int("n", 10, "Number of recent entries to show")
	if err := fs.Parse(args); err != nil {
		return err
	}

	log, err := envfile.LoadHistory(*historyPath)
	if err != nil {
		return fmt.Errorf("load history: %w", err)
	}

	entries := log.Entries
	if len(entries) == 0 {
		fmt.Println("No history found.")
		return nil
	}

	if *limit > 0 && len(entries) > *limit {
		entries = entries[len(entries)-*limit:]
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tSOURCE\tADDED\tREMOVED\tCHANGED")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%d\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Source,
			len(e.Added),
			len(e.Removed),
			len(e.Changed),
		)
	}
	return w.Flush()
}
