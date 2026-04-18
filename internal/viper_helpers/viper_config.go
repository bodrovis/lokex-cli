package viper_helpers

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

var userHomeDir = os.UserHomeDir

func NewConfigViper(configFile, envPrefix string) *viper.Viper {
	v := viper.New()

	if strings.TrimSpace(configFile) != "" {
		v.SetConfigFile(configFile)
	} else {
		v.SetConfigName("lokex")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		if home, err := userHomeDir(); err == nil {
			v.AddConfigPath(home + "/.config/lokex-cli")
		}
	}

	if strings.TrimSpace(envPrefix) != "" {
		v.SetEnvPrefix(envPrefix)
	}

	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	v.AutomaticEnv()

	return v
}
