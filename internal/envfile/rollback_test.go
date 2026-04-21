package envfile_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/envfile"
)

func writeTempHistory(t *testing.T, path string, entries []envfile.HistoryEntry) {
	t.Helper()
	for _, e := range entries {
		if err := envfile.AppendHistory(path, e); err != nil {
			t.Fatalf("writeTempHistory: %v", err)
		}
	}
}

func TestListRollbackPoints_ReturnsAll(t *testing.T) {
	dir := t.TempDir()
	hPath := filepath.Join(dir, "history.json")
	entries := []envfile.HistoryEntry{
		{Timestamp: time.Now().Format(time.RFC3339), Secrets: map[string]string{"A": "1"}},
		{Timestamp: time.Now().Format(time.RFC3339), Secrets: map[string]string{"A": "2", "B": "3"}},
	}
	writeTempHistory(t, hPath, entries)
	points, err := envfile.ListRollbackPoints(hPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(points) != 2 {
		t.Fatalf("expected 2 points, got %d", len(points))
	}
	if points[1].Index != 1 {
		t.Errorf("expected index 1, got %d", points[1].Index)
	}
}

func TestRollback_RestoresSecrets(t *testing.T) {
	dir := t.TempDir()
	hPath := filepath.Join(dir, "history.json")
	target := filepath.Join(dir, ".env")

	original := map[string]string{"KEY": "original"}
	if err := envfile.Write(target, original); err != nil {
		t.Fatalf("write: %v", err)
	}

	entry := envfile.HistoryEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Secrets:   map[string]string{"KEY": "restored", "NEW": "value"},
	}
	writeTempHistory(t, hPath, []envfile.HistoryEntry{entry})

	if err := envfile.Rollback(hPath, target, 0); err != nil {
		t.Fatalf("rollback: %v", err)
	}
	result, err := envfile.Read(target)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if result["KEY"] != "restored" {
		t.Errorf("expected KEY=restored, got %s", result["KEY"])
	}
}

func TestRollback_IndexOutOfRange(t *testing.T) {
	dir := t.TempDir()
	hPath := filepath.Join(dir, "history.json")
	target := filepath.Join(dir, ".env")
	_ = os.WriteFile(target, []byte("KEY=val\n"), 0600)

	entry := envfile.HistoryEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Secrets:   map[string]string{"KEY": "v"},
	}
	writeTempHistory(t, hPath, []envfile.HistoryEntry{entry})

	if err := envfile.Rollback(hPath, target, 5); err == nil {
		t.Error("expected error for out-of-range index")
	}
}

func TestFormatRollbackList_Empty(t *testing.T) {
	out := envfile.FormatRollbackList(nil)
	if out != "no rollback points available" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatRollbackList_ShowsKeys(t *testing.T) {
	points := []envfile.RollbackEntry{
		{Index: 0, Timestamp: "2024-01-01T00:00:00Z", Secrets: map[string]string{"FOO": "bar"}},
	}
	out := envfile.FormatRollbackList(points)
	if out == "" {
		t.Error("expected non-empty output")
	}
}
