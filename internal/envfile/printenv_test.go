package envfile

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintEnv_BasicOutput(t *testing.T) {
	var buf bytes.Buffer
	secrets := map[string]string{"APP_NAME": "myapp", "PORT": "8080"}
	if err := PrintEnv(&buf, secrets, nil); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("missing APP_NAME line in %q", out)
	}
	if !strings.Contains(out, "PORT=8080") {
		t.Errorf("missing PORT line in %q", out)
	}
}

func TestPrintEnv_RedactsSecrets(t *testing.T) {
	var buf bytes.Buffer
	secrets := map[string]string{"API_TOKEN": "supersecret", "APP": "ok"}
	err := PrintEnv(&buf, secrets, &PrintOptions{Redact: true})
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if strings.Contains(out, "supersecret") {
		t.Error("sensitive value should have been redacted")
	}
	if !strings.Contains(out, "APP=ok") {
		t.Error("non-sensitive value should be present")
	}
}

func TestPrintEnv_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	secrets := map[string]string{"ZZZ": "last", "AAA": "first"}
	PrintEnv(&buf, secrets, nil)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 || !strings.HasPrefix(lines[0], "AAA") {
		t.Errorf("expected sorted output, got %v", lines)
	}
}

func TestPrintEnv_QuotesValuesWithSpaces(t *testing.T) {
	var buf bytes.Buffer
	secrets := map[string]string{"GREETING": "hello world"}
	PrintEnv(&buf, secrets, nil)
	out := buf.String()
	if !strings.Contains(out, `"hello world"`) {
		t.Errorf("expected quoted value in %q", out)
	}
}
