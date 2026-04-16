package cli

import (
	"strings"

	"github.com/spf13/cobra"
)

func SetString[P ~map[string]any](params P, key, value string) {
	if strings.TrimSpace(value) != "" {
		params[key] = value
	}
}

func SetStringSlice[P ~map[string]any](params P, key string, value []string) {
	if len(value) > 0 {
		params[key] = value
	}
}

func SetChangedBool[P ~map[string]any](cmd *cobra.Command, params P, flagName, paramName string, value bool) {
	if cmd.Flags().Changed(flagName) {
		params[paramName] = value
	}
}
