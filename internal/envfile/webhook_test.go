package envfile_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/your-org/vaultpull/internal/envfile"
)

func TestSendWebhook_Success(t *testing.T) {
	var received envfile.WebhookEvent

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cfg := envfile.WebhookConfig{URL: srv.URL, TimeoutSec: 5}
	event := envfile.WebhookEvent{
		Event:  "sync",
		Source: "test",
	}

	if err := envfile.SendWebhook(cfg, event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Event != "sync" {
		t.Errorf("expected event=sync, got %q", received.Event)
	}
	if received.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestSendWebhook_EmptyURL_NoOp(t *testing.T) {
	cfg := envfile.WebhookConfig{}
	if err := envfile.SendWebhook(cfg, envfile.WebhookEvent{Event: "sync"}); err != nil {
		t.Fatalf("expected nil error for empty URL, got: %v", err)
	}
}

func TestSendWebhook_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	cfg := envfile.WebhookConfig{URL: srv.URL}
	err := envfile.SendWebhook(cfg, envfile.WebhookEvent{Event: "sync"})
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestSendWebhook_SecretHeader(t *testing.T) {
	var gotSecret string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSecret = r.Header.Get("X-Vaultpull-Secret")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cfg := envfile.WebhookConfig{URL: srv.URL, Secret: "mysecret"}
	if err := envfile.SendWebhook(cfg, envfile.WebhookEvent{Event: "sync"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotSecret != "mysecret" {
		t.Errorf("expected secret header mysecret, got %q", gotSecret)
	}
}

func TestSendWebhook_InvalidURL(t *testing.T) {
	cfg := envfile.WebhookConfig{URL: "://bad-url"}
	err := envfile.SendWebhook(cfg, envfile.WebhookEvent{Event: "sync"})
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
