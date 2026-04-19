package envfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChain_ResolveMergesLayers(t *testing.T) {
	c := NewChain()
	c.Add("base", map[string]string{"APP_ENV": "base", "DB_HOST": "localhost"})
	c.Add("staging", map[string]string{"APP_ENV": "staging"})

	out, err := c.Resolve()
	require.NoError(t, err)
	assert.Equal(t, "staging", out["APP_ENV"])
	assert.Equal(t, "localhost", out["DB_HOST"])
}

func TestChain_ResolveEmptyChain(t *testing.T) {
	c := NewChain()
	out, err := c.Resolve()
	require.NoError(t, err)
	assert.Empty(t, out)
}

func TestChain_ResolveInvalidKey(t *testing.T) {
	c := NewChain()
	c.Add("base", map[string]string{"VALID_KEY": "ok"})
	c.Add("bad", map[string]string{"1INVALID": "nope"})

	_, err := c.Resolve()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "bad")
}

func TestChain_DoesNotMutateInput(t *testing.T) {
	original := map[string]string{"KEY": "value"}
	c := NewChain()
	c.Add("layer", original)
	original["KEY"] = "mutated"

	out, err := c.Resolve()
	require.NoError(t, err)
	assert.Equal(t, "value", out["KEY"])
}

func TestChain_LayerNames(t *testing.T) {
	c := NewChain()
	c.Add("base", map[string]string{})
	c.Add("staging", map[string]string{})
	c.Add("prod", map[string]string{})

	assert.Equal(t, []string{"base", "staging", "prod"}, c.LayerNames())
}

func TestChain_LaterLayerWins(t *testing.T) {
	c := NewChain()
	c.Add("a", map[string]string{"X": "1", "Y": "2"})
	c.Add("b", map[string]string{"X": "99"})
	c.Add("c", map[string]string{"X": "final"})

	out, err := c.Resolve()
	require.NoError(t, err)
	assert.Equal(t, "final", out["X"])
	assert.Equal(t, "2", out["Y"])
}
