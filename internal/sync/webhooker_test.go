package sync_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/envfile"
	syncer "github.com/your-org/vaultpull/internal/sync"
)

func TestDispatchSyncEvent_SendsPayload(t *testing.T) {
	var received envfile.WebhookEvent
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	var buf bytes.Buffer
	cfg := syncer.WebhookDispatchConfig{
		Webhook: envfile.WebhookConfig{URL: srv.URL},
		Source:  "vault/prod",
		Out:     &buf,
	}
	diff := []envfile.DiffEntry{{Key: "FOO", Action: "added"}}

	if err := syncer.DispatchSyncEvent(cfg, nil, diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Event != "sync" {
		t.Errorf("expected event=sync, got %q", received.Event)
	}
	if received.Source != "vault/prod" {
		t.Errorf("expected source vault/prod, got %q", received.Source)
	}
	if !strings.Contains(buf.String(), "dispatched sync event") {
		t.Errorf("expected log output, got: %q", buf.String())
	}
}

func TestDispatchSyncEvent_NoURL_NoOp(t *testing.T) {
	cfg := syncer.WebhookDispatchConfig{Source: "test"}
	if err := syncer.DispatchSyncEvent(cfg, nil, nil); err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
}

func TestDispatchSyncEvent_ServerError_ReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	cfg := syncer.WebhookDispatchConfig{
		Webhook: envfile.WebhookConfig{URL: srv.URL},
	}
	err := syncer.DispatchSyncEvent(cfg, nil, nil)
	if err == nil {
		t.Fatal("expected error for bad gateway")
	}
}
