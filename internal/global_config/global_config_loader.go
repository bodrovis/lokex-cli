package global_config

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bodrovis/lokex-cli/internal/params"
)

func LoadGlobalConfig(
	v *viper.Viper,
	cmd *cobra.Command,
	cfg *GlobalConfig,
) error {
	if v == nil {
		return errors.New("viper is nil")
	}

	if cmd == nil {
		return errors.New("command is nil")
	}

	if cfg == nil {
		return errors.New("global config is nil")
	}

	specs := paramSpecs(cfg.UserAgent)

	params.ApplyDefaults(v, "", specs)

	if err := params.BindEnv(
		v,
		"",
		specs,
	); err != nil {
		return fmt.Errorf("bind global environment: %w", err)
	}

	if err := params.ApplyChangedFlags(
		v,
		cmd,
		"",
		specs,
	); err != nil {
		return fmt.Errorf("apply global flags: %w", err)
	}

	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("decode global config: %w", err)
	}

	return nil
}
