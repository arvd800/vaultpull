package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// ExportFormat defines the output format for exported secrets.
type ExportFormat string

const (
	FormatDotenv ExportFormat = "dotenv"
	FormatJSON   ExportFormat = "json"
	FormatExport ExportFormat = "export"
)

// Export serializes a secrets map into the given format as a string.
func Export(secrets map[string]string, format ExportFormat) (string, error) {
	switch format {
	case FormatDotenv:
		return exportDotenv(secrets), nil
	case FormatJSON:
		return exportJSON(secrets)
	case FormatExport:
		return exportShell(secrets), nil
	default:
		return "", fmt.Errorf("unsupported export format: %q", format)
	}
}

// ExportToFile writes the exported content to a file.
func ExportToFile(secrets map[string]string, format ExportFormat, path string) error {
	data, err := Export(secrets, format)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(data), 0600)
}

func exportDotenv(secrets map[string]string) string {
	keys := sortedKeys(secrets)
	var sb strings.Builder
	for _, k := range keys {
		v := secrets[k]
		if need(v) {
			fmt.Fprintf(&sb, "%s=%q\n", k, v)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)	}
	return sb.String()
}

func exportJSON(secrets map[string]string) (string, error) {
	b, err := json.MarshalIndent(secrets, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}

func exportShell(secrets map[string]string) string {
	keys := sortedKeys(secrets)
	var sb strings.Builder
	for _, k := range keys {
		v := secrets[k]
		fmt.Fprintf(&sb, "export %s=%q\n", k, v)
	}
	return sb.String()
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
