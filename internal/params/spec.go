package params

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ParamSpec[TFlags any, TCfg any, TReq any] struct {
	FlagName  string
	ConfigKey string

	BindFlag      func(cmd *cobra.Command, flags *TFlags)
	ApplyDefault  func(cmd *cobra.Command, flags *TFlags, cfg *TCfg)
	LoadFromViper func(v *viper.Viper, cfg *TCfg)

	ApplyToRequest func(cmd *cobra.Command, flags *TFlags, defaults *TCfg, req TReq) error
}

func BindFlags[TFlags any, TCfg any, TReq any](
	cmd *cobra.Command,
	flags *TFlags,
	specs []ParamSpec[TFlags, TCfg, TReq],
) {
	cmd.Flags().SortFlags = false

	for _, p := range specs {
		if p.BindFlag != nil {
			p.BindFlag(cmd, flags)
		}
	}
}

func ApplyDefaults[TFlags any, TCfg any, TReq any](
	cmd *cobra.Command,
	flags *TFlags,
	cfg *TCfg,
	specs []ParamSpec[TFlags, TCfg, TReq],
) {
	if cfg == nil {
		return
	}

	for _, p := range specs {
		if p.ApplyDefault != nil {
			p.ApplyDefault(cmd, flags, cfg)
		}
	}
}

func LoadFromViper[TFlags any, TCfg any, TReq any](
	v *viper.Viper,
	cfg *TCfg,
	specs []ParamSpec[TFlags, TCfg, TReq],
) {
	if cfg == nil {
		return
	}

	for _, p := range specs {
		if p.LoadFromViper != nil {
			p.LoadFromViper(v, cfg)
		}
	}
}

func ConfigKeys[TFlags any, TCfg any, TReq any](
	specs []ParamSpec[TFlags, TCfg, TReq],
) []string {
	keys := make([]string, 0, len(specs))

	for _, p := range specs {
		key := strings.TrimSpace(p.ConfigKey)
		if key != "" {
			keys = append(keys, key)
		}
	}

	return keys
}

func ApplyToRequest[TFlags any, TCfg any, TReq any](
	cmd *cobra.Command,
	flags *TFlags,
	defaults *TCfg,
	req TReq,
	specs []ParamSpec[TFlags, TCfg, TReq],
) error {
	if defaults == nil {
		defaults = new(TCfg)
	}

	for _, p := range specs {
		if p.ApplyToRequest != nil {
			if err := p.ApplyToRequest(cmd, flags, defaults, req); err != nil {
				return err
			}
		}
	}

	return nil
}
