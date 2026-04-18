package envfile

import "fmt"

// EncryptMap encrypts all values in a secrets map using the given passphrase.
// Returns a new map with encrypted values.
func EncryptMap(secrets map[string]string, passphrase string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		enc, err := Encrypt(v, passphrase)
		if err != nil {
			return nil, fmt.Errorf("encrypting key %q: %w", k, err)
		}
		out[k] = enc
	}
	return out, nil
}

// DecryptMap decrypts all values in a secrets map using the given passphrase.
// Returns a new map with plaintext values.
func DecryptMap(secrets map[string]string, passphrase string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		dec, err := Decrypt(v, passphrase)
		if err != nil {
			return nil, fmt.Errorf("decrypting key %q: %w", k, err)
		}
		out[k] = dec
	}
	return out, nil
}
