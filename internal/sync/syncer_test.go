package sync

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMockVault(t *testing.T, secret string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(secret))
	}))
}

func TestSyncer_Run(t *testing.T) {
	kvv1Payload := `{"data":{"API_KEY":"abc123","DB_PASS":"secret"}}`
	server := newMockVault(t, kvv1Payload)
	defer server.Close()

	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, ".env")

	s, err := New(server.URL, "fake-token", "secret/myapp", outPath)
	require.NoError(t, err)

	err = s.Run()
	require.NoError(t, err)

	data, err := os.ReadFile(outPath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "API_KEY")
	assert.Contains(t, string(data), "DB_PASS")
}

func TestSyncer_Run_BadPath(t *testing.T) {
	kvv1Payload := `{"data":{"KEY":"val"}}`
	server := newMockVault(t, kvv1Payload)
	defer server.Close()

	s, err := New(server.URL, "fake-token", "secret/myapp", "/nonexistent/dir/.env")
	require.NoError(t, err)

	err = s.Run()
	assert.Error(t, err)
}
