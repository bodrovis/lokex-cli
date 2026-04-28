package viper_helpers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var (
	userConfigDir = os.UserConfigDir
	userHomeDir   = os.UserHomeDir
)

func NewConfigViper(configFile, envPrefix string) *viper.Viper {
	v := viper.New()

	configFile = strings.TrimSpace(configFile)
	envPrefix = strings.TrimSpace(envPrefix)

	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.SetConfigName("lokex")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")

		if configDir, err := userConfigDir(); err == nil && strings.TrimSpace(configDir) != "" {
			v.AddConfigPath(filepath.Join(configDir, "lokex-cli"))
		} else if home, err := userHomeDir(); err == nil && strings.TrimSpace(home) != "" {
			v.AddConfigPath(filepath.Join(home, ".config", "lokex-cli"))
		}
	}

	if envPrefix != "" {
		v.SetEnvPrefix(envPrefix)
	}

	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	v.AutomaticEnv()

	return v
}
