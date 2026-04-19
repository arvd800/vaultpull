package sync_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/envfile"
	"github.com/yourusername/vaultpull/internal/sync"
)

func TestApplyTags_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tags.json")

	err := sync.ApplyTags(sync.TagOptions{
		TagsFilePath: path,
		Source:       "secret/app",
		ExtraTags:    map[string]string{"env": "staging"},
	})
	if err != nil {
		t.Fatalf("ApplyTags: %v", err)
	}

	record, err := envfile.LoadTags(path)
	if err != nil {
		t.Fatalf("LoadTags: %v", err)
	}
	if record.Source != "secret/app" {
		t.Errorf("Source: got %q", record.Source)
	}
	if record.Tags["env"] != "staging" {
		t.Errorf("env tag: got %q", record.Tags["env"])
	}
	if record.FetchedAt.IsZero() {
		t.Error("FetchedAt should not be zero")
	}
}

func TestApplyTags_MergesExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tags.json")

	initial := envfile.TagRecord{
		Source:    "old",
		FetchedAt: time.Now(),
		Tags:      map[string]string{"team": "infra"},
	}
	_ = envfile.SaveTags(path, initial)

	_ = sync.ApplyTags(sync.TagOptions{
		TagsFilePath: path,
		Source:       "secret/new",
		ExtraTags:    map[string]string{"env": "prod"},
	})

	record, _ := envfile.LoadTags(path)
	if record.Tags["team"] != "infra" {
		t.Errorf("expected existing tag preserved, got %q", record.Tags["team"])
	}
	if record.Tags["env"] != "prod" {
		t.Errorf("expected new tag, got %q", record.Tags["env"])
	}
}

func TestApplyTags_NoPath_NoOp(t *testing.T) {
	err := sync.ApplyTags(sync.TagOptions{})
	if err != nil {
		t.Errorf("expected no error for empty path, got %v", err)
	}
}
