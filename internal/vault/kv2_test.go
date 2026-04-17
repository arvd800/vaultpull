package vault

import (
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func TestInjectDataSegment(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"secret/myapp", "secret/data/myapp"},
		{"secret/data/myapp", "secret/data/myapp"},
		{"kv/prod/db", "kv/data/prod/db"},
		{"nomount", "nomount"},
	}
	for _, tc := range cases {
		got := injectDataSegment(tc.in)
		if got != tc.want {
			t.Errorf("injectDataSegment(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestReadSecretKVv2(t *testing.T) {
	srv := newMockVault(t, map[string]interface{}{
		"data": map[string]interface{}{
			"DB_PASSWORD": "s3cr3t",
			"API_KEY":     "abc123",
		},
	})
	defer srv.Close()

	cfg := vaultapi.DefaultConfig()
	cfg.Address = srv.URL
	raw, _ := vaultapi.NewClient(cfg)
	c := &Client{logical: raw.Logical()}

	// newMockVault registers the path without the data segment, so we pass the
	// raw path and rely on injectDataSegment to build the correct request path.
	secrets, err := c.ReadSecretKVv2("secret/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if secrets["DB_PASSWORD"] != "s3cr3t" {
		t.Errorf("DB_PASSWORD = %q, want %q", secrets["DB_PASSWORD"], "s3cr3t")
	}
	if secrets["API_KEY"] != "abc123" {
		t.Errorf("API_KEY = %q, want %q", secrets["API_KEY"], "abc123")
	}
}

func TestReadSecretKVv2_NotFound(t *testing.T) {
	srv := newMockVault(t, map[string]interface{}{})
	defer srv.Close()

	cfg := vaultapi.DefaultConfig()
	cfg.Address = srv.URL
	raw, _ := vaultapi.NewClient(cfg)
	c := &Client{logical: raw.Logical()}

	_, err := c.ReadSecretKVv2("secret/missing")
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}
