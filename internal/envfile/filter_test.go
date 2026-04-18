package envfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter_IncludePrefix(t *testing.T) {
	input := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db",
	}
	got := Filter(input, FilterOptions{IncludePrefix: "APP_"})
	assert.Equal(t, map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	}, got)
}

func TestFilter_ExcludeKeys(t *testing.T) {
	input := map[string]string{
		"SECRET": "s3cr3t",
		"TOKEN":  "tok",
		"HOST":   "localhost",
	}
	got := Filter(input, FilterOptions{ExcludeKeys: []string{"SECRET", "TOKEN"}})
	assert.Equal(t, map[string]string{"HOST": "localhost"}, got)
}

func TestFilter_StripPrefix(t *testing.T) {
	input := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db",
	}
	got := Filter(input, FilterOptions{IncludePrefix: "APP_", StripPrefix: true})
	assert.Equal(t, map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
	}, got)
}

func TestFilter_StripPrefix_ExactMatch_Skipped(t *testing.T) {
	input := map[string]string{
		"APP_": "empty",
		"APP_X": "x",
	}
	got := Filter(input, FilterOptions{IncludePrefix: "APP_", StripPrefix: true})
	assert.Equal(t, map[string]string{"X": "x"}, got)
}

func TestFilter_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"A": "1", "B": "2"}
	orig := map[string]string{"A": "1", "B": "2"}
	Filter(input, FilterOptions{ExcludeKeys: []string{"A"}})
	assert.Equal(t, orig, input)
}

func TestFilter_EmptyOptions(t *testing.T) {
	input := map[string]string{"A": "1", "B": "2"}
	got := Filter(input, FilterOptions{})
	assert.Equal(t, input, got)
}
