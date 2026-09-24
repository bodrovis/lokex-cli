package global_config

import (
	"time"

	"github.com/bodrovis/lokex-cli/internal/params"
)

func paramSpecs(userAgent string) []params.Spec {
	return []params.Spec{
		{
			Name:  "token",
			Kind:  params.KindString,
			Usage: "Lokalise API token",
		},
		{
			Name:  "project-id",
			Kind:  params.KindString,
			Usage: "Lokalise project ID",
		},
		{
			Name:  "base-url",
			Kind:  params.KindString,
			Usage: "Override Lokalise API base URL",
		},
		{
			Name:    "user-agent",
			Kind:    params.KindString,
			Usage:   "User-Agent header",
			Default: userAgent,
		},
		{
			Name:    "http-timeout",
			Kind:    params.KindDuration,
			Usage:   "HTTP client timeout (e.g. 30s, 1m). 0 means library default",
			Default: time.Duration(0),
		},
		{
			Name:    "retries",
			Kind:    params.KindInt,
			Usage:   "Number of retries after the first attempt. -1 means library default",
			Default: -1,
		},
		{
			Name:    "backoff-initial",
			Kind:    params.KindDuration,
			Usage:   "Initial retry backoff (e.g. 400ms, 1s). 0 means library default",
			Default: time.Duration(0),
		},
		{
			Name:    "backoff-max",
			Kind:    params.KindDuration,
			Usage:   "Maximum retry backoff (e.g. 5s, 10s). 0 means library default",
			Default: time.Duration(0),
		},
		{
			Name:    "poll-initial-wait",
			Kind:    params.KindDuration,
			Usage:   "Initial wait between polling rounds (e.g. 1s, 2s). 0 means library default",
			Default: time.Duration(0),
		},
		{
			Name:    "poll-max-wait",
			Kind:    params.KindDuration,
			Usage:   "Maximum total wait for polling (e.g. 120s, 5m). 0 means library default",
			Default: time.Duration(0),
		},
		{
			Name:    "context-timeout",
			Kind:    params.KindDuration,
			Usage:   "Overall command timeout (e.g. 30s, 2m). 0 disables the timeout",
			Default: 150 * time.Second,
		},
	}
}
