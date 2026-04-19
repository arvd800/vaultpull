package envfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPromote_AllKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"C": "3"}
	out, res, err := Promote(src, dst, PromoteOptions{})
	require.NoError(t, err)
	assert.Equal(t, "1", out["A"])
	assert.Equal(t, "2", out["B"])
	assert.Equal(t, "3", out["C"])
	assert.Len(t, res.Promoted, 2)
	assert.Empty(t, res.Skipped)
}

func TestPromote_SelectedKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{}
	out, res, err := Promote(src, dst, PromoteOptions{Keys: []string{"A", "C"}})
	require.NoError(t, err)
	assert.Equal(t, "1", out["A"])
	assert.Equal(t, "3", out["C"])
	_, hasB := out["B"]
	assert.False(t, hasB)
	assert.Len(t, res.Promoted, 2)
}

func TestPromote_SkipExisting(t *testing.T) {
	src := map[string]string{"A": "new", "B": "2"}
	dst := map[string]string{"A": "old"}
	out, res, err := Promote(src, dst, PromoteOptions{SkipExisting: true})
	require.NoError(t, err)
	assert.Equal(t, "old", out["A"], "existing key should not be overwritten")
	assert.Equal(t, "2", out["B"])
	assert.Contains(t, res.Skipped, "A")
	assert.Contains(t, res.Promoted, "B")
}

func TestPromote_NilSrc(t *testing.T) {
	_, _, err := Promote(nil, map[string]string{}, PromoteOptions{})
	assert.Error(t, err)
}

func TestPromote_NilDst(t *testing.T) {
	src := map[string]string{"X": "1"}
	out, res, err := Promote(src, nil, PromoteOptions{})
	require.NoError(t, err)
	assert.Equal(t, "1", out["X"])
	assert.Len(t, res.Promoted, 1)
}

func TestPromote_DoesNotMutateInputs(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{"B": "2"}
	Promote(src, dst, PromoteOptions{})
	assert.Len(t, src, 1)
	assert.Len(t, dst, 1)
}
