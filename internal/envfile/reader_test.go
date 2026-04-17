package envfile

import (
	"os"
	"testing"
)

func TestRead_BasicParsing(t *testing.T) {
	content := `# comment
DB_HOST=localhost
DB_PORT=5432
API_KEY="secret123"
`
	f, err := os.CreateTemp("", "*.env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString(content)
	f.Close()

	got, err := Read(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"API_KEY": "secret123",
	}
	for k, v := range expected {
		if got[k] != v {
			t.Errorf("key %s: want %q got %q", k, v, got[k])
		}
	}
}

func TestRead_NonExistentFile(t *testing.T) {
	got, err := Read("/tmp/does_not_exist_vaultpull.env")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestRead_EmptyLines(t *testing.T) {
	content := "\n\n# only comments\n\n"
	f, err := os.CreateTemp("", "*.env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString(content)
	f.Close()

	got, err := Read(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}
