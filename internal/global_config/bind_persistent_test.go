package global_config

import (
	"testing"
	"time"

	"github.com/spf13/pflag"
)

func TestBindPersistentFlags(t *testing.T) {
	t.Parallel()

	cfg := &GlobalConfig{}
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	BindPersistentFlags(fs, cfg)

	if err := fs.Parse([]string{
		"--token=test-token",
		"--project-id=test-project",
		"--base-url=https://example.com/api/",
		"--user-agent=lokex-cli/test",
		"--http-timeout=45s",
		"--retries=5",
		"--backoff-initial=500ms",
		"--backoff-max=10s",
		"--poll-initial-wait=2s",
		"--poll-max-wait=30s",
		"--context-timeout=100s",
	}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	if cfg.Token != "test-token" {
		t.Fatalf("unexpected token: %q", cfg.Token)
	}
	if cfg.ProjectID != "test-project" {
		t.Fatalf("unexpected project id: %q", cfg.ProjectID)
	}
	if cfg.BaseURL != "https://example.com/api/" {
		t.Fatalf("unexpected base url: %q", cfg.BaseURL)
	}
	if cfg.UserAgent != "lokex-cli/test" {
		t.Fatalf("unexpected user agent: %q", cfg.UserAgent)
	}
	if cfg.HTTPTimeout != 45*time.Second {
		t.Fatalf("unexpected http timeout: %v", cfg.HTTPTimeout)
	}
	if cfg.ContextTimeout != 100*time.Second {
		t.Fatalf("unexpected context timeout: %v", cfg.ContextTimeout)
	}
	if cfg.MaxRetries != 5 {
		t.Fatalf("unexpected retries: %d", cfg.MaxRetries)
	}
	if cfg.InitialBackoff != 500*time.Millisecond {
		t.Fatalf("unexpected initial backoff: %v", cfg.InitialBackoff)
	}
	if cfg.MaxBackoff != 10*time.Second {
		t.Fatalf("unexpected max backoff: %v", cfg.MaxBackoff)
	}
	if cfg.PollInitialWait != 2*time.Second {
		t.Fatalf("unexpected poll initial wait: %v", cfg.PollInitialWait)
	}
	if cfg.PollMaxWait != 30*time.Second {
		t.Fatalf("unexpected poll max wait: %v", cfg.PollMaxWait)
	}
}

func TestBindPersistentFlags_Defaults(t *testing.T) {
	t.Parallel()

	cfg := &GlobalConfig{}
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	BindPersistentFlags(fs, cfg)

	if err := fs.Parse(nil); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	if cfg.Token != "" {
		t.Fatalf("expected empty token, got %q", cfg.Token)
	}
	if cfg.ProjectID != "" {
		t.Fatalf("expected empty project id, got %q", cfg.ProjectID)
	}
	if cfg.BaseURL != "" {
		t.Fatalf("expected empty base url, got %q", cfg.BaseURL)
	}
	if cfg.UserAgent != "" {
		t.Fatalf("expected empty user agent, got %q", cfg.UserAgent)
	}
	if cfg.HTTPTimeout != 0 {
		t.Fatalf("expected zero http timeout, got %v", cfg.HTTPTimeout)
	}
	if cfg.ContextTimeout != 150*time.Second {
		t.Fatalf("expected 150s context timeout, got %v", cfg.ContextTimeout)
	}
	if cfg.MaxRetries != -1 {
		t.Fatalf("expected default retries -1, got %d", cfg.MaxRetries)
	}
	if cfg.InitialBackoff != 0 {
		t.Fatalf("expected zero initial backoff, got %v", cfg.InitialBackoff)
	}
	if cfg.MaxBackoff != 0 {
		t.Fatalf("expected zero max backoff, got %v", cfg.MaxBackoff)
	}
	if cfg.PollInitialWait != 0 {
		t.Fatalf("expected zero poll initial wait, got %v", cfg.PollInitialWait)
	}
	if cfg.PollMaxWait != 0 {
		t.Fatalf("expected zero poll max wait, got %v", cfg.PollMaxWait)
	}
}
