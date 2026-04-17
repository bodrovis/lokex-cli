package global_config

import (
	"testing"
	"time"
)

func TestNewClient_Minimal(t *testing.T) {
	t.Parallel()

	cfg := &GlobalConfig{
		Token:     "test-token",
		ProjectID: "test-project",
	}

	client, err := cfg.NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if client.Token != "test-token" {
		t.Fatalf("unexpected token: %q", client.Token)
	}
	if client.ProjectID != "test-project" {
		t.Fatalf("unexpected project id: %q", client.ProjectID)
	}
	if client.HTTPClient == nil {
		t.Fatal("expected non-nil HTTP client")
	}
}

func TestNewClient_WithOverrides(t *testing.T) {
	t.Parallel()

	cfg := &GlobalConfig{
		Token:           "test-token",
		ProjectID:       "test-project",
		BaseURL:         "https://example.com/api/",
		UserAgent:       "lokex-cli/test",
		HTTPTimeout:     45 * time.Second,
		MaxRetries:      5,
		InitialBackoff:  500 * time.Millisecond,
		MaxBackoff:      10 * time.Second,
		PollInitialWait: 2 * time.Second,
		PollMaxWait:     30 * time.Second,
	}

	client, err := cfg.NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if client.BaseURL != "https://example.com/api/" {
		t.Fatalf("unexpected base url: %q", client.BaseURL)
	}
	if client.UserAgent != "lokex-cli/test" {
		t.Fatalf("unexpected user agent: %q", client.UserAgent)
	}
	if client.HTTPClient == nil {
		t.Fatal("expected non-nil HTTP client")
	}
	if client.HTTPClient.Timeout != 45*time.Second {
		t.Fatalf("unexpected http timeout: %v", client.HTTPClient.Timeout)
	}
	if client.MaxRetries != 5 {
		t.Fatalf("unexpected retries: %d", client.MaxRetries)
	}
	if client.InitialBackoff != 500*time.Millisecond {
		t.Fatalf("unexpected initial backoff: %v", client.InitialBackoff)
	}
	if client.MaxBackoff != 10*time.Second {
		t.Fatalf("unexpected max backoff: %v", client.MaxBackoff)
	}
	if client.PollInitialWait != 2*time.Second {
		t.Fatalf("unexpected poll initial wait: %v", client.PollInitialWait)
	}
	if client.PollMaxWait != 30*time.Second {
		t.Fatalf("unexpected poll max wait: %v", client.PollMaxWait)
	}
}
