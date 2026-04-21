package envfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/envfile"
)

func TestSaveAndLoadDeprecations_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "deprecations.json")

	dm := envfile.DeprecationMap{}
	dm = envfile.MarkDeprecated(dm, "OLD_API_KEY", "use new auth system", "NEW_API_KEY")
	dm = envfile.MarkDeprecated(dm, "LEGACY_DB_URL", "migrated to secrets manager", "")

	if err := envfile.SaveDeprecations(path, dm); err != nil {
		t.Fatalf("SaveDeprecations: %v", err)
	}

	loaded, err := envfile.LoadDeprecations(path)
	if err != nil {
		t.Fatalf("LoadDeprecations: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(loaded))
	}
	if loaded["OLD_API_KEY"].ReplacedBy != "NEW_API_KEY" {
		t.Errorf("expected ReplacedBy=NEW_API_KEY, got %q", loaded["OLD_API_KEY"].ReplacedBy)
	}
}

func TestLoadDeprecations_NonExistent(t *testing.T) {
	dm, err := envfile.LoadDeprecations("/nonexistent/path.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(dm) != 0 {
		t.Errorf("expected empty map, got %d entries", len(dm))
	}
}

func TestSaveDeprecations_EmptyPath(t *testing.T) {
	if err := envfile.SaveDeprecations("", envfile.DeprecationMap{}); err != nil {
		t.Errorf("expected no error for empty path, got %v", err)
	}
}

func TestSaveDeprecations_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "deprecations.json")
	dm := envfile.MarkDeprecated(envfile.DeprecationMap{}, "KEY", "old", "")
	if err := envfile.SaveDeprecations(path, dm); err != nil {
		t.Fatalf("SaveDeprecations: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected perm 0600, got %v", info.Mode().Perm())
	}
}

func TestCheckDeprecations_WarnsOnMatch(t *testing.T) {
	dm := envfile.MarkDeprecated(envfile.DeprecationMap{}, "OLD_KEY", "use NEW_KEY", "NEW_KEY")
	secrets := map[string]string{
		"OLD_KEY": "value",
		"SAFE_KEY": "ok",
	}
	warnings := envfile.CheckDeprecations(secrets, dm)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if warnings[0] == "" {
		t.Error("expected non-empty warning message")
	}
}

func TestCheckDeprecations_NoMatches(t *testing.T) {
	dm := envfile.MarkDeprecated(envfile.DeprecationMap{}, "OLD_KEY", "old", "")
	secrets := map[string]string{"SAFE_KEY": "value"}
	warnings := envfile.CheckDeprecations(secrets, dm)
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %d", len(warnings))
	}
}

func TestMarkDeprecated_DoesNotMutateInput(t *testing.T) {
	orig := envfile.DeprecationMap{}
	result := envfile.MarkDeprecated(orig, "KEY", "reason", "")
	if len(orig) != 0 {
		t.Error("MarkDeprecated mutated the input map")
	}
	if len(result) != 1 {
		t.Errorf("expected 1 entry in result, got %d", len(result))
	}
}
