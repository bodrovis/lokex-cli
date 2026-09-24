package download

import (
	"errors"
	"fmt"

	"github.com/bodrovis/lokex-cli/internal/params"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LoadDownloadConfig(
	v *viper.Viper,
	cmd *cobra.Command,
	cfg *DownloadConfig,
) error {
	if v == nil {
		return errors.New("viper is nil")
	}

	if cmd == nil {
		return errors.New("download command is nil")
	}

	if cfg == nil {
		return errors.New("download config is nil")
	}

	params.ApplyDefaults(
		v,
		configPrefix,
		downloadParamSpecs,
	)

	if err := params.BindEnv(
		v,
		configPrefix,
		downloadParamSpecs,
	); err != nil {
		return fmt.Errorf("bind download environment: %w", err)
	}

	if err := params.ApplyChangedFlags(
		v,
		cmd,
		configPrefix,
		downloadParamSpecs,
	); err != nil {
		return fmt.Errorf("apply download flags: %w", err)
	}

	var root struct {
		Download DownloadConfig `mapstructure:"download"`
	}

	if err := v.Unmarshal(&root); err != nil {
		return fmt.Errorf("decode download config: %w", err)
	}

	*cfg = root.Download

	return nil
}
