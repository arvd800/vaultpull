package envfile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var sampleSecrets = map[string]string{
	"DB_HOST": "localhost",
	"DB_PASS": "s3cr3t word",
	"API_KEY": "abc123",
}

func TestExport_Dotenv(t *testing.T) {
	out, err := Export(sampleSecrets, FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST line, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_PASS=") {
		t.Errorf("expected DB_PASS line, got:\n%s", out)
	}
}

func TestExport_JSON(t *testing.T) {
	out, err := Export(sampleSecrets, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]string
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if parsed["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %s", parsed["API_KEY"])
	}
}

func TestExport_Shell(t *testing.T) {
	out, err := Export(sampleSecrets, FormatExport)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export API_KEY=") {
		t.Errorf("expected export prefix, got:\n%s", out)
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	_, err := Export(sampleSecrets, ExportFormat("xml"))
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestExportToFile_WritesCorrectly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.env")

	if err := ExportToFile(sampleSecrets, FormatDotenv, path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if !strings.Contains(string(data), "API_KEY=abc123") {
		t.Errorf("expected API_KEY in file, got:\n%s", string(data))
	}
}

func TestExportToFile_Permissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secrets.env")

	if err := ExportToFile(sampleSecrets, FormatDotenv, path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
