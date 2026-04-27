package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/yourusername/vaultpull/internal/envfile"
)

// runSignature handles the "signature" subcommand with actions: sign, verify, show.
func runSignature(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: vaultpull signature <sign|verify|show> [flags]")
	}

	action := args[0]
	envPath := envFlagOrDefault(args[1:], "--env", ".env")
	sigPath := envFlagOrDefault(args[1:], "--sig", ".env.sig")
	passphrase := os.Getenv("VAULTPULL_SIG_PASSPHRASE")

	switch action {
	case "sign":
		return runSignatureSign(envPath, sigPath, passphrase)
	case "verify":
		return runSignatureVerify(envPath, sigPath, passphrase)
	case "show":
		return runSignatureShow(sigPath)
	default:
		return fmt.Errorf("unknown signature action %q", action)
	}
}

func runSignatureSign(envPath, sigPath, passphrase string) error {
	if passphrase == "" {
		return fmt.Errorf("VAULTPULL_SIG_PASSPHRASE must be set")
	}
	secrets, err := envfile.Read(envPath)
	if err != nil {
		return fmt.Errorf("read env: %w", err)
	}
	sig, err := envfile.SignMap(secrets, passphrase)
	if err != nil {
		return err
	}
	record := envfile.SignatureRecord{
		Signature: sig,
		SignedAt:  time.Now().UTC(),
		KeyCount:  len(secrets),
	}
	if err := envfile.SaveSignature(sigPath, record); err != nil {
		return fmt.Errorf("save signature: %w", err)
	}
	fmt.Printf("signed %d keys → %s\n", len(secrets), sigPath)
	return nil
}

func runSignatureVerify(envPath, sigPath, passphrase string) error {
	if passphrase == "" {
		return fmt.Errorf("VAULTPULL_SIG_PASSPHRASE must be set")
	}
	secrets, err := envfile.Read(envPath)
	if err != nil {
		return fmt.Errorf("read env: %w", err)
	}
	record, err := envfile.LoadSignature(sigPath)
	if err != nil {
		return fmt.Errorf("load signature: %w", err)
	}
	if record.Signature == "" {
		return fmt.Errorf("no signature found at %s", sigPath)
	}
	if err := envfile.VerifySignature(secrets, passphrase, record); err != nil {
		return err
	}
	fmt.Println("signature valid ✓")
	return nil
}

func runSignatureShow(sigPath string) error {
	record, err := envfile.LoadSignature(sigPath)
	if err != nil {
		return fmt.Errorf("load signature: %w", err)
	}
	if record.Signature == "" {
		fmt.Println("no signature file found")
		return nil
	}
	out, _ := json.MarshalIndent(record, "", "  ")
	fmt.Println(string(out))
	return nil
}

func envFlagOrDefault(args []string, flag, def string) string {
	for i, a := range args {
		if a == flag && i+1 < len(args) {
			return args[i+1]
		}
	}
	return def
}
