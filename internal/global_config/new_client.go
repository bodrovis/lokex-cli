package global_config

import (
	"strings"

	lokexclient "github.com/bodrovis/lokex/v2/client"
)

func (cfg *GlobalConfig) NewClient() (*lokexclient.Client, error) {
	token := strings.TrimSpace(cfg.Token)
	projectID := strings.TrimSpace(cfg.ProjectID)
	baseURL := strings.TrimSpace(cfg.BaseURL)
	userAgent := strings.TrimSpace(cfg.UserAgent)

	opts := make([]lokexclient.Option, 0, 6)

	if baseURL != "" {
		opts = append(opts, lokexclient.WithBaseURL(baseURL))
	}
	if userAgent != "" {
		opts = append(opts, lokexclient.WithUserAgent(userAgent))
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

	return lokexclient.NewClient(token, projectID, opts...)
}
