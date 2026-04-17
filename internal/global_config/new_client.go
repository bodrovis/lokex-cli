package global_config

import (
	"strings"

	lokexclient "github.com/bodrovis/lokex/v2/client"
)

func (cfg *GlobalConfig) NewClient() (*lokexclient.Client, error) {
	opts := make([]lokexclient.Option, 0, 6)

	if strings.TrimSpace(cfg.BaseURL) != "" {
		opts = append(opts, lokexclient.WithBaseURL(cfg.BaseURL))
	}
	if strings.TrimSpace(cfg.UserAgent) != "" {
		opts = append(opts, lokexclient.WithUserAgent(cfg.UserAgent))
	}
	if cfg.HTTPTimeout != 0 {
		opts = append(opts, lokexclient.WithHTTPTimeout(cfg.HTTPTimeout))
	}
	if cfg.MaxRetries >= 0 {
		opts = append(opts, lokexclient.WithMaxRetries(cfg.MaxRetries))
	}
	if cfg.InitialBackoff != 0 || cfg.MaxBackoff != 0 {
		opts = append(opts, lokexclient.WithBackoff(cfg.InitialBackoff, cfg.MaxBackoff))
	}
	if cfg.PollInitialWait != 0 || cfg.PollMaxWait != 0 {
		opts = append(opts, lokexclient.WithPollWait(cfg.PollInitialWait, cfg.PollMaxWait))
	}

	return lokexclient.NewClient(cfg.Token, cfg.ProjectID, opts...)
}
