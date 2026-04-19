package envfile_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envfile"
)

func TestCompare_OnlyInA(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"FOO": "1"}
	r := envfile.Compare(a, b)
	if len(r.OnlyInA) != 1 || r.OnlyInA[0] != "BAR" {
		t.Errorf("expected OnlyInA=[BAR], got %v", r.OnlyInA)
	}
}

func TestCompare_OnlyInB(t *testing.T) {
	a := map[string]string{"FOO": "1"}
	b := map[string]string{"FOO": "1", "BAZ": "3"}
	r := envfile.Compare(a, b)
	if len(r.OnlyInB) != 1 || r.OnlyInB[0] != "BAZ" {
		t.Errorf("expected OnlyInB=[BAZ], got %v", r.OnlyInB)
	}
}

func TestCompare_Differ(t *testing.T) {
	a := map[string]string{"FOO": "old"}
	b := map[string]string{"FOO": "new"}
	r := envfile.Compare(a, b)
	if len(r.Differ) != 1 || r.Differ[0] != "FOO" {
		t.Errorf("expected Differ=[FOO], got %v", r.Differ)
	}
}

func TestCompare_Match(t *testing.T) {
	a := map[string]string{"FOO": "same"}
	b := map[string]string{"FOO": "same"}
	r := envfile.Compare(a, b)
	if len(r.Match) != 1 || r.Match[0] != "FOO" {
		t.Errorf("expected Match=[FOO], got %v", r.Match)
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	r := envfile.Compare(map[string]string{}, map[string]string{})
	if len(r.OnlyInA)+len(r.OnlyInB)+len(r.Differ)+len(r.Match) != 0 {
		t.Error("expected all empty slices for empty maps")
	}
}

func TestCompare_SortedOutput(t *testing.T) {
	a := map[string]string{"Z": "1", "A": "1", "M": "1"}
	b := map[string]string{}
	r := envfile.Compare(a, b)
	for i := 1; i < len(r.OnlyInA); i++ {
		if r.OnlyInA[i] < r.OnlyInA[i-1] {
			t.Error("OnlyInA is not sorted")
		}
	}
}
