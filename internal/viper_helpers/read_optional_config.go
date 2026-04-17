package viper_helpers

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

func ReadOptionalConfig(v *viper.Viper, explicitPath string) error {
	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if strings.TrimSpace(explicitPath) != "" {
			return err
		}
		if !errors.As(err, &notFound) {
			return err
		}
	}
	return nil
}
