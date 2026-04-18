package global_config

import (
	"errors"
	"strings"
)

func (cfg *GlobalConfig) ValidateClientConfig() error {
	if strings.TrimSpace(cfg.Token) == "" {
		return errors.New("--token is required")
	}
	if strings.TrimSpace(cfg.ProjectID) == "" {
		return errors.New("--project-id is required")
	}
	if cfg.HTTPTimeout < 0 {
		return errors.New("--http-timeout must be >= 0")
	}
	if cfg.MaxRetries < -1 {
		return errors.New("--retries must be >= -1")
	}
	if cfg.InitialBackoff < 0 {
		return errors.New("--backoff-initial must be >= 0")
	}
	if cfg.MaxBackoff < 0 {
		return errors.New("--backoff-max must be >= 0")
	}
	if cfg.PollInitialWait < 0 {
		return errors.New("--poll-initial-wait must be >= 0")
	}
	if cfg.PollMaxWait < 0 {
		return errors.New("--poll-max-wait must be >= 0")
	}
	if cfg.ContextTimeout < 0 {
		return errors.New("--context-timeout must be >= 0")
	}
	if cfg.InitialBackoff > 0 && cfg.MaxBackoff > 0 && cfg.MaxBackoff < cfg.InitialBackoff {
		return errors.New("--backoff-max must be >= --backoff-initial")
	}
	if cfg.PollInitialWait > 0 && cfg.PollMaxWait > 0 && cfg.PollMaxWait < cfg.PollInitialWait {
		return errors.New("--poll-max-wait must be >= --poll-initial-wait")
	}
	return nil
}
