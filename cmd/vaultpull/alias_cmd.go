package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yourusername/vaultpull/internal/envfile"
)

// runAlias handles the `vaultpull alias` subcommand.
// Usage:
//
//	vaultpull alias set <alias> <canonical>  --file aliases.json
//	vaultpull alias list                     --file aliases.json
func runAlias(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("alias: expected subcommand: set | list")
	}

	aliasFile := os.Getenv("VAULTPULL_ALIAS_FILE")
	if aliasFile == "" {
		aliasFile = ".vaultpull.aliases.json"
	}

	switch args[0] {
	case "set":
		if len(args) < 3 {
			return fmt.Errorf("alias set: usage: alias set <alias> <canonical>")
		}
		aliasName := args[1]
		canonical := args[2]

		aliases, err := envfile.LoadAliases(aliasFile)
		if err != nil {
			return fmt.Errorf("alias set: load: %w", err)
		}
		aliases[aliasName] = canonical
		if err := envfile.SaveAliases(aliasFile, aliases); err != nil {
			return fmt.Errorf("alias set: save: %w", err)
		}
		fmt.Printf("alias %q -> %q saved to %s\n", aliasName, canonical, aliasFile)

	case "list":
		aliases, err := envfile.LoadAliases(aliasFile)
		if err != nil {
			return fmt.Errorf("alias list: %w", err)
		}
		if len(aliases) == 0 {
			fmt.Println("no aliases defined")
			return nil
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(aliases)

	default:
		return fmt.Errorf("alias: unknown subcommand %q", args[0])
	}
	return nil
}
