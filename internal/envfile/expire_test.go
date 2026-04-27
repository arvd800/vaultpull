package envfile_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/envfile"
)

func TestSaveAndLoadExpiry_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "expiry.json")

	expiries := envfile.ExpiryMap{
		"API_KEY": {Key: "API_KEY", ExpiresAt: time.Now().UTC().Add(time.Hour), Note: "rotate soon"},
	}

	if err := envfile.SaveExpiry(path, expiries); err != nil {
		t.Fatalf("SaveExpiry: %v", err)
	}

	loaded, err := envfile.LoadExpiry(path)
	if err != nil {
		t.Fatalf("LoadExpiry: %v", err)
	}

	rec, ok := loaded["API_KEY"]
	if !ok {
		t.Fatal("expected API_KEY in loaded expiries")
	}
	if rec.Note != "rotate soon" {
		t.Errorf("note mismatch: got %q", rec.Note)
	}
}

func TestLoadExpiry_NonExistent(t *testing.T) {
	loaded, err := envfile.LoadExpiry("/nonexistent/expiry.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(loaded) != 0 {
		t.Errorf("expected empty map, got %v", loaded)
	}
}

func TestSaveExpiry_EmptyPath(t *testing.T) {
	if err := envfile.SaveExpiry("", envfile.ExpiryMap{}); err != nil {
		t.Fatalf("expected no-op, got %v", err)
	}
}

func TestSetExpiry_DoesNotMutateInput(t *testing.T) {
	original := envfile.ExpiryMap{}
	updated := envfile.SetExpiry(original, "DB_PASS", time.Hour, "")
	if len(original) != 0 {
		t.Error("original map was mutated")
	}
	if _, ok := updated["DB_PASS"]; !ok {
		t.Error("DB_PASS not found in updated map")
	}
}

func TestCheckExpiry_Expired(t *testing.T) {
	secrets := map[string]string{"OLD_TOKEN": "abc", "FRESH_KEY": "xyz"}
	expiries := envfile.ExpiryMap{
		"OLD_TOKEN": {Key: "OLD_TOKEN", ExpiresAt: time.Now().UTC().Add(-time.Minute)},
		"FRESH_KEY": {Key: "FRESH_KEY", ExpiresAt: time.Now().UTC().Add(time.Hour)},
	}
	expired := envfile.CheckExpiry(secrets, expiries)
	if len(expired) != 1 || expired[0] != "OLD_TOKEN" {
		t.Errorf("expected [OLD_TOKEN], got %v", expired)
	}
}

func TestCheckExpiry_NoneExpired(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	expiries := envfile.ExpiryMap{
		"KEY": {Key: "KEY", ExpiresAt: time.Now().UTC().Add(time.Hour)},
	}
	expired := envfile.CheckExpiry(secrets, expiries)
	if len(expired) != 0 {
		t.Errorf("expected no expired keys, got %v", expired)
	}
}

func TestSaveExpiry_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "expiry.json")
	if err := envfile.SaveExpiry(path, envfile.ExpiryMap{}); err != nil {
		t.Fatalf("SaveExpiry: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
}
