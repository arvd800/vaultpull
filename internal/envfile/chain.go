package envfile

import "fmt"

// ChainEntry represents a single environment in a chain (e.g. base -> staging -> prod).
type ChainEntry struct {
	Name    string            `json:"name"`
	Secrets map[string]string `json:"secrets"`
}

// Chain holds an ordered list of environments whose secrets are merged in order.
type Chain struct {
	Entries []ChainEntry
}

// NewChain creates an empty Chain.
func NewChain() *Chain {
	return &Chain{}
}

// Add appends an environment layer to the chain.
func (c *Chain) Add(name string, secrets map[string]string) {
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	c.Entries = append(c.Entries, ChainEntry{Name: name, Secrets: copy})
}

// Resolve merges all layers in order; later layers override earlier ones.
// Returns an error if any layer contains invalid keys.
func (c *Chain) Resolve() (map[string]string, error) {
	result := make(map[string]string)
	for _, entry := range c.Entries {
		if err := Validate(entry.Secrets); err != nil {
			return nil, fmt.Errorf("chain layer %q: %w", entry.Name, err)
		}
		for k, v := range entry.Secrets {
			result[k] = v
		}
	}
	return result, nil
}

// LayerNames returns the names of all layers in order.
func (c *Chain) LayerNames() []string {
	names := make([]string, len(c.Entries))
	for i, e := range c.Entries {
		names[i] = e.Name
	}
	return names
}
