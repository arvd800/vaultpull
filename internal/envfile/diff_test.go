package envfile

import (
	"testing"
)

func TestDiff_AddedKeys(t *testing.T) {
	existing := map[string]string{"FOO": "bar"}
	incoming := map[string]string{"FOO": "bar", "NEW": "value"}

	result := Diff(existing, incoming)

	if len(result.Added) != 1 || result.Added["NEW"] != "value" {
		t.Errorf("expected NEW=value in Added, got %v", result.Added)
	}
	if len(result.Removed) != 0 || len(result.Changed) != 0 {
		t.Errorf("unexpected changes: removed=%v changed=%v", result.Removed, result.Changed)
	}
}

func TestDiff_RemovedKeys(t *testing.T) {
	existing := map[string]string{"FOO": "bar", "OLD": "gone"}
	incoming := map[string]string{"FOO": "bar"}

	result := Diff(existing, incoming)

	if len(result.Removed) != 1 || result.Removed["OLD"] != "gone" {
		t.Errorf("expected OLD in Removed, got %v", result.Removed)
	}
}

func TestDiff_ChangedKeys(t *testing.T) {
	existing := map[string]string{"FOO": "old"}
	incoming := map[string]string{"FOO": "new"}

	result := Diff(existing, incoming)

	if len(result.Changed) != 1 || result.Changed["FOO"] != "new" {
		t.Errorf("expected FOO=new in Changed, got %v", result.Changed)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	existing := map[string]string{"FOO": "bar", "BAZ": "qux"}
	incoming := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := Diff(existing, incoming)

	if result.HasChanges() {
		t.Errorf("expected no changes, got %+v", result)
	}
}

func TestDiff_EmptyExisting(t *testing.T) {
	existing := map[string]string{}
	incoming := map[string]string{"A": "1", "B": "2"}

	result := Diff(existing, incoming)

	if len(result.Added) != 2 {
		t.Errorf("expected 2 added keys, got %d", len(result.Added))
	}
	if result.HasChanges() == false {
		t.Error("expected HasChanges to be true")
	}
}

func TestDiff_BothEmpty(t *testing.T) {
	result := Diff(map[string]string{}, map[string]string{})
	if result.HasChanges() {
		t.Error("expected no changes for two empty maps")
	}
}
