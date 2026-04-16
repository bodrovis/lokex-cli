package cli

import (
	"errors"
	"strings"
	"time"

	lokexclient "github.com/bodrovis/lokex/v2/client"
	"github.com/spf13/pflag"
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

func BindPersistentFlags(fs *pflag.FlagSet, cfg *GlobalConfig) {
	fs.StringVar(&cfg.Token, "token", "", "Lokalise API token")
	fs.StringVar(&cfg.ProjectID, "project-id", "", "Lokalise project ID")
	fs.StringVar(&cfg.BaseURL, "base-url", "", "Override Lokalise API base URL")
	fs.StringVar(&cfg.UserAgent, "user-agent", cfg.UserAgent, "User-Agent header")
	fs.DurationVar(&cfg.HTTPTimeout, "http-timeout", 0, "HTTP client timeout (e.g. 30s, 1m). 0 means library default")
	fs.IntVar(&cfg.MaxRetries, "retries", -1, "Number of retries after the first attempt. -1 means library default")
	fs.DurationVar(&cfg.InitialBackoff, "backoff-initial", 0, "Initial retry backoff (e.g. 400ms, 1s). 0 means library default")
	fs.DurationVar(&cfg.MaxBackoff, "backoff-max", 0, "Maximum retry backoff (e.g. 5s, 10s). 0 means library default")
	fs.DurationVar(&cfg.PollInitialWait, "poll-initial-wait", 0, "Initial wait between polling rounds (e.g. 1s, 2s). 0 means library default")
	fs.DurationVar(&cfg.PollMaxWait, "poll-max-wait", 0, "Maximum total wait for polling (e.g. 120s, 5m). 0 means library default")
}

func (cfg *GlobalConfig) ValidateClientConfig() error {
	if strings.TrimSpace(cfg.Token) == "" {
		return errors.New("--token is required")
	}
	if strings.TrimSpace(cfg.ProjectID) == "" {
		return errors.New("--project-id is required")
	}
	return nil
}

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
