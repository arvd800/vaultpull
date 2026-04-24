package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestComputeChecksum_Deterministic(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	d1 := ComputeChecksum(secrets)
	d2 := ComputeChecksum(secrets)
	if d1 != d2 {
		t.Fatalf("expected deterministic digest, got %q vs %q", d1, d2)
	}
}

func TestComputeChecksum_OrderIndependent(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAR": "2", "FOO": "1"}
	if ComputeChecksum(a) != ComputeChecksum(b) {
		t.Fatal("checksum should be order-independent")
	}
}

func TestComputeChecksum_DifferentValues(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "baz"}
	if ComputeChecksum(a) == ComputeChecksum(b) {
		t.Fatal("different values should produce different digests")
	}
}

func TestSaveAndLoadChecksum_RoundTrip(t *testing.T) {
	secrets := map[string]string{"KEY": "value", "OTHER": "data"}
	path := filepath.Join(t.TempDir(), "checksum.json")

	if err := SaveChecksum(path, secrets); err != nil {
		t.Fatalf("SaveChecksum: %v", err)
	}

	rec, err := LoadChecksum(path)
	if err != nil {
		t.Fatalf("LoadChecksum: %v", err)
	}
	if rec.Digest != ComputeChecksum(secrets) {
		t.Errorf("digest mismatch: got %q", rec.Digest)
	}
	if rec.KeyCount != 2 {
		t.Errorf("expected key_count 2, got %d", rec.KeyCount)
	}
	if rec.CreatedAt.IsZero() {
		t.Error("expected non-zero created_at")
	}
}

func TestLoadChecksum_NonExistent(t *testing.T) {
	rec, err := LoadChecksum(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Digest != "" {
		t.Errorf("expected empty digest for missing file, got %q", rec.Digest)
	}
}

func TestSaveChecksum_EmptyPath(t *testing.T) {
	if err := SaveChecksum("", map[string]string{"A": "1"}); err != nil {
		t.Fatalf("expected no error for empty path, got %v", err)
	}
}

func TestVerifyChecksum_Match(t *testing.T) {
	secrets := map[string]string{"X": "y"}
	path := filepath.Join(t.TempDir(), "cs.json")
	_ = SaveChecksum(path, secrets)

	ok, err := VerifyChecksum(path, secrets)
	if err != nil {
		t.Fatalf("VerifyChecksum: %v", err)
	}
	if !ok {
		t.Error("expected checksum to match")
	}
}

func TestVerifyChecksum_Mismatch(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cs.json")
	_ = SaveChecksum(path, map[string]string{"X": "original"})

	ok, err := VerifyChecksum(path, map[string]string{"X": "changed"})
	if err != nil {
		t.Fatalf("VerifyChecksum: %v", err)
	}
	if ok {
		t.Error("expected checksum mismatch")
	}
}

func TestVerifyChecksum_NoFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "none.json")
	ok, err := VerifyChecksum(path, map[string]string{"A": "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false when no checksum file exists")
	}
}

func TestSaveChecksum_FilePermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cs.json")
	_ = SaveChecksum(path, map[string]string{"K": "v"})

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected mode 0600, got %v", info.Mode().Perm())
	}
}
