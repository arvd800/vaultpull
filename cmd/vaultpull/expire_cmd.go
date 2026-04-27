package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/your-org/vaultpull/internal/envfile"
)

// runExpire handles the `vaultpull expire` subcommand.
//
// Usage:
//   vaultpull expire set   <key> <duration> [note]   -- mark a key with an expiry TTL
//   vaultpull expire check <envfile> <expiryfile>    -- report expired keys
func runExpire(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("expire: subcommand required (set|check)")
	}

	switch args[0] {
	case "set":
		return runExpireSet(args[1:])
	case "check":
		return runExpireCheck(args[1:])
	default:
		return fmt.Errorf("expire: unknown subcommand %q", args[0])
	}
}

func runExpireSet(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("expire set: usage: <expiry-file> <key> <duration> [note]")
	}
	path, key, rawDur := args[0], args[1], args[2]
	note := ""
	if len(args) >= 4 {
		note = args[3]
	}

	// Accept seconds as plain integer or Go duration string.
	var dur time.Duration
	if secs, err := strconv.ParseInt(rawDur, 10, 64); err == nil {
		dur = time.Duration(secs) * time.Second
	} else {
		var perr error
		dur, perr = time.ParseDuration(rawDur)
		if perr != nil {
			return fmt.Errorf("expire set: invalid duration %q: %w", rawDur, perr)
		}
	}

	expiries, err := envfile.LoadExpiry(path)
	if err != nil {
		return fmt.Errorf("expire set: load: %w", err)
	}

	expiries = envfile.SetExpiry(expiries, key, dur, note)

	if err := envfile.SaveExpiry(path, expiries); err != nil {
		return fmt.Errorf("expire set: save: %w", err)
	}
	fmt.Fprintf(os.Stdout, "expire: set expiry for %q (%s)\n", key, dur)
	return nil
}

func runExpireCheck(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("expire check: usage: <env-file> <expiry-file>")
	}
	envPath, expiryPath := args[0], args[1]

	secrets, err := envfile.Read(envPath)
	if err != nil {
		return fmt.Errorf("expire check: read env: %w", err)
	}

	expiries, err := envfile.LoadExpiry(expiryPath)
	if err != nil {
		return fmt.Errorf("expire check: load expiry: %w", err)
	}

	expired := envfile.CheckExpiry(secrets, expiries)
	if len(expired) == 0 {
		fmt.Fprintln(os.Stdout, "expire: no expired keys")
		return nil
	}

	fmt.Fprintf(os.Stderr, "expire: %d expired key(s):\n", len(expired))
	for _, k := range expired {
		fmt.Fprintf(os.Stderr, "  - %s (expired %s)\n", k, expiries[k].ExpiresAt.Format(time.RFC3339))
	}
	return fmt.Errorf("expire: %d key(s) have expired", len(expired))
}
