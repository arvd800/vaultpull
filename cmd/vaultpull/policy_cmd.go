package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yourusername/vaultpull/internal/envfile"
)

func runPolicy(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: vaultpull policy <check|add|show> [flags]")
	}
	switch args[0] {
	case "check":
		return runPolicyCheck(args[1:])
	case "add":
		return runPolicyAdd(args[1:])
	case "show":
		return runPolicyShow(args[1:])
	default:
		return fmt.Errorf("unknown policy subcommand: %s", args[0])
	}
}

func runPolicyCheck(args []string) error {
	policyPath := ".vaultpull.policy.json"
	envPath := ".env"
	if len(args) >= 1 {
		policyPath = args[0]
	}
	if len(args) >= 2 {
		envPath = args[1]
	}
	secrets, err := envfile.Read(envPath)
	if err != nil {
		return fmt.Errorf("read env: %w", err)
	}
	policy, err := envfile.LoadPolicy(policyPath)
	if err != nil {
		return fmt.Errorf("load policy: %w", err)
	}
	violations := envfile.EnforcePolicy(secrets, policy)
	if len(violations) == 0 {
		fmt.Println("policy check passed: no violations")
		return nil
	}
	fmt.Fprintf(os.Stderr, "%d policy violation(s):\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(os.Stderr, "  - %s\n", v.Error())
	}
	return fmt.Errorf("policy check failed")
}

func runPolicyAdd(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: vaultpull policy add <policy-file> <key> [--required] [--deny] [--pattern=<re>]")
	}
	policyPath := args[0]
	key := args[1]
	rule := envfile.PolicyRule{Key: key}
	for _, a := range args[2:] {
		switch {
		case a == "--required":
			rule.Required = true
		case a == "--deny":
			rule.Deny = true
		case len(a) > 10 && a[:10] == "--pattern=":
			rule.Pattern = a[10:]
		}
	}
	policy, err := envfile.LoadPolicy(policyPath)
	if err != nil {
		return err
	}
	policy.Rules = append(policy.Rules, rule)
	if err := envfile.SavePolicy(policyPath, policy); err != nil {
		return err
	}
	fmt.Printf("added policy rule for key %q\n", key)
	return nil
}

func runPolicyShow(args []string) error {
	policyPath := ".vaultpull.policy.json"
	if len(args) >= 1 {
		policyPath = args[0]
	}
	policy, err := envfile.LoadPolicy(policyPath)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(policy)
}
