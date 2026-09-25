package manifest

import (
	"errors"
	"fmt"

	"github.com/bodrovis/lokex-cli/internal/params"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LoadGenerateConfig(
	v *viper.Viper,
	cmd *cobra.Command,
	cfg *GenerateConfig,
) error {
	if v == nil {
		return errors.New("viper is nil")
	}

	if cmd == nil {
		return errors.New("manifest generate command is nil")
	}

	if cfg == nil {
		return errors.New("manifest generate config is nil")
	}

	params.ApplyDefaults(
		v,
		generateConfigPrefix,
		generateParamSpecs,
	)

	if err := params.BindEnv(
		v,
		generateConfigPrefix,
		generateParamSpecs,
	); err != nil {
		return fmt.Errorf(
			"bind manifest generate environment: %w",
			err,
		)
	}

	if err := params.ApplyChangedFlags(
		v,
		cmd,
		generateConfigPrefix,
		generateParamSpecs,
	); err != nil {
		return fmt.Errorf(
			"apply manifest generate flags: %w",
			err,
		)
	}

	var root struct {
		Manifest struct {
			Generate GenerateConfig `mapstructure:"generate"`
		} `mapstructure:"manifest"`
	}

	if err := v.Unmarshal(&root); err != nil {
		return fmt.Errorf(
			"decode manifest generate config: %w",
			err,
		)
	}

	*cfg = root.Manifest.Generate

	return nil
}
