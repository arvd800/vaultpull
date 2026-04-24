package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/vaultpull/internal/envfile"
)

// runPreset handles the `vaultpull preset` subcommand family.
// Usage:
//
//	vaultpull preset list
//	vaultpull preset add <name> KEY=VALUE ...
//	vaultpull preset apply <name>
func runPreset(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("preset: expected subcommand: list, add, apply")
	}
	dir, _ := os.Getwd()
	switch args[0] {
	case "list":
		return presetList(dir)
	case "add":
		if len(args) < 3 {
			return fmt.Errorf("preset add: usage: preset add <name> KEY=VALUE ...")
		}
		return presetAdd(dir, args[1], args[2:])
	case "apply":
		if len(args) < 2 {
			return fmt.Errorf("preset apply: usage: preset apply <name>")
		}
		return presetApply(dir, args[1])
	default:
		return fmt.Errorf("preset: unknown subcommand %q", args[0])
	}
}

func presetList(dir string) error {
	presets, err := envfile.LoadPresets(dir)
	if err != nil {
		return err
	}
	if len(presets) == 0 {
		fmt.Println("no presets defined")
		return nil
	}
	for _, p := range presets {
		fmt.Printf("  %s (%d keys)\n", p.Name, len(p.Values))
	}
	return nil
}

func presetAdd(dir, name string, pairs []string) error {
	values := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("preset add: invalid pair %q", pair)
		}
		values[parts[0]] = parts[1]
	}
	presets, err := envfile.LoadPresets(dir)
	if err != nil {
		return err
	}
	presets = append(presets, envfile.Preset{Name: name, Values: values})
	if err := envfile.SavePresets(dir, presets); err != nil {
		return err
	}
	fmt.Printf("preset %q saved with %d key(s)\n", name, len(values))
	return nil
}

func presetApply(dir, name string) error {
	presets, err := envfile.LoadPresets(dir)
	if err != nil {
		return err
	}
	secrets := map[string]string{}
	result, err := envfile.ApplyPreset(name, presets, secrets)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}
