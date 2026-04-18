package global_config

import (
	"testing"
	"time"
)

func TestValidateClientConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     GlobalConfig
		wantErr string
	}{
		{
			name: "ok minimal",
			cfg: GlobalConfig{
				Token:     "token",
				ProjectID: "project-id",
			},
		},
		{
			name: "ok with optional values",
			cfg: GlobalConfig{
				Token:           "token",
				ProjectID:       "project-id",
				HTTPTimeout:     30 * time.Second,
				MaxRetries:      3,
				InitialBackoff:  500 * time.Millisecond,
				MaxBackoff:      5 * time.Second,
				PollInitialWait: 1 * time.Second,
				PollMaxWait:     30 * time.Second,
				ContextTimeout:  150 * time.Second,
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
		{
			name: "negative http timeout",
			cfg: GlobalConfig{
				Token:       "token",
				ProjectID:   "project-id",
				HTTPTimeout: -1 * time.Second,
			},
			wantErr: "--http-timeout must be >= 0",
		},
		{
			name: "negative context timeout",
			cfg: GlobalConfig{
				Token:          "token",
				ProjectID:      "project-id",
				ContextTimeout: -1 * time.Second,
			},
			wantErr: "--context-timeout must be >= 0",
		},
		{
			name: "zero context timeout is allowed",
			cfg: GlobalConfig{
				Token:          "token",
				ProjectID:      "project-id",
				ContextTimeout: 0,
			},
		},
		{
			name: "positive context timeout is allowed",
			cfg: GlobalConfig{
				Token:          "token",
				ProjectID:      "project-id",
				ContextTimeout: 10 * time.Second,
			},
		},
		{
			name: "retries less than minus one",
			cfg: GlobalConfig{
				Token:      "token",
				ProjectID:  "project-id",
				MaxRetries: -2,
			},
			wantErr: "--retries must be >= -1",
		},
		{
			name: "retries minus one is allowed",
			cfg: GlobalConfig{
				Token:      "token",
				ProjectID:  "project-id",
				MaxRetries: -1,
			},
		},
		{
			name: "negative initial backoff",
			cfg: GlobalConfig{
				Token:          "token",
				ProjectID:      "project-id",
				InitialBackoff: -1 * time.Second,
			},
			wantErr: "--backoff-initial must be >= 0",
		},
		{
			name: "negative max backoff",
			cfg: GlobalConfig{
				Token:      "token",
				ProjectID:  "project-id",
				MaxBackoff: -1 * time.Second,
			},
			wantErr: "--backoff-max must be >= 0",
		},
		{
			name: "negative poll initial wait",
			cfg: GlobalConfig{
				Token:           "token",
				ProjectID:       "project-id",
				PollInitialWait: -1 * time.Second,
			},
			wantErr: "--poll-initial-wait must be >= 0",
		},
		{
			name: "negative poll max wait",
			cfg: GlobalConfig{
				Token:       "token",
				ProjectID:   "project-id",
				PollMaxWait: -1 * time.Second,
			},
			wantErr: "--poll-max-wait must be >= 0",
		},
		{
			name: "max backoff less than initial backoff",
			cfg: GlobalConfig{
				Token:          "token",
				ProjectID:      "project-id",
				InitialBackoff: 5 * time.Second,
				MaxBackoff:     1 * time.Second,
			},
			wantErr: "--backoff-max must be >= --backoff-initial",
		},
		{
			name: "poll max wait less than poll initial wait",
			cfg: GlobalConfig{
				Token:           "token",
				ProjectID:       "project-id",
				PollInitialWait: 10 * time.Second,
				PollMaxWait:     2 * time.Second,
			},
			wantErr: "--poll-max-wait must be >= --poll-initial-wait",
		},
		{
			name: "equal backoff bounds are allowed",
			cfg: GlobalConfig{
				Token:          "token",
				ProjectID:      "project-id",
				InitialBackoff: 2 * time.Second,
				MaxBackoff:     2 * time.Second,
			},
		},
		{
			name: "equal poll bounds are allowed",
			cfg: GlobalConfig{
				Token:           "token",
				ProjectID:       "project-id",
				PollInitialWait: 3 * time.Second,
				PollMaxWait:     3 * time.Second,
			},
		},
		{
			name: "zero backoff pair is allowed",
			cfg: GlobalConfig{
				Token:          "token",
				ProjectID:      "project-id",
				InitialBackoff: 0,
				MaxBackoff:     0,
			},
		},
		{
			name: "zero poll pair is allowed",
			cfg: GlobalConfig{
				Token:           "token",
				ProjectID:       "project-id",
				PollInitialWait: 0,
				PollMaxWait:     0,
			},
		},
		{
			name: "only initial backoff set is allowed",
			cfg: GlobalConfig{
				Token:          "token",
				ProjectID:      "project-id",
				InitialBackoff: 1 * time.Second,
			},
		},
		{
			name: "only max backoff set is allowed",
			cfg: GlobalConfig{
				Token:      "token",
				ProjectID:  "project-id",
				MaxBackoff: 5 * time.Second,
			},
		},
		{
			name: "only poll initial wait set is allowed",
			cfg: GlobalConfig{
				Token:           "token",
				ProjectID:       "project-id",
				PollInitialWait: 1 * time.Second,
			},
		},
		{
			name: "only poll max wait set is allowed",
			cfg: GlobalConfig{
				Token:       "token",
				ProjectID:   "project-id",
				PollMaxWait: 10 * time.Second,
			},
		},
	}

	for _, tt := range tests {
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
