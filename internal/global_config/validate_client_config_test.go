package global_config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     GlobalConfig
		wantErr string
	}{
		{
			name: "ok minimal",
			cfg:  GlobalConfig{},
		},
		{
			name: "ok with optional values",
			cfg: GlobalConfig{
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
			name: "negative http timeout",
			cfg: GlobalConfig{
				HTTPTimeout: -1 * time.Second,
			},
			wantErr: "http-timeout must be >= 0",
		},
		{
			name: "negative context timeout",
			cfg: GlobalConfig{
				ContextTimeout: -1 * time.Second,
			},
			wantErr: "context-timeout must be >= 0",
		},
		{
			name: "zero context timeout is allowed",
			cfg: GlobalConfig{
				ContextTimeout: 0,
			},
		},
		{
			name: "positive context timeout is allowed",
			cfg: GlobalConfig{
				ContextTimeout: 10 * time.Second,
			},
		},
		{
			name: "retries less than minus one",
			cfg: GlobalConfig{
				MaxRetries: -2,
			},
			wantErr: "retries must be >= -1",
		},
		{
			name: "retries minus one is allowed",
			cfg: GlobalConfig{
				MaxRetries: -1,
			},
		},
		{
			name: "negative initial backoff",
			cfg: GlobalConfig{
				InitialBackoff: -1 * time.Second,
			},
			wantErr: "backoff-initial must be >= 0",
		},
		{
			name: "negative max backoff",
			cfg: GlobalConfig{
				MaxBackoff: -1 * time.Second,
			},
			wantErr: "backoff-max must be >= 0",
		},
		{
			name: "negative poll initial wait",
			cfg: GlobalConfig{
				PollInitialWait: -1 * time.Second,
			},
			wantErr: "poll-initial-wait must be >= 0",
		},
		{
			name: "negative poll max wait",
			cfg: GlobalConfig{
				PollMaxWait: -1 * time.Second,
			},
			wantErr: "poll-max-wait must be >= 0",
		},
		{
			name: "max backoff less than initial backoff",
			cfg: GlobalConfig{
				InitialBackoff: 5 * time.Second,
				MaxBackoff:     1 * time.Second,
			},
			wantErr: "backoff-max must be >= backoff-initial",
		},
		{
			name: "poll max wait less than poll initial wait",
			cfg: GlobalConfig{
				PollInitialWait: 10 * time.Second,
				PollMaxWait:     2 * time.Second,
			},
			wantErr: "poll-max-wait must be >= poll-initial-wait",
		},
		{
			name: "equal backoff bounds are allowed",
			cfg: GlobalConfig{
				InitialBackoff: 2 * time.Second,
				MaxBackoff:     2 * time.Second,
			},
		},
		{
			name: "equal poll bounds are allowed",
			cfg: GlobalConfig{
				PollInitialWait: 3 * time.Second,
				PollMaxWait:     3 * time.Second,
			},
		},
		{
			name: "zero backoff pair is allowed",
			cfg: GlobalConfig{
				InitialBackoff: 0,
				MaxBackoff:     0,
			},
		},
		{
			name: "zero poll pair is allowed",
			cfg: GlobalConfig{
				PollInitialWait: 0,
				PollMaxWait:     0,
			},
		},
		{
			name: "only initial backoff set is allowed",
			cfg: GlobalConfig{
				InitialBackoff: 1 * time.Second,
			},
		},
		{
			name: "only max backoff set is allowed",
			cfg: GlobalConfig{
				MaxBackoff: 5 * time.Second,
			},
		},
		{
			name: "only poll initial wait set is allowed",
			cfg: GlobalConfig{
				PollInitialWait: 1 * time.Second,
			},
		},
		{
			name: "only poll max wait set is allowed",
			cfg: GlobalConfig{
				PollMaxWait: 10 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.cfg.Validate()

			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.wantErr)
		})
	}
}

func TestValidateToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     GlobalConfig
		wantErr string
	}{
		{
			name: "token present",
			cfg: GlobalConfig{
				Token: "token",
			},
		},
		{
			name:    "token missing",
			cfg:     GlobalConfig{},
			wantErr: "token is required",
		},
		{
			name: "token whitespace only",
			cfg: GlobalConfig{
				Token: "   ",
			},
			wantErr: "token is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.cfg.ValidateToken()

			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.wantErr)
		})
	}
}

func TestValidateProjectAccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     GlobalConfig
		wantErr string
	}{
		{
			name: "token and project id present",
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
			wantErr: "token is required",
		},
		{
			name: "missing project id",
			cfg: GlobalConfig{
				Token: "token",
			},
			wantErr: "project-id is required",
		},
		{
			name: "project id whitespace only",
			cfg: GlobalConfig{
				Token:     "token",
				ProjectID: "   ",
			},
			wantErr: "project-id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.cfg.ValidateProjectAccess()

			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.wantErr)
		})
	}
}

func TestValidateToken_ValidatesBaseConfigFirst(t *testing.T) {
	t.Parallel()

	cfg := GlobalConfig{
		Token:       "token",
		HTTPTimeout: -time.Second,
	}

	require.EqualError(
		t,
		cfg.ValidateToken(),
		"http-timeout must be >= 0",
	)
}

func TestValidateProjectAccess_ValidatesTokenFirst(t *testing.T) {
	t.Parallel()

	cfg := GlobalConfig{
		ProjectID: "project-id",
	}

	require.EqualError(
		t,
		cfg.ValidateProjectAccess(),
		"token is required",
	)
}
