package envfile

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAppendAndLoadHistory_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	entry := HistoryEntry{
		Timestamp: time.Now().UTC(),
		Source:    "vault/secret/app",
		Added:     []string{"DB_HOST"},
		Removed:   []string{},
		Changed:   []string{"API_KEY"},
		Snapshot:  map[string]string{"DB_HOST": "localhost", "API_KEY": "new"},
	}

	if err := AppendHistory(path, entry); err != nil {
		t.Fatalf("AppendHistory: %v", err)
	}

	log, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	if log.Entries[0].Source != "vault/secret/app" {
		t.Errorf("unexpected source: %s", log.Entries[0].Source)
	}
}

func TestAppendHistory_Accumulates(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	for i := 0; i < 3; i++ {
		_ = AppendHistory(path, HistoryEntry{Source: "src", Snapshot: map[string]string{}})
	}
	log, _ := LoadHistory(path)
	if len(log.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(log.Entries))
	}
}

func TestLoadHistory_NonExistent(t *testing.T) {
	log, err := LoadHistory("/tmp/nonexistent-history-xyz.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(log.Entries) != 0 {
		t.Errorf("expected empty log")
	}
}

func TestAppendHistory_EmptyPath(t *testing.T) {
	 err := AppendHistory("", HistoryEntry{}); err != nil {
		t.Errorf("expected no error for empty path, got %v", err)
	}
}

func TestBuildHistoryEntry_PopulatesFields(t *testing.T) {
	d := DiffResult{
		Added:   map[string]string{"NEW": "val"},
		Removed: map[string]string{"OLD": "gone"},
		Changed: map[string]string{},
	}
	current := map[string]string{"NEW": "val", "KEEP": "same"}
	entry := BuildHistoryEntry("vault/test", d, current)

	if entry.Source != "vault/test" {
		t.Errorf("unexpected source")
	}
	if len(entry.Added) != 1 || entry.Added[0] != "NEW" {
		t.Errorf("expected Added=[NEW]")
	}
	if entry.Snapshot["KEEP"] != "same" {
		t.Errorf("snapshot missing KEEP")
	}
	if entry.Timestamp.IsZero() {
		t.Errorf("timestamp not set")
	}
}

func TestSaveHistory_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")
	_ = AppendHistory(path, HistoryEntry{Source: "x", Snapshot: map[string]string{}})
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
