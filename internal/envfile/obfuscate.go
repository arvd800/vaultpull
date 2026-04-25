package envfile

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

// ObfuscateRule defines how a key should be obfuscated.
type ObfuscateRule struct {
	// Token is the stable opaque token replacing the real value.
	Token string `json:"token"`
	// OriginalKey is the source key name.
	OriginalKey string `json:"original_key"`
}

// ObfuscateMap replaces secret values with stable opaque tokens.
// It returns the obfuscated map and a lookup table mapping token -> original value.
func ObfuscateMap(secrets map[string]string) (obfuscated map[string]string, lookup map[string]string, err error) {
	obfuscated = make(map[string]string, len(secrets))
	lookup = make(map[string]string, len(secrets))

	for k, v := range secrets {
		token, genErr := generateToken()
		if genErr != nil {
			return nil, nil, fmt.Errorf("obfuscate: generate token for %q: %w", k, genErr)
		}
		obfuscated[k] = token
		lookup[token] = v
	}
	return obfuscated, lookup, nil
}

// DeobfuscateMap reverses ObfuscateMap using the provided lookup table.
// Keys whose token is not found in the lookup are left as-is.
func DeobfuscateMap(obfuscated map[string]string, lookup map[string]string) map[string]string {
	out := make(map[string]string, len(obfuscated))
	for k, token := range obfuscated {
		if original, ok := lookup[token]; ok {
			out[k] = original
		} else {
			out[k] = token
		}
	}
	return out
}

// IsObfuscatedToken reports whether s looks like a token produced by ObfuscateMap.
func IsObfuscatedToken(s string) bool {
	return strings.HasPrefix(s, "vp_") && len(s) == 35
}

func generateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "vp_" + hex.EncodeToString(b), nil
}
