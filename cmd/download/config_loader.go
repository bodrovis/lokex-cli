package download

import (
	"errors"
	"fmt"

	params "github.com/bodrovis/lokex-cli/internal/params"
	vh "github.com/bodrovis/lokex-cli/internal/viper_helpers"
	"github.com/spf13/viper"
)

func LoadDownloadConfig(cfg *DownloadConfig, configFile, envPrefix string) error {
	if cfg == nil {
		return errors.New("download config is nil")
	}

	v := vh.NewConfigViper(configFile, envPrefix)

	if err := bindDownloadEnv(v); err != nil {
		return fmt.Errorf("bind download env: %w", err)
	}

	if err := vh.ReadOptionalConfig(v, configFile); err != nil {
		return fmt.Errorf("read download config: %w", err)
	}

	params.LoadFromViper(v, cfg, downloadParamSpecs)
	return nil
}

func bindDownloadEnv(v *viper.Viper) error {
	return vh.BindEnvKeys(v, params.ConfigKeys(downloadParamSpecs))
}
