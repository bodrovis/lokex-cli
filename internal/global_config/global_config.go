package global_config

import (
	"time"
)

type GlobalConfig struct {
	Token           string
	ProjectID       string
	BaseURL         string
	UserAgent       string
	HTTPTimeout     time.Duration
	MaxRetries      int
	InitialBackoff  time.Duration
	MaxBackoff      time.Duration
	PollInitialWait time.Duration
	PollMaxWait     time.Duration
}
