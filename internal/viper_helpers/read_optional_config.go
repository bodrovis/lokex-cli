package viper_helpers

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

func ReadOptionalConfig(v *viper.Viper, explicitPath string) error {
	if err := v.ReadInConfig(); err != nil {
		if strings.TrimSpace(explicitPath) != "" {
			return err
		}
		if _, ok := errors.AsType[viper.ConfigFileNotFoundError](err); !ok {
			return err
		}
	}
	return nil
}
