package vault

import (
	"context"
	"fmt"
)

// ReadSecretAuto reads a secret from Vault, automatically detecting the KV version.
// It falls back to KVv1 if version detection fails.
func (c *Client) ReadSecretAuto(ctx context.Context, secretPath string) (map[string]string, error) {
	version, err := c.DetectKVVersion(ctx, secretPath)
	if err != nil {
		// Fall back to KVv1 on detection failure
		version = KVVersion1
	}

	switch version {
	case KVVersion2:
		return c.ReadSecretKVv2(ctx, secretPath)
	case KVVersion1:
		return c.ReadSecret(ctx, secretPath)
	default:
		return nil, fmt.Errorf("unsupported KV version: %d", version)
	}
}
