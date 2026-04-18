package envfile_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/envfile"
)

// TestPrintEnv_WithRedact_Integration exercises Redact + PrintEnv together.
func TestPrintEnv_WithRedact_Integration(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "topsecret",
		"DB_HOST":     "localhost",
		"API_KEY":     "key-abc",
		"APP_ENV":     "production",
	}

	var buf bytes.Buffer
	err := envfile.PrintEnv(&buf, secrets, &envfile.PrintOptions{Redact: true})
	if err != nil {
		t.Fatalf("PrintEnv error: %v", err)
	}

	out := buf.String()

	for _, sensitive := range []string{"topsecret", "key-abc"} {
		if strings.Contains(out, sensitive) {
			t.Errorf("sensitive value %q should not appear in output", sensitive)
		}
	}

	for _, plain := range []string{"DB_HOST=localhost", "APP_ENV=production"} {
		if !strings.Contains(out, plain) {
			t.Errorf("expected %q in output, got:\n%s", plain, out)
		}
	}
}
