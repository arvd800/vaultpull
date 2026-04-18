package vault

import (
	"context"
	"testing"
)

func TestExtractMountPath(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"secret/myapp/prod", "secret"},
		{"/secret/myapp/prod", "secret"},
		{"kv/data/foo", "kv"},
		{"single", "single"},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := extractMountPath(tc.input)
			if got != tc.expected {
				t.Errorf("extractMountPath(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestDetectKVVersion_NotFound(t *testing.T) {
	srv, addr := newMockVault(t)
	_ = srv

	client, err := NewClient(addr, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = client.DetectKVVersion(context.Background(), "nonexistent/path")
	if err == nil {
		t.Error("expected error for unknown mount, got nil")
	}
}
