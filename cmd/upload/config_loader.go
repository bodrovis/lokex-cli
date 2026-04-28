package upload

import (
	"errors"
	"fmt"

	params "github.com/bodrovis/lokex-cli/internal/params"
	vh "github.com/bodrovis/lokex-cli/internal/viper_helpers"
	"github.com/spf13/viper"
)

func LoadUploadConfig(cfg *UploadConfig, configFile, envPrefix string) error {
	if cfg == nil {
		return errors.New("upload config is nil")
	}

	v := vh.NewConfigViper(configFile, envPrefix)

	if err := bindUploadEnv(v); err != nil {
		return fmt.Errorf("bind upload env: %w", err)
	}

	if err := vh.ReadOptionalConfig(v, configFile); err != nil {
		return fmt.Errorf("read upload config: %w", err)
	}

	params.LoadFromViper(v, cfg, uploadParamSpecs)
	return nil
}

func bindUploadEnv(v *viper.Viper) error {
	return vh.BindEnvKeys(v, params.ConfigKeys(uploadParamSpecs))
}
