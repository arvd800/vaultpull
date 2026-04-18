package envfile

import (
	"fmt"
	"os"
)

// WriteEncrypted writes secrets to path after encrypting all values with passphrase.
// The resulting file is a valid .env file with base64-encoded encrypted values.
func WriteEncrypted(path string, secrets map[string]string, passphrase string) error {
	encrypted, err := EncryptMap(secrets, passphrase)
	if err != nil {
		return fmt.Errorf("encrypting secrets: %w", err)
	}
	return Write(path, encrypted)
}

// ReadDecrypted reads an encrypted .env file and returns plaintext secrets.
func ReadDecrypted(path string, passphrase string) (map[string]string, error) {
	raw, err := Read(path)
	if err != nil {
		return nil, fmt.Errorf("reading encrypted file: %w", err)
	}
	decrypted, err := DecryptMap(raw, passphrase)
	if err != nil {
		return nil, fmt.Errorf("decrypting secrets: %w", err)
	}
	return decrypted, nil
}

// WriteEncryptedFile writes an encrypted .env to path with restricted permissions.
func WriteEncryptedFile(path string, secrets map[string]string, passphrase string) error {
	encrypted, err := EncryptMap(secrets, passphrase)
	if err != nil {
		return err
	}
	content := Format(encrypted)
	return os.WriteFile(path, []byte(content), 0600)
}
