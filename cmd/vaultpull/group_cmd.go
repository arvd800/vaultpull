package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/your-org/vaultpull/internal/envfile"
	"github.com/your-org/vaultpull/internal/sync"
)

// runGroup handles the `vaultpull group` subcommand.
// Usage:
//
//	vaultpull group --groups groups.json --list
//	vaultpull group --groups groups.json --add <name> --keys KEY1,KEY2
func runGroup(args []string) error {
	fs := flag.NewFlagSet("group", flag.ContinueOnError)
	groupsPath := fs.String("groups", ".vaultpull-groups.json", "path to groups file")
	list := fs.Bool("list", false, "list all defined groups")
	add := fs.String("add", "", "name of group to add or update")
	keys := fs.String("keys", "", "comma-separated list of keys for --add")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *list {
		return groupList(*groupsPath)
	}

	if *add != "" {
		if *keys == "" {
			return fmt.Errorf("group: --keys is required when using --add")
		}
		keySlice := strings.Split(*keys, ",")
		for i, k := range keySlice {
			keySlice[i] = strings.TrimSpace(k)
		}
		if err := sync.RegisterGroup(*groupsPath, *add, keySlice); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "group %q saved to %s\n", *add, *groupsPath)
		return nil
	}

	fs.Usage()
	return fmt.Errorf("group: no action specified; use --list or --add")
}

func groupList(path string) error {
	groups, err := envfile.LoadGroups(path)
	if err != nil {
		return fmt.Errorf("group: %w", err)
	}
	if len(groups) == 0 {
		fmt.Println("no groups defined")
		return nil
	}
	names := make([]string, 0, len(groups))
	for n := range groups {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		fmt.Printf("%s: %s\n", n, strings.Join(groups[n], ", "))
	}
	return nil
}
