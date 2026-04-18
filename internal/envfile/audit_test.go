package envfile

import (
	"strings"
	"testing"
)

func TestBuildAuditLog_AllActions(t *testing.T) {
	d := DiffResult{
		Added:   []string{"NEW_KEY"},
		Removed: []string{"OLD_KEY"},
		Changed: []string{"MOD_KEY"},
	}
	log := BuildAuditLog(d, "secret/app")
	if len(log) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(log))
	}

	actions := map[string]bool{}
	for _, e := range log {
		actions[e.Action] = true
		if e.Source != "secret/app" {
			t.Errorf("expected source 'secret/app', got %q", e.Source)
		}
		if e.Timestamp.IsZero() {
			t.Error("timestamp should not be zero")
		}
	}
	for _, a := range []string{"added", "removed", "changed"} {
		if !actions[a] {
			t.Errorf("missing action %q in log", a)
		}
	}
}

func TestBuildAuditLog_NoChanges(t *testing.T) {
	d := DiffResult{}
	log := BuildAuditLog(d, "secret/app")
	if len(log) != 0 {
		t.Errorf("expected empty log, got %d entries", len(log))
	}
	if log.Format() != "(no changes)" {
		t.Errorf("unexpected format output: %q", log.Format())
	}
}

func TestAuditLog_Format_ContainsKeyAndAction(t *testing.T) {
	d := DiffResult{
		Added:   []string{"DB_URL"},
		Changed: []string{"API_KEY"},
	}
	log := BuildAuditLog(d, "vault/secret")
	formatted := log.Format()

	if !strings.Contains(formatted, "DB_URL") {
		t.Error("expected DB_URL in formatted output")
	}
	if !strings.Contains(formatted, "added") {
		t.Error("expected 'added' in formatted output")
	}
	if !strings.Contains(formatted, "API_KEY") {
		t.Error("expected API_KEY in formatted output")
	}
	if !strings.Contains(formatted, "changed") {
		t.Error("expected 'changed' in formatted output")
	}
}

func TestAuditEntry_String(t *testing.T) {
	d := DiffResult{Added: []string{"MY_KEY"}}
	log := BuildAuditLog(d, "kv/prod")
	s := log[0].String()
	if !strings.Contains(s, "MY_KEY") {
		t.Errorf("String() missing key: %q", s)
	}
	if !strings.Contains(s, "added") {
		t.Errorf("String() missing action: %q", s)
	}
	if !strings.Contains(s, "kv/prod") {
		t.Errorf("String() missing source: %q", s)
	}
}
