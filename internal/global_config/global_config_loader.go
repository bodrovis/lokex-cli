package global_config

import (
	"fmt"

	vh "github.com/bodrovis/lokex-cli/internal/viper_helpers"
	"github.com/spf13/cobra"
)

type LoadOptions struct {
	ConfigFile string
	EnvPrefix  string
}

func LoadGlobalConfigInput(userAgent string, opts LoadOptions) (*GlobalConfigInput, error) {
	v := vh.NewConfigViper(opts.ConfigFile, opts.EnvPrefix)

	v.SetDefault("user-agent", userAgent)

	if err := vh.ReadOptionalConfig(v, opts.ConfigFile); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var input GlobalConfigInput
	if err := v.Unmarshal(&input); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	return &input, nil
}

func ApplyGlobalInput(cmd *cobra.Command, cfg *GlobalConfig, input *GlobalConfigInput) {
	if input == nil {
		return
	}

	if !cmd.Flags().Changed("token") && input.Token != nil {
		cfg.Token = *input.Token
	}

	if !cmd.Flags().Changed("project-id") && input.ProjectID != nil {
		cfg.ProjectID = *input.ProjectID
	}

	if !cmd.Flags().Changed("base-url") && input.BaseURL != nil {
		cfg.BaseURL = *input.BaseURL
	}

	if !cmd.Flags().Changed("user-agent") && input.UserAgent != nil {
		cfg.UserAgent = *input.UserAgent
	}

	if !cmd.Flags().Changed("http-timeout") && input.HTTPTimeout != nil {
		cfg.HTTPTimeout = *input.HTTPTimeout
	}

	if !cmd.Flags().Changed("retries") && input.MaxRetries != nil {
		cfg.MaxRetries = *input.MaxRetries
	}

	if !cmd.Flags().Changed("backoff-initial") && input.InitialBackoff != nil {
		cfg.InitialBackoff = *input.InitialBackoff
	}

	if !cmd.Flags().Changed("backoff-max") && input.MaxBackoff != nil {
		cfg.MaxBackoff = *input.MaxBackoff
	}

	if !cmd.Flags().Changed("poll-initial-wait") && input.PollInitialWait != nil {
		cfg.PollInitialWait = *input.PollInitialWait
	}

	if !cmd.Flags().Changed("poll-max-wait") && input.PollMaxWait != nil {
		cfg.PollMaxWait = *input.PollMaxWait
	}

	if !cmd.Flags().Changed("context-timeout") && input.ContextTimeout != nil {
		cfg.ContextTimeout = *input.ContextTimeout
	}
}
