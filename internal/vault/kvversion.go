package vault

import (
	"context"
	"fmt"
	"strings"
)

// KVVersion represents the KV secrets engine version.
type KVVersion int

const (
	KVVersionUnknown KVVersion = 0
	KVVersion1       KVVersion = 1
	KVVersion2       KVVersion = 2
)

// DetectKVVersion attempts to detect the KV engine version for the given path.
// It queries the sys/mounts endpoint and inspects the mount options.
func (c *Client) DetectKVVersion(ctx context.Context, secretPath string) (KVVersion, error) {
	mountPath := extractMountPath(secretPath)

	sys := c.vault.Sys()
	mounts, err := sys.ListMounts()
	if err != nil {
		return KVVersionUnknown, fmt.Errorf("list mounts: %w", err)
	}

	for mount, info := range mounts {
		normalized := strings.Trim(mount, "/")
		if normalized == mountPath {
			if info.Options != nil {
				if v, ok := info.Options["version"]; ok {
					switch v {
					case "1":
						return KVVersion1, nil
					case "2":
						return KVVersion2, nil
					}
				}
			}
			// Default to v1 if no version option present
			return KVVersion1, nil
		}
	}

	return KVVersionUnknown, fmt.Errorf("mount not found for path: %s", secretPath)
}

// extractMountPath returns the first path segment (mount) from a secret path.
func extractMountPath(secretPath string) string {
	parts := strings.SplitN(strings.TrimPrefix(secretPath, "/"), "/", 2)
	return parts[0]
}
