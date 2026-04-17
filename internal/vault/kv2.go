package vault

import (
	"fmt"
	"strings"
)

// ReadSecretKVv2 reads a secret from a KV v2 mount, handling the data/metadata
// wrapper that Vault uses for versioned secrets.
func (c *Client) ReadSecretKVv2(path string) (map[string]string, error) {
	// KV v2 secrets are stored under <mount>/data/<path>
	// Normalise: if the caller already included "data/" we skip adding it.
	dataPath := injectDataSegment(path)

	secret, err := c.logical.Read(dataPath)
	if err != nil {
		return nil, fmt.Errorf("kv2 read %q: %w", dataPath, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("kv2 read %q: secret not found", dataPath)
	}

	// KV v2 wraps values under secret.Data["data"]
	raw, ok := secret.Data["data"]
	if !ok {
		return nil, fmt.Errorf("kv2 read %q: missing 'data' key in response", dataPath)
	}

	nested, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("kv2 read %q: unexpected data format", dataPath)
	}

	out := make(map[string]string, len(nested))
	for k, v := range nested {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out, nil
}

// injectDataSegment ensures the path contains the /data/ segment required by
// the KV v2 API (e.g. "secret/myapp" -> "secret/data/myapp").
func injectDataSegment(path string) string {
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return path
	}
	mount, rest := parts[0], parts[1]
	if strings.HasPrefix(rest, "data/") {
		return path
	}
	return mount + "/data/" + rest
}
