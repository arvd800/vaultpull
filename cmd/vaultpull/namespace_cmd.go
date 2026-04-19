package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultpull/internal/envfile"
)

func runNamespace(args []string) error {
	fs := flag.NewFlagSet("namespace", flag.ContinueOnError)
	nsPath := fs.String("ns-file", ".vaultpull.namespaces.json", "path to namespace definitions file")
	add := fs.String("add", "", "add a namespace: name:prefix")
	remove := fs.String("remove", "", "remove a namespace by name")
	list := fs.Bool("list", false, "list all namespaces")

	if err := fs.Parse(args); err != nil {
		return err
	}

	nsMap, err := envfile.LoadNamespaces(*nsPath)
	if err != nil {
		return fmt.Errorf("load namespaces: %w", err)
	}

	if *list {
		if len(nsMap) == 0 {
			fmt.Println("no namespaces defined")
			return nil
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(nsMap)
	}

	if *add != "" {
		var name, prefix string
		if _, err := fmt.Sscanf(*add, "%s", &name); err != nil {
			return fmt.Errorf("invalid --add format, use name:prefix")
		}
		for i, c := range *add {
			if c == ':' {
				name = (*add)[:i]
				prefix = (*add)[i+1:]
				break
			}
		}
		if name == "" {
			return fmt.Errorf("namespace name required")
		}
		nsMap[name] = envfile.Namespace{Name: name, Prefix: prefix}
		fmt.Printf("added namespace %q with prefix %q\n", name, prefix)
	}

	if *remove != "" {
		delete(nsMap, *remove)
		fmt.Printf("removed namespace %q\n", *remove)
	}

	return envfile.SaveNamespaces(*nsPath, nsMap)
}
