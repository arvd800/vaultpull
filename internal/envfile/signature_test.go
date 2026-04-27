package envfile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSignMap_Deterministic(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	sig1, err := SignMap(secrets, "secret")
	if err != nil {
		t.Fatal(err)
	}
	sig2, err := SignMap(secrets, "secret")
	if err != nil {
		t.Fatal(err)
	}
	if sig1 != sig2 {
		t.Errorf("expected deterministic signatures, got %q and %q", sig1, sig2)
	}
}

func TestSignMap_DifferentPassphrase(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	sig1, _ := SignMap(secrets, "pass1")
	sig2, _ := SignMap(secrets, "pass2")
	if sig1 == sig2 {
		t.Error("expected different signatures for different passphrases")
	}
}

func TestSignMap_EmptyPassphrase_Error(t *testing.T) {
	_, err := SignMap(map[string]string{"K": "v"}, "")
	if err == nil {
		t.Error("expected error for empty passphrase")
	}
}

func TestVerifySignature_Valid(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	sig, _ := SignMap(secrets, "mypass")
	record := SignatureRecord{Signature: sig, SignedAt: time.Now(), KeyCount: 2}
	if err := VerifySignature(secrets, "mypass", record); err != nil {
		t.Errorf("expected valid signature, got: %v", err)
	}
}

func TestVerifySignature_Tampered(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	sig, _ := SignMap(secrets, "mypass")
	record := SignatureRecord{Signature: sig}
	tampered := map[string]string{"A": "2"}
	if err := VerifySignature(tampered, "mypass", record); err == nil {
		t.Error("expected error for tampered secrets")
	}
}

func TestSaveAndLoadSignature_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sig.json")
	record := SignatureRecord{
		Signature: "abc123",
		SignedAt:  time.Now().UTC().Truncate(time.Second),
		KeyCount:  5,
	}
	if err := SaveSignature(path, record); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadSignature(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Signature != record.Signature {
		t.Errorf("signature mismatch: %q vs %q", loaded.Signature, record.Signature)
	}
	if loaded.KeyCount != record.KeyCount {
		t.Errorf("key count mismatch: %d vs %d", loaded.KeyCount, record.KeyCount)
	}
}

func TestLoadSignature_NonExistent(t *testing.T) {
	record, err := LoadSignature("/nonexistent/path/sig.json")
	if err != nil {
		t.Errorf("expected nil error for missing file, got: %v", err)
	}
	if record.Signature != "" {
		t.Error("expected empty record")
	}
}

func TestSaveSignature_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sig.json")
	_ = SaveSignature(path, SignatureRecord{Signature: "x"})
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}

func TestSaveSignature_EmptyPath_NoOp(t *testing.T) {
	if err := SaveSignature("", SignatureRecord{}); err != nil {
		t.Errorf("expected no error for empty path, got: %v", err)
	}
}

func TestMarshalSorted_OrderIndependent(t *testing.T) {
	a := marshalSorted(map[string]string{"Z": "1", "A": "2"})
	b := marshalSorted(map[string]string{"A": "2", "Z": "1"})
	if a != b {
		t.Errorf("expected order-independent marshalling")
	}
	_ = json.Marshal // suppress unused import
}
