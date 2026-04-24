package envfile

import (
	"testing"
)

func TestPrune_RemoveKeys(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PASS": "secret",
		"API_KEY": "abc123",
	}
	res, err := Prune(secrets, PruneOptions{RemoveKeys: []string{"DB_PASS", "API_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(res.Removed))
	}
	if _, ok := res.Retained["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to be retained")
	}
	if _, ok := res.Retained["DB_PASS"]; ok {
		t.Error("expected DB_PASS to be removed")
	}
}

func TestPrune_KeepKeys(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PASS": "secret",
		"API_KEY": "abc123",
	}
	res, err := Prune(secrets, PruneOptions{KeepKeys: []string{"DB_HOST"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Retained) != 1 {
		t.Errorf("expected 1 retained, got %d", len(res.Retained))
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(res.Removed))
	}
}

func TestPrune_DryRun_DoesNotModify(t *testing.T) {
	secrets := map[string]string{
		"A": "1",
		"B": "2",
	}
	res, err := Prune(secrets, PruneOptions{RemoveKeys: []string{"A"}, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Retained) != 2 {
		t.Errorf("dry run: expected all keys retained, got %d", len(res.Retained))
	}
	if len(res.Removed) != 1 || res.Removed[0] != "A" {
		t.Errorf("dry run: expected removed=[A], got %v", res.Removed)
	}
}

func TestPrune_MutuallyExclusiveOptions(t *testing.T) {
	_, err := Prune(map[string]string{"X": "1"}, PruneOptions{
		KeepKeys:   []string{"X"},
		RemoveKeys: []string{"X"},
	})
	if err == nil {
		t.Error("expected error when both KeepKeys and RemoveKeys are set")
	}
}

func TestPrune_NilInput(t *testing.T) {
	res, err := Prune(nil, PruneOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Retained) != 0 {
		t.Errorf("expected empty retained for nil input")
	}
}

func TestPrune_NoOptions_ReturnsAll(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res, err := Prune(secrets, PruneOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Retained) != 2 {
		t.Errorf("expected all keys retained, got %d", len(res.Retained))
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected no removed keys, got %v", res.Removed)
	}
}
