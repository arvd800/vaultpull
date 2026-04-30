package envfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMergeWithStrategy_VaultWins(t *testing.T) {
	existing := map[string]string{"A": "old", "B": "keep"}
	incoming := map[string]string{"A": "new", "C": "added"}

	result, err := MergeWithStrategy(existing, incoming, MergeOptions{Strategy: StrategyVaultWins})
	require.NoError(t, err)
	assert.Equal(t, "new", result["A"])
	assert.Equal(t, "keep", result["B"])
	assert.Equal(t, "added", result["C"])
}

func TestMergeWithStrategy_LocalWins(t *testing.T) {
	existing := map[string]string{"A": "local", "B": "keep"}
	incoming := map[string]string{"A": "vault", "C": "added"}

	result, err := MergeWithStrategy(existing, incoming, MergeOptions{Strategy: StrategyLocalWins})
	require.NoError(t, err)
	assert.Equal(t, "local", result["A"], "local value should be preserved")
	assert.Equal(t, "keep", result["B"])
	assert.Equal(t, "added", result["C"], "new key should be added")
}

func TestMergeWithStrategy_ConflictTracking(t *testing.T) {
	existing := map[string]string{"A": "old", "B": "same"}
	incoming := map[string]string{"A": "new", "B": "same", "C": "extra"}

	var conflicts []string
	_, err := MergeWithStrategy(existing, incoming, MergeOptions{
		Strategy:     StrategyVaultWins,
		ConflictKeys: &conflicts,
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"A"}, conflicts)
}

func TestMergeWithStrategy_DefaultStrategy(t *testing.T) {
	existing := map[string]string{"A": "old"}
	incoming := map[string]string{"A": "new"}

	result, err := MergeWithStrategy(existing, incoming, MergeOptions{})
	require.NoError(t, err)
	assert.Equal(t, "new", result["A"], "default should be vault-wins")
}

func TestMergeWithStrategy_UnknownStrategy(t *testing.T) {
	_, err := MergeWithStrategy(nil, nil, MergeOptions{Strategy: "bad"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown merge strategy")
}

func TestMergeWithStrategy_DoesNotMutateInputs(t *testing.T) {
	existing := map[string]string{"A": "orig"}
	incoming := map[string]string{"A": "new"}

	_, err := MergeWithStrategy(existing, incoming, MergeOptions{Strategy: StrategyVaultWins})
	require.NoError(t, err)
	assert.Equal(t, "orig", existing["A"])
}

func TestMergeWithStrategy_PromptOnConflict_AddsNewKeys(t *testing.T) {
	existing := map[string]string{"A": "local"}
	incoming := map[string]string{"A": "vault", "B": "new"}

	var conflicts []string
	result, err := MergeWithStrategy(existing, incoming, MergeOptions{
		Strategy:     StrategyPromptOnConflict,
		ConflictKeys: &conflicts,
	})
	require.NoError(t, err)
	assert.Equal(t, "local", result["A"], "conflict key should remain local")
	assert.Equal(t, "new", result["B"])
	assert.Equal(t, []string{"A"}, conflicts)
}
