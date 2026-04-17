package envfile_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/envfile"
)

func TestWrite_BasicSecrets(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env")

	secrets := map[string]string{
		"API_KEY": "abc123",
		"DB_URL":  "postgres://localhost/db",
	}

	if err := envfile.Write(out, secrets); err != nil {
		t.Fatalf("Write: %v", err)
	}

	b, _ := os.ReadFile(out)
	content := string(b)

	if !strings.Contains(content, "API_KEY=abc123") {
		t.Errorf("missing API_KEY line, got:\n%s", content)
	}
	if !strings.Contains(content, "DB_URL=postgres://localhost/db") {
		t.Errorf("missing DB_URL line, got:\n%s", content)
	}
}

func TestWrite_QuotesValueWithSpaces(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env")

	if err := envfile.Write(out, map[string]string{"MSG": "hello world"}); err != nil {
		t.Fatalf("Write: %v", err)
	}

	b, _ := os.ReadFile(out)
	if !strings.Contains(string(b), `MSG="hello world"`) {
		t.Errorf("expected quoted value, got:\n%s", string(b))
	}
}

func TestWrite_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env")
	envfile.Write(out, map[string]string{"X": "1"})

	info, err := os.Stat(out)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 perms, got %v", info.Mode().Perm())
	}
}
