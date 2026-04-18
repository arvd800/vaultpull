package envfile

import (
	"testing"
)

func TestEncryptMap_RoundTrip(t *testing.T) {
	secrets := map[string]string{
		"DB_PASS": "hunter2",
		"API_KEY": "abc123",
	}
	passphrase := "test-pass"

	encrypted, err := EncryptMap(secrets, passphrase)
	if err != nil {
		t.Fatalf("EncryptMap failed: %v", err)
	}
	for k, v := range secrets {
		if encrypted[k] == v {
			t.Errorf("key %q: encrypted value should differ from plaintext", k)
		}
	}

	decrypted, err := DecryptMap(encrypted, passphrase)
	if err != nil {
		t.Fatalf("DecryptMap failed: %v", err)
	}
	for k, want := range secrets {
		if got := decrypted[k]; got != want {
			t.Errorf("key %q: expected %q, got %q", k, want, got)
		}
	}
}

func TestEncryptMap_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	_, _ = EncryptMap(secrets, "pass")
	if secrets["KEY"] != "val" {
		t.Fatal("EncryptMap mutated input map")
	}
}

func TestDecryptMap_WrongPassphrase(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	encrypted, _ := EncryptMap(secrets, "correct")
	_, err := DecryptMap(encrypted, "wrong")
	if err == nil {
		t.Fatal("expected error when decrypting with wrong passphrase")
	}
}

func TestEncryptMap_EmptyMap(t *testing.T) {
	out, err := EncryptMap(map[string]string{}, "pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Fatal("expected empty output map")
	}
}
