package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// Format takes a map of key-value pairs and returns a sorted, formatted
// .env file content string.
func Format(secrets map[string]string) string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v := secrets[k]
		if needsQuoting(v) {
			sb.WriteString(fmt.Sprintf("%s=%q\n", k, v))
		} else {
			sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
		}
	}
	return sb.String()
}

// ParseLine parses a single KEY=VALUE line from a .env file.
// Returns the key, value, and whether the line was valid.
func ParseLine(line string) (key, value string, ok bool) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false
	}
	idx := strings.IndexByte(line, '=')
	if idx < 1 {
		return "", "", false
	}
	key = strings.TrimSpace(line[:idx])
	value = strings.TrimSpace(line[idx+1:])
	// Strip surrounding quotes if present
	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
		value = value[1 : len(value)-1]
	}
	return key, value, true
}
