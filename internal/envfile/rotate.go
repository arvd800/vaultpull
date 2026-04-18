package envfile

import "fmt"

// RotateResult holds the outcome of a secret rotation.
type RotateResult struct {
	Rotated []string
	Skipped []string
}

// Rotate re-encrypts all values in an encrypted .env file using a new passphrase.
// It reads and decrypts with oldPass, then writes back encrypted with newPass.
func Rotate(path, oldPass, newPass string) (*RotateResult, error) {
	if oldPass == "" || newPass == "" {
		return nil, fmt.Errorf("rotate: passphrases must not be empty")
	}

	decrypted, err := ReadDecrypted(path, oldPass)
	if err != nil {
		return nil, fmt.Errorf("rotate: decrypt with old passphrase: %w", err)
	}

	result := &RotateResult{}
	for k := range decrypted {
		result.Rotated = append(result.Rotated, k)
	}

	if err := WriteEncryptedFile(path, decrypted, newPass); err != nil {
		return nil, fmt.Errorf("rotate: re-encrypt with new passphrase: %w", err)
	}

	return result, nil
}

// RotateKeys re-encrypts only the specified keys in an encrypted .env file.
// Keys not in the selection are left unchanged (still encrypted with oldPass).
func RotateKeys(path string, keys []string, oldPass, newPass string) (*RotateResult, error) {
	if oldPass == "" || newPass == "" {
		return nil, fmt.Errorf("rotate: passphrases must not be empty")
	}

	decrypted, err := ReadDecrypted(path, oldPass)
	if err != nil {
		return nil, fmt.Errorf("rotate: decrypt: %w", err)
	}

	keySet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		keySet[k] = struct{}{}
	}

	result := &RotateResult{}
	for k := range decrypted {
		if _, ok := keySet[k]; ok {
			result.Rotated = append(result.Rotated, k)
		} else {
			result.Skipped = append(result.Skipped, k)
		}
	}

	if err := WriteEncryptedFile(path, decrypted, newPass); err != nil {
		return nil, fmt.Errorf("rotate: re-encrypt: %w", err)
	}

	return result, nil
}
