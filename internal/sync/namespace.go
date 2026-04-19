package sync

import (
	"fmt"
	"log"

	"github.com/yourusername/vaultpull/internal/envfile"
)

// NamespacedSync applies a namespace filter to incoming secrets before merging.
// It returns the filtered secrets ready for writing.
func NamespacedSync(
	incoming map[string]string,
	existing map[string]string,
	ns envfile.Namespace,
	nsPath string,
) (map[string]string, error) {
	if ns.Name == "" {
		return incoming, nil
	}

	filtered := envfile.ApplyNamespace(incoming, ns)
	if len(filtered) == 0 {
		log.Printf("[namespace] no keys matched prefix %q for namespace %q", ns.Prefix, ns.Name)
	}

	if nsPath != "" {
		nsMap, err := envfile.LoadNamespaces(nsPath)
		if err != nil {
			return nil, fmt.Errorf("load namespaces: %w", err)
		}
		nsMap[ns.Name] = ns
		if err := envfile.SaveNamespaces(nsPath, nsMap); err != nil {
			return nil, fmt.Errorf("save namespaces: %w", err)
		}
	}

	return envfile.Merge(existing, filtered), nil
}
