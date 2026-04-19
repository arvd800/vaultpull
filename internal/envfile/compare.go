package envfile

import "sort"

// CompareResult holds the result of comparing two env maps across environments.
type CompareResult struct {
	OnlyInA  []string
	OnlyInB  []string
	Differ   []string
	Match    []string
}

// Compare compares two secret maps (e.g. staging vs production) and returns
// keys that are unique to each, differ in value, or match.
func Compare(a, b map[string]string) CompareResult {
	seen := make(map[string]struct{})
	result := CompareResult{}

	for k, va := range a {
		seen[k] = struct{}{}
		vb, ok := b[k]
		if !ok {
			result.OnlyInA = append(result.OnlyInA, k)
		} else if va != vb {
			result.Differ = append(result.Differ, k)
		} else {
			result.Match = append(result.Match, k)
		}
	}

	for k := range b {
		if _, ok := seen[k]; !ok {
			result.OnlyInB = append(result.OnlyInB, k)
		}
	}

	sort.Strings(result.OnlyInA)
	sort.Strings(result.OnlyInB)
	sort.Strings(result.Differ)
	sort.Strings(result.Match)

	return result
}
