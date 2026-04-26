package envfile

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// DiffSummary holds aggregated statistics about a diff between two secret maps.
type DiffSummary struct {
	Added   []string
	Removed []string
	Changed []string
	Unchanged []string
}

// Summarize produces a DiffSummary from two secret maps.
// It does not expose secret values — only key-level changes are tracked.
func Summarize(existing, incoming map[string]string) DiffSummary {
	summary := DiffSummary{}

	allKeys := make(map[string]struct{})
	for k := range existing {
		allKeys[k] = struct{}{}
	}
	for k := range incoming {
		allKeys[k] = struct{}{}
	}

	for k := range allKeys {
		oldVal, inOld := existing[k]
		newVal, inNew := incoming[k]

		switch {
		case inNew && !inOld:
			summary.Added = append(summary.Added, k)
		case inOld && !inNew:
			summary.Removed = append(summary.Removed, k)
		case oldVal != newVal:
			summary.Changed = append(summary.Changed, k)
		default:
			summary.Unchanged = append(summary.Unchanged, k)
		}
	}

	sort.Strings(summary.Added)
	sort.Strings(summary.Removed)
	sort.Strings(summary.Changed)
	sort.Strings(summary.Unchanged)

	return summary
}

// HasChanges returns true if any keys were added, removed, or changed.
func (s DiffSummary) HasChanges() bool {
	return len(s.Added)+len(s.Removed)+len(s.Changed) > 0
}

// TotalKeys returns the total number of unique keys across both maps.
func (s DiffSummary) TotalKeys() int {
	return len(s.Added) + len(s.Removed) + len(s.Changed) + len(s.Unchanged)
}

// Format writes a human-readable summary to w.
// Each section is only printed if it contains entries.
func (s DiffSummary) Format(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}

	if !s.HasChanges() {
		fmt.Fprintln(w, "No changes detected.")
		return
	}

	if len(s.Added) > 0 {
		fmt.Fprintf(w, "Added (%d):\n", len(s.Added))
		for _, k := range s.Added {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}

	if len(s.Removed) > 0 {
		fmt.Fprintf(w, "Removed (%d):\n", len(s.Removed))
		for _, k := range s.Removed {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(s.Changed) > 0 {
		fmt.Fprintf(w, "Changed (%d):\n", len(s.Changed))
		for _, k := range s.Changed {
			fmt.Fprintf(w, "  ~ %s\n", k)
		}
	}

	fmt.Fprintf(w, "\nSummary: %d added, %d removed, %d changed, %d unchanged (total: %d)\n",
		len(s.Added), len(s.Removed), len(s.Changed), len(s.Unchanged), s.TotalKeys())
}

// FormatString returns the formatted summary as a string.
func (s DiffSummary) FormatString() string {
	var sb strings.Builder
	s.Format(&sb)
	return sb.String()
}
