package params

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func BindEnv(
	v *viper.Viper,
	prefix string,
	specs []Spec,
) error {
	for _, spec := range specs {
		key := configKey(prefix, spec.Name)

		if err := v.BindEnv(key); err != nil {
			return fmt.Errorf(
				"bind env key %q: %w",
				key,
				err,
			)
		}
	}

	return nil
}

func ApplyChangedFlags(
	v *viper.Viper,
	cmd *cobra.Command,
	prefix string,
	specs []Spec,
) error {
	for _, spec := range specs {
		if !cmd.Flags().Changed(spec.Name) {
			continue
		}

		value, err := flagValue(cmd, spec)
		if err != nil {
			return err
		}

		v.Set(
			configKey(prefix, spec.Name),
			value,
		)
	}

	return nil
}

func ApplyDefaults(
	v *viper.Viper,
	prefix string,
	specs []Spec,
) {
	for _, spec := range specs {
		if spec.Default == nil {
			continue
		}

		v.SetDefault(
			configKey(prefix, spec.Name),
			spec.Default,
		)
	}
}
