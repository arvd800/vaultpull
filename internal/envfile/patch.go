package envfile

import (
	"fmt"
	"strings"
)

// PatchOp represents a single patch operation to apply to a secrets map.
type PatchOp struct {
	// Op is the operation type: "set", "delete", or "rename".
	Op string
	// Key is the target key.
	Key string
	// Value is used for "set" operations.
	Value string
	// NewKey is used for "rename" operations.
	NewKey string
}

// PatchResult summarises the outcome of a Patch call.
type PatchResult struct {
	Applied []string
	Skipped []string
	Errors  []string
}

// Patch applies a slice of PatchOps to src and returns a new map with the
// changes applied. src is never mutated. Unknown op types are recorded in
// PatchResult.Errors but do not abort processing.
func Patch(src map[string]string, ops []PatchOp) (map[string]string, PatchResult) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}

	var result PatchResult

	for _, op := range ops {
		switch strings.ToLower(op.Op) {
		case "set":
			if op.Key == "" {
				result.Errors = append(result.Errors, "set op missing key")
				continue
			}
			out[op.Key] = op.Value
			result.Applied = append(result.Applied, fmt.Sprintf("set %s", op.Key))

		case "delete":
			if op.Key == "" {
				result.Errors = append(result.Errors, "delete op missing key")
				continue
			}
			if _, ok := out[op.Key]; !ok {
				result.Skipped = append(result.Skipped, fmt.Sprintf("delete %s (not found)", op.Key))
				continue
			}
			delete(out, op.Key)
			result.Applied = append(result.Applied, fmt.Sprintf("delete %s", op.Key))

		case "rename":
			if op.Key == "" || op.NewKey == "" {
				result.Errors = append(result.Errors, "rename op requires both key and new_key")
				continue
			}
			val, ok := out[op.Key]
			if !ok {
				result.Skipped = append(result.Skipped, fmt.Sprintf("rename %s (not found)", op.Key))
				continue
			}
			out[op.NewKey] = val
			delete(out, op.Key)
			result.Applied = append(result.Applied, fmt.Sprintf("rename %s -> %s", op.Key, op.NewKey))

		default:
			result.Errors = append(result.Errors, fmt.Sprintf("unknown op %q for key %q", op.Op, op.Key))
		}
	}

	return out, result
}
