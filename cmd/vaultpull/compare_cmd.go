package main

import (
	"fmt"
	"os"

	"github.com/your-org/vaultpull/internal/envfile"
)

// runCompare loads two .env files and prints a diff summary between them.
func runCompare(fileA, fileB string) error {
	a, err := envfile.Read(fileA)
	if err != nil {
		return fmt.Errorf("reading %s: %w", fileA, err)
	}
	b, err := envfile.Read(fileB)
	if err != nil {
		return fmt.Errorf("reading %s: %w", fileB, err)
	}

	r := envfile.Compare(a, b)

	w := os.Stdout

	if len(r.OnlyInA) > 0 {
		fmt.Fprintf(w, "Only in %s:\n", fileA)
		for _, k := range r.OnlyInA {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(r.OnlyInB) > 0 {
		fmt.Fprintf(w, "Only in %s:\n", fileB)
		for _, k := range r.OnlyInB {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}

	if len(r.Differ) > 0 {
		fmt.Fprintf(w, "Values differ:\n")
		for _, k := range r.Differ {
			fmt.Fprintf(w, "  ~ %s\n", k)
		}
	}

	fmt.Fprintf(w, "Matching keys: %d\n", len(r.Match))
	return nil
}
