package envfile

import (
	"fmt"
	"io"
	"sort"
)

// PrintOptions controls output behaviour of PrintEnv.
type PrintOptions struct {
	Redact  bool
	RedactOpts *RedactOptions
}

// PrintEnv writes secrets to w in KEY=VALUE format, optionally redacting
// sensitive values before printing.
func PrintEnv(w io.Writer, secrets map[string]string, opts *PrintOptions) error {
	data := secrets
	if opts != nil && opts.Redact {
		data = Redact(secrets, opts.RedactOpts)
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := data[k]
		if needsQuoting(v) {
			_, err := fmt.Fprintf(w, "%s=%q\n", k, v)
			if err != nil {
				return err
			}
		} else {
			_, err := fmt.Fprintf(w, "%s=%s\n", k, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
