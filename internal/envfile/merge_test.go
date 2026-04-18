package envfile

import "testing"

func TestMerge_NewKeysAdded(t *testing.T) {
	existing := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}
	incoming := map[string]string{
		"FOO":     "newbar",
		"NEWKEY":  "newval",
	}

	result := Merge(existing, incoming)

	if result["FOO"] != "newbar" {
		t.Errorf("expected FOO=newbar, got %s", result["FOO"])
	}
	if result["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %s", result["BAZ"])
	}
	if result["NEWKEY"] != "newval" {
		t.Errorf("expected NEWKEY=newval, got %s", result["NEWKEY"])
	}
}

func TestMerge_EmptyExisting(t *testing.T) {
	existing := map[string]string{}
	incoming := map[string]string{"A": "1", "B": "2"}

	result := Merge(existing, incoming)

	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestMerge_EmptyIncoming(t *testing.T) {
	existing := map[string]string{"A": "1"}
	incoming := map[string]string{}

	result := Merge(existing, incoming)

	if len(result) != 1 {
		t.Errorf("expected 1 key, got %d", len(result))
	}
	if result["A"] != "1" {
		t.Errorf("expected A=1, got %s", result["A"])
	}
}

func TestMerge_DoesNotMutateInputs(t *testing.T) {
	existing := map[string]string{"X": "original"}
	incoming := map[string]string{"X": "updated"}

	Merge(existing, incoming)

	if existing["X"] != "original" {
		t.Errorf("Merge mutated existing map")
	}
}

func TestMerge_BothEmpty(t *testing.T) {
	existing := map[string]string{}
	incoming := map[string]string{}

	result := Merge(existing, incoming)

	if len(result) != 0 {
		t.Errorf("expected empty result, got %d keys", len(result))
	}
}
