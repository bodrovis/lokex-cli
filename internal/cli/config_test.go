package cli

import (
	"testing"
	"time"

	"github.com/spf13/pflag"
)

func TestValidateClientConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     GlobalConfig
		wantErr string
	}{
		{
			name: "ok",
			cfg: GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
		},
		{
			name: "missing token",
			cfg: GlobalConfig{
				ProjectID: "project-id",
			},
			wantErr: "--token is required",
		},
		{
			name: "missing project id",
			cfg: GlobalConfig{
				Token: "token",
			},
			wantErr: "--project-id is required",
		},
		{
			name: "whitespace only values",
			cfg: GlobalConfig{
				Token:     "   ",
				ProjectID: "   ",
			},
			wantErr: "--token is required",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.cfg.ValidateClientConfig()
			if tt.wantErr == "" && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
				}
			}
		})
	}
}

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
