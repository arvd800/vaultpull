package envfile

import "fmt"

// PatchOp represents a single patch operation.
type PatchOp struct {
	Op    string // "set", "delete", "rename"
	Key   string
	Value string // used by "set"
	To    string // used by "rename"
}

// PatchResult holds the outcome of a single patch operation.
type PatchResult struct {
	Op      string
	Key     string
	Applied bool
	Note    string
}

// Patch applies a list of patch operations to a secrets map and returns
// the modified copy along with a result log. It never mutates the input.
func Patch(secrets map[string]string, ops []PatchOp) (map[string]string, []PatchResult, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	results := make([]PatchResult, 0, len(ops))

	for _, op := range ops {
		switch op.Op {
		case "set":
			_, existed := out[op.Key]
			out[op.Key] = op.Value
			note := "created"
			if existed {
				note = "updated"
			}
			results = append(results, PatchResult{Op: op.Op, Key: op.Key, Applied: true, Note: note})

		case "delete":
			_, existed := out[op.Key]
			if existed {
				delete(out, op.Key)
				results = append(results, PatchResult{Op: op.Op, Key: op.Key, Applied: true, Note: "deleted"})
			} else {
				results = append(results, PatchResult{Op: op.Op, Key: op.Key, Applied: false, Note: "key not found"})
			}

		case "rename":
			if op.To == "" {
				return nil, nil, fmt.Errorf("rename op for key %q missing 'to' field", op.Key)
			}
			val, existed := out[op.Key]
			if !existed {
				results = append(results, PatchResult{Op: op.Op, Key: op.Key, Applied: false, Note: "key not found"})
				continue
			}
			delete(out, op.Key)
			out[op.To] = val
			results = append(results, PatchResult{Op: op.Op, Key: op.Key, Applied: true, Note: fmt.Sprintf("renamed to %s", op.To)})

		default:
			return nil, nil, fmt.Errorf("unknown patch op %q", op.Op)
		}
	}

	return out, results, nil
}
