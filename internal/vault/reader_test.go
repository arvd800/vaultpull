package vault

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReadSecretAuto_FallsBackToKVv1(t *testing.T) {
	// Mock server that serves a KVv1-style secret and returns empty mounts
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/sys/mounts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Return empty mounts to trigger fallback
		w.Write([]byte(`{}`))
	})

	mux.HandleFunc("/v1/secret/myapp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"API_KEY":"abc123","DB_PASS":"secret"}}`))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	client, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	got, err := client.ReadSecretAuto(context.Background(), "secret/myapp")
	if err != nil {
		t.Fatalf("ReadSecretAuto: %v", err)
	}

	if got["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", got["API_KEY"])
	}
	if got["DB_PASS"] != "secret" {
		t.Errorf("expected DB_PASS=secret, got %q", got["DB_PASS"])
	}
}
