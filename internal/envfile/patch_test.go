package envfile

import (
	"testing"
)

func baseSecrets() map[string]string {
	return map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"API_KEY": "secret",
	}
}

func TestPatch_SetNewKey(t *testing.T) {
	out, results, err := Patch(baseSecrets(), []PatchOp{{Op: "set", Key: "NEW_KEY", Value: "hello"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", out["NEW_KEY"])
	}
	if results[0].Note != "created" {
		t.Errorf("expected note 'created', got %q", results[0].Note)
	}
}

func TestPatch_SetExistingKey(t *testing.T) {
	out, results, err := Patch(baseSecrets(), []PatchOp{{Op: "set", Key: "DB_HOST", Value: "remotehost"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "remotehost" {
		t.Errorf("expected DB_HOST=remotehost, got %q", out["DB_HOST"])
	}
	if results[0].Note != "updated" {
		t.Errorf("expected note 'updated', got %q", results[0].Note)
	}
}

func TestPatch_DeleteKey(t *testing.T) {
	out, results, err := Patch(baseSecrets(), []PatchOp{{Op: "delete", Key: "API_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["API_KEY"]; ok {
		t.Error("expected API_KEY to be deleted")
	}
	if !results[0].Applied {
		t.Error("expected Applied=true")
	}
}

func TestPatch_DeleteMissingKey(t *testing.T) {
	_, results, err := Patch(baseSecrets(), []PatchOp{{Op: "delete", Key: "MISSING"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Applied {
		t.Error("expected Applied=false for missing key")
	}
}

func TestPatch_RenameKey(t *testing.T) {
	out, results, err := Patch(baseSecrets(), []PatchOp{{Op: "rename", Key: "DB_PORT", To: "DATABASE_PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_PORT"]; ok {
		t.Error("expected DB_PORT to be removed")
	}
	if out["DATABASE_PORT"] != "5432" {
		t.Errorf("expected DATABASE_PORT=5432, got %q", out["DATABASE_PORT"])
	}
	if !results[0].Applied {
		t.Error("expected Applied=true")
	}
}

func TestPatch_RenameKeyMissingTo(t *testing.T) {
	_, _, err := Patch(baseSecrets(), []PatchOp{{Op: "rename", Key: "DB_PORT"}})
	if err == nil {
		t.Error("expected error for rename without 'to'")
	}
}

func TestPatch_UnknownOp(t *testing.T) {
	_, _, err := Patch(baseSecrets(), []PatchOp{{Op: "upsert", Key: "X"}})
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestPatch_DoesNotMutateInput(t *testing.T) {
	src := baseSecrets()
	origLen := len(src)
	Patch(src, []PatchOp{
		{Op: "set", Key: "EXTRA", Value: "val"},
		{Op: "delete", Key: "DB_HOST"},
	})
	if len(src) != origLen {
		t.Errorf("input was mutated: expected %d keys, got %d", origLen, len(src))
	}
	if src["DB_HOST"] != "localhost" {
		t.Error("input DB_HOST was mutated")
	}
}
