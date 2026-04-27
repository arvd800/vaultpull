package envfile

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewDraft_PopulatesFields(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	d := NewDraft(secrets, "initial draft")

	if d.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if d.Message != "initial draft" {
		t.Errorf("expected message 'initial draft', got %q", d.Message)
	}
	if len(d.Secrets) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(d.Secrets))
	}
	if d.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestNewDraft_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"KEY": "value"}
	d := NewDraft(secrets, "")
	d.Secrets["KEY"] = "mutated"

	if secrets["KEY"] != "value" {
		t.Error("NewDraft mutated the input map")
	}
}

func TestSaveAndLoadDraft_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "draft.json")

	orig := NewDraft(map[string]string{"A": "1", "B": "2"}, "test draft")
	orig.CreatedAt = orig.CreatedAt.Truncate(time.Second)

	if err := SaveDraft(path, orig); err != nil {
		t.Fatalf("SaveDraft: %v", err)
	}

	loaded, err := LoadDraft(path)
	if err != nil {
		t.Fatalf("LoadDraft: %v", err)
	}

	if loaded.ID != orig.ID {
		t.Errorf("ID mismatch: got %q, want %q", loaded.ID, orig.ID)
	}
	if loaded.Message != orig.Message {
		t.Errorf("Message mismatch: got %q, want %q", loaded.Message, orig.Message)
	}
	if loaded.Secrets["A"] != "1" || loaded.Secrets["B"] != "2" {
		t.Errorf("Secrets mismatch: %v", loaded.Secrets)
	}
}

func TestLoadDraft_NonExistent(t *testing.T) {
	d, err := LoadDraft("/nonexistent/draft.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if d.ID != "" {
		t.Errorf("expected empty Draft, got ID=%q", d.ID)
	}
}

func TestSaveDraft_EmptyPath_NoOp(t *testing.T) {
	err := SaveDraft("", Draft{})
	if err != nil {
		t.Fatalf("expected no error for empty path, got: %v", err)
	}
}

func TestDiscardDraft_RemovesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "draft.json")

	if err := SaveDraft(path, NewDraft(map[string]string{"X": "y"}, "")); err != nil {
		t.Fatalf("SaveDraft: %v", err)
	}
	if err := DiscardDraft(path); err != nil {
		t.Fatalf("DiscardDraft: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected file to be removed")
	}
}

func TestDiscardDraft_NonExistent_NoError(t *testing.T) {
	if err := DiscardDraft("/nonexistent/draft.json"); err != nil {
		t.Fatalf("expected no error discarding non-existent file, got: %v", err)
	}
}
