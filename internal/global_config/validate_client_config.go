package global_config

import (
	"errors"
	"strings"
)

func (cfg *GlobalConfig) Validate() error {
	if cfg.HTTPTimeout < 0 {
		return errors.New("http-timeout must be >= 0")
	}

	if cfg.MaxRetries < -1 {
		return errors.New("retries must be >= -1")
	}

	if cfg.InitialBackoff < 0 {
		return errors.New("backoff-initial must be >= 0")
	}

	if cfg.MaxBackoff < 0 {
		return errors.New("backoff-max must be >= 0")
	}

	if cfg.PollInitialWait < 0 {
		return errors.New("poll-initial-wait must be >= 0")
	}

	if cfg.PollMaxWait < 0 {
		return errors.New("poll-max-wait must be >= 0")
	}

	if cfg.ContextTimeout < 0 {
		return errors.New("context-timeout must be >= 0")
	}

	if cfg.InitialBackoff > 0 &&
		cfg.MaxBackoff > 0 &&
		cfg.MaxBackoff < cfg.InitialBackoff {
		return errors.New("backoff-max must be >= backoff-initial")
	}

	if cfg.PollInitialWait > 0 &&
		cfg.PollMaxWait > 0 &&
		cfg.PollMaxWait < cfg.PollInitialWait {
		return errors.New("poll-max-wait must be >= poll-initial-wait")
	}

	return nil
}

func (c *GlobalConfig) ValidateToken() error {
	if err := c.Validate(); err != nil {
		return err
	}

	if strings.TrimSpace(c.Token) == "" {
		return errors.New("token is required")
	}

	return nil
}

func (c *GlobalConfig) ValidateProjectAccess() error {
	if err := c.ValidateToken(); err != nil {
		return err
	}

	if strings.TrimSpace(c.ProjectID) == "" {
		return errors.New("project-id is required")
	}

	return nil
}
