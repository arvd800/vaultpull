package envfile

import (
	"strings"
	"testing"
)

func TestObfuscateMap_ReplacesValues(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "abc123",
	}

	obf, lookup, err := ObfuscateMap(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for k, token := range obf {
		if token == secrets[k] {
			t.Errorf("key %q: token should not equal original value", k)
		}
		if !IsObfuscatedToken(token) {
			t.Errorf("key %q: token %q does not match expected format", k, token)
		}
		if orig, ok := lookup[token]; !ok || orig != secrets[k] {
			t.Errorf("key %q: lookup mismatch: got %q, want %q", k, orig, secrets[k])
		}
	}
}

func TestObfuscateMap_ProducesUniqueTokens(t *testing.T) {
	secrets := map[string]string{"A": "val", "B": "val"}
	obf, _, err := ObfuscateMap(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obf["A"] == obf["B"] {
		t.Error("expected unique tokens for different keys even with same value")
	}
}

func TestDeobfuscateMap_RoundTrip(t *testing.T) {
	secrets := map[string]string{
		"SECRET_ONE": "hello",
		"SECRET_TWO": "world",
	}

	obf, lookup, err := ObfuscateMap(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	restored := DeobfuscateMap(obf, lookup)
	for k, want := range secrets {
		if got := restored[k]; got != want {
			t.Errorf("key %q: got %q, want %q", k, got, want)
		}
	}
}

func TestDeobfuscateMap_UnknownTokenLeftAsIs(t *testing.T) {
	obf := map[string]string{"KEY": "vp_unknowntoken0000000000000000000"}
	lookup := map[string]string{}
	out := DeobfuscateMap(obf, lookup)
	if out["KEY"] != obf["KEY"] {
		t.Errorf("expected unknown token to be left as-is")
	}
}

func TestObfuscateMap_EmptyMap(t *testing.T) {
	obf, lookup, err := ObfuscateMap(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(obf) != 0 || len(lookup) != 0 {
		t.Error("expected empty outputs for empty input")
	}
}

func TestIsObfuscatedToken(t *testing.T) {
	valid := "vp_" + strings.Repeat("a", 32)
	if !IsObfuscatedToken(valid) {
		t.Errorf("expected %q to be recognised as obfuscated token", valid)
	}
	if IsObfuscatedToken("plaintext") {
		t.Error("expected plain string not to be recognised as token")
	}
}
