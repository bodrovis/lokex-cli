package params

import (
	"strings"

	"github.com/spf13/cobra"
)

func SetString[P ~map[string]any](params P, key, value string) {
	v := strings.TrimSpace(value)
	if v != "" {
		params[key] = v
	}
}

func SetStringSlice[P ~map[string]any](params P, key string, value []string) {
	clean := make([]string, 0, len(value))
	for _, v := range value {
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			clean = append(clean, trimmed)
		}
	}
	if len(clean) > 0 {
		params[key] = clean
	}
}

func SetBoolWithDefault[P ~map[string]any](
	cmd *cobra.Command,
	params P,
	flagName, paramName string,
	flagValue bool,
	cfgValue *bool,
) {
	if cmd.Flags().Changed(flagName) {
		params[paramName] = flagValue
		return
	}

	if cfgValue != nil {
		params[paramName] = *cfgValue
	}
}

func SetInt64WithDefault[P ~map[string]any](
	cmd *cobra.Command,
	params P,
	flagName, paramName string,
	flagValue int64,
	cfgValue *int64,
) {
	if cmd.Flags().Changed(flagName) {
		params[paramName] = flagValue
		return
	}

	if cfgValue != nil {
		params[paramName] = *cfgValue
	}
}
