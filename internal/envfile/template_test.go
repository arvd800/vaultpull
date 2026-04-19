package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempTemplate(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestRenderTemplate_Basic(t *testing.T) {
	tmplPath := writeTempTemplate(t, "DB_HOST={{ index . \"DB_HOST\" }}\nDB_PORT={{ index . \"DB_PORT\" }}\n")
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}

	out, err := RenderTemplate(tmplPath, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "DB_HOST=localhost\nDB_PORT=5432\n"
	if out != expected {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestRenderTemplate_MissingKey(t *testing.T) {
	tmplPath := writeTempTemplate(t, "VAL={{ index . \"MISSING_KEY\" }}\n")
	secrets := map[string]string{}

	_, err := RenderTemplate(tmplPath, secrets)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestRenderTemplate_NonExistentFile(t *testing.T) {
	_, err := RenderTemplate("/nonexistent/path.tmpl", map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing template file")
	}
}

func TestRenderTemplateToFile_WritesOutput(t *testing.T) {
	tmplPath := writeTempTemplate(t, "SECRET={{ index . \"MY_SECRET\" }}\n")
	destPath := filepath.Join(t.TempDir(), ".env")
	secrets := map[string]string{"MY_SECRET": "hunter2"}

	if err := RenderTemplateToFile(tmplPath, destPath, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "SECRET=hunter2\n" {
		t.Errorf("unexpected file content: %q", string(data))
	}

	info, _ := os.Stat(destPath)
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected perm 0600, got %v", info.Mode().Perm())
	}
}
