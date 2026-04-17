package vault

import (
	"fmt"
	"net/http"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client.
type Client struct {
	logical *vaultapi.Logical
}

// NewClient creates an authenticated Vault client from addr and token.
func NewClient(addr, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = addr
	cfg.HttpClient = &http.Client{Timeout: 10 * time.Second}

	c, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("vault: create client: %w", err)
	}
	c.SetToken(token)

	return &Client{logical: c.Logical()}, nil
}

// ReadSecret reads a KV v2 secret at the given path and returns key/value pairs.
func (c *Client) ReadSecret(path string) (map[string]string, error) {
	secret, err := c.logical.Read(path)
	if err != nil {
		return nil, fmt.Errorf("vault: read %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("vault: no secret found at %q", path)
	}

	data, ok := secret.Data["data"]
	if !ok {
		// KV v1 fallback
		data = secret.Data
	}

	raw, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("vault: unexpected data format at %q", path)
	}

	result := make(map[string]string, len(raw))
	for k, v := range raw {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result, nil
}
