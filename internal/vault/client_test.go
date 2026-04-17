package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/your-org/vaultpull/internal/vault"
)

func newMockVault(t *testing.T, payload map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": payload})
	}))
}

func TestReadSecret_KVv1(t *testing.T) {
	srv := newMockVault(t, map[string]interface{}{"DB_URL": "postgres://localhost/test"})
	defer srv.Close()

	c, err := vault.NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	secrets, err := c.ReadSecret("secret/myapp")
	if err != nil {
		t.Fatalf("ReadSecret: %v", err)
	}

	if secrets["DB_URL"] != "postgres://localhost/test" {
		t.Errorf("expected DB_URL, got %v", secrets)
	}
}

func TestNewClient_InvalidAddr(t *testing.T) {
	_, err := vault.NewClient("://bad-addr", "token")
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}
