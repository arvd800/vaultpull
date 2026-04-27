package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/your-org/vaultpull/internal/envfile"
)

// runOverride handles the `vaultpull override` subcommand.
// Usage:
//   vaultpull override add   <set-name> <key> <value> [--condition missing|always] [--file path]
//   vaultpull override apply <set-name> <env-file> [--overrides path]
func runOverride(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: vaultpull override <add|apply> [args]")
	}
	switch args[0] {
	case "add":
		return runOverrideAdd(args[1:])
	case "apply":
		return runOverrideApply(args[1:])
	default:
		return fmt.Errorf("unknown override subcommand: %q", args[0])
	}
}

func runOverrideAdd(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: vaultpull override add <set-name> <key> <value> [--condition missing|always] [--file path]")
	}
	setName, key, value := args[0], args[1], args[2]
	condition := "always"
	filePath := ".vaultpull-overrides.json"
	for i := 3; i < len(args)-1; i++ {
		switch args[i] {
		case "--condition":
			condition = args[i+1]
		case "--file":
			filePath = args[i+1]
		}
	}
	sets, err := envfile.LoadOverrides(filePath)
	if err != nil {
		return err
	}
	ov := envfile.Override{Key: key, Value: value, Condition: condition}
	for i, s := range sets {
		if s.Name == setName {
			sets[i].Overrides = append(sets[i].Overrides, ov)
			return envfile.SaveOverrides(filePath, sets)
		}
	}
	sets = append(sets, envfile.OverrideSet{Name: setName, Overrides: []envfile.Override{ov}})
	if err := envfile.SaveOverrides(filePath, sets); err != nil {
		return err
	}
	fmt.Printf("override: added key %q to set %q\n", key, setName)
	return nil
}

func runOverrideApply(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: vaultpull override apply <set-name> <env-file> [--overrides path]")
	}
	setName, envFile := args[0], args[1]
	overridesPath := ".vaultpull-overrides.json"
	for i := 2; i < len(args)-1; i++ {
		if args[i] == "--overrides" {
			overridesPath = args[i+1]
		}
	}
	secrets, err := envfile.Read(envFile)
	if err != nil {
		return fmt.Errorf("override apply: read env: %w", err)
	}
	sets, err := envfile.LoadOverrides(overridesPath)
	if err != nil {
		return err
	}
	result, err := envfile.ApplyOverrides(secrets, sets, setName)
	if err != nil {
		return err
	}
	out, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintln(os.Stdout, string(out))
	return nil
}
