package global_config

import (
	"time"
)

type GlobalConfigInput struct {
	Token           *string        `mapstructure:"token"`
	ProjectID       *string        `mapstructure:"project-id"`
	BaseURL         *string        `mapstructure:"base-url"`
	UserAgent       *string        `mapstructure:"user-agent"`
	HTTPTimeout     *time.Duration `mapstructure:"http-timeout"`
	MaxRetries      *int           `mapstructure:"retries"`
	InitialBackoff  *time.Duration `mapstructure:"backoff-initial"`
	MaxBackoff      *time.Duration `mapstructure:"backoff-max"`
	PollInitialWait *time.Duration `mapstructure:"poll-initial-wait"`
	PollMaxWait     *time.Duration `mapstructure:"poll-max-wait"`
}
