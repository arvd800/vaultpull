package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/your-org/vaultpull/internal/envfile"
)

// runObfuscate replaces secret values in a .env file with opaque tokens and
// writes the lookup table to a separate JSON file for later restoration.
func runObfuscate(args []string) error {
	fs := flag.NewFlagSet("obfuscate", flag.ContinueOnError)
	input := fs.String("input", ".env", "path to source .env file")
	output := fs.String("output", ".env.obfuscated", "path to write obfuscated .env file")
	lookupOut := fs.String("lookup", ".env.lookup.json", "path to write token lookup table (keep secret!)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	secrets, err := envfile.Read(*input)
	if err != nil {
		return fmt.Errorf("obfuscate: read %q: %w", *input, err)
	}

	obfuscated, lookup, err := envfile.ObfuscateMap(secrets)
	if err != nil {
		return fmt.Errorf("obfuscate: %w", err)
	}

	if err := envfile.Write(*output, obfuscated); err != nil {
		return fmt.Errorf("obfuscate: write output: %w", err)
	}

	f, err := os.OpenFile(*lookupOut, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("obfuscate: open lookup file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(lookup); err != nil {
		return fmt.Errorf("obfuscate: encode lookup: %w", err)
	}

	fmt.Fprintf(os.Stdout, "obfuscated %d keys → %s (lookup: %s)\n",
		len(obfuscated), *output, *lookupOut)
	return nil
}
