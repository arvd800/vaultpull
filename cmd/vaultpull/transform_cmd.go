package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/your-org/vaultpull/internal/envfile"
)

// runTransform reads a .env file, applies named built-in transforms, and prints the result.
// Usage: vaultpull transform --file=.env --rule=APP_:upper --rule=DB_:trim
func runTransform(args []string) error {
	filePath := ".env"
	format := "dotenv"
	var ruleArgs []string

	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "--file="):
			filePath = strings.TrimPrefix(arg, "--file=")
		case strings.HasPrefix(arg, "--format="):
			format = strings.TrimPrefix(arg, "--format=")
		case strings.HasPrefix(arg, "--rule="):
			ruleArgs = append(ruleArgs, strings.TrimPrefix(arg, "--rule="))
		}
	}

	secrets, err := envfile.Read(filePath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", filePath, err)
	}

	var rules []envfile.TransformRule
	for _, r := range ruleArgs {
		parts := strings.SplitN(r, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid rule %q: expected prefix:transform", r)
		}
		prefix, name := parts[0], parts[1]
		var fn envfile.TransformFunc
		switch name {
		case "upper":
			fn = envfile.UpperCase
		case "lower":
			fn = envfile.LowerCase
		case "trim":
			fn = envfile.TrimSpace
		default:
			return fmt.Errorf("unknown transform %q (supported: upper, lower, trim)", name)
		}
		rules = append(rules, envfile.TransformRule{KeyPrefix: prefix, Transform: fn})
	}

	result, err := envfile.Transform(secrets, rules)
	if err != nil {
		return fmt.Errorf("transforming secrets: %w", err)
	}

	switch format {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(result)
	default:
		return envfile.Export(result, "dotenv", os.Stdout)
	}
}
