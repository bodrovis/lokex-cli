package global_config

import (
	"time"

	"github.com/spf13/pflag"
)

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
	fs.DurationVar(&cfg.ContextTimeout, "context-timeout", 150*time.Second, "Overall command timeout (e.g. 30s, 2m). 0 disables the timeout")
}
