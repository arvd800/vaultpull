package envfile

import (
	"testing"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	plaintext := "super-secret-value"
	passphrase := "my-passphrase"

	encrypted, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	if encrypted == plaintext {
		t.Fatal("encrypted value should differ from plaintext")
	}

	decrypted, err := Decrypt(encrypted, passphrase)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}
	if decrypted != plaintext {
		t.Fatalf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestEncrypt_ProducesUniqueOutputs(t *testing.T) {
	plaintext := "value"
	passphrase := "pass"
	a, _ := Encrypt(plaintext, passphrase)
	b, _ := Encrypt(plaintext, passphrase)
	if a == b {
		t.Fatal("two encryptions of the same value should differ (random nonce)")
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	encrypted, _ := Encrypt("secret", "correct")
	_, err := Decrypt(encrypted, "wrong")
	if err == nil {
		t.Fatal("expected error decrypting with wrong passphrase")
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	_, err := Decrypt("!!!not-base64!!!", "pass")
	if err == nil {
		t.Fatal("expected error on invalid base64")
	}
}

func TestDecrypt_TooShort(t *testing.T) {
	import64 := "aGk=" // "hi" in base64 — too short for nonce
	_, err := Decrypt(import64, "pass")
	if err == nil {
		t.Fatal("expected error on too-short ciphertext")
	}
}
