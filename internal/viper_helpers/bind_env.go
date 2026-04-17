package viper_helpers

import (
	"github.com/spf13/viper"
)

func BindEnvKeys(v *viper.Viper, keys []string) error {
	for _, key := range keys {
		if err := v.BindEnv(key); err != nil {
			return err
		}
	}
	return nil
}
