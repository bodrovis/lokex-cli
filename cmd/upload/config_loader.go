package upload

import (
	"errors"
	"fmt"

	"github.com/bodrovis/lokex-cli/internal/params"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LoadUploadConfig(
	v *viper.Viper,
	cmd *cobra.Command,
	cfg *UploadConfig,
) error {
	if v == nil {
		return errors.New("viper is nil")
	}

	if cmd == nil {
		return errors.New("upload command is nil")
	}

	if cfg == nil {
		return errors.New("upload config is nil")
	}

	params.ApplyDefaults(
		v,
		configPrefix,
		uploadParamSpecs,
	)

	if err := params.BindEnv(
		v,
		configPrefix,
		uploadParamSpecs,
	); err != nil {
		return fmt.Errorf("bind upload environment: %w", err)
	}

	if err := params.ApplyChangedFlags(
		v,
		cmd,
		configPrefix,
		uploadParamSpecs,
	); err != nil {
		return fmt.Errorf("apply upload flags: %w", err)
	}

	var root struct {
		Upload UploadConfig `mapstructure:"upload"`
	}

	if err := v.Unmarshal(&root); err != nil {
		return fmt.Errorf("decode upload config: %w", err)
	}

	*cfg = root.Upload

	return nil
}
