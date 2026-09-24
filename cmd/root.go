package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	downloadcmd "github.com/bodrovis/lokex-cli/cmd/download"
	uploadcmd "github.com/bodrovis/lokex-cli/cmd/upload"
	"github.com/bodrovis/lokex-cli/internal/appstate"
	"github.com/bodrovis/lokex-cli/internal/global_config"
	"github.com/bodrovis/lokex-cli/internal/viper_helpers"
)

const skipConfigAnnotation = "skipConfig"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type loadGlobalConfigFunc func(
	*viper.Viper,
	*cobra.Command,
	*global_config.GlobalConfig,
) error

func RootCmd() *cobra.Command {
	return newRootCmd(global_config.LoadGlobalConfig)
}

func newRootCmd(
	loadGlobal loadGlobalConfigFunc,
) *cobra.Command {
	cfg := &global_config.GlobalConfig{
		UserAgent: fmt.Sprintf("lokex-cli/%s", version),
	}

	state := &appstate.State{}

	var configFile string

	cmd := &cobra.Command{
		Use:   "lokex-cli",
		Short: "CLI for uploading and downloading files with Lokalise",
		Long: `lokex-cli is a focused CLI built on top of Lokex for fast file exchange with Lokalise.

It is intentionally limited to two core operations:

  - upload files
  - download files

This tool is optimized for import/export workflows and direct access to file-related API parameters.
`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: newPersistentPreRunE(
			cfg,
			state,
			&configFile,
			loadGlobal,
		),
	}

	global_config.BindPersistentFlags(
		cmd.PersistentFlags(),
		cfg.UserAgent,
	)

	cmd.PersistentFlags().StringVar(
		&configFile,
		"config",
		"",
		"Path to YAML config file",
	)

	versionCmd := newVersionCmd()
	markSkipConfig(versionCmd)
	cmd.AddCommand(versionCmd)

	genDocsCmd := newGenDocsCmd(cmd)
	markSkipConfig(genDocsCmd)
	cmd.AddCommand(genDocsCmd)

	cmd.AddCommand(downloadcmd.NewCommand(cfg, state))
	cmd.AddCommand(uploadcmd.NewCommand(cfg, state))

	return cmd
}

func newPersistentPreRunE(
	cfg *global_config.GlobalConfig,
	state *appstate.State,
	configFile *string,
	loadGlobal loadGlobalConfigFunc,
) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		if shouldSkipConfig(cmd) {
			return nil
		}

		v := viper_helpers.NewConfigViper(
			*configFile,
			"LOKEX",
		)

		if err := viper_helpers.ReadOptionalConfig(
			v,
			*configFile,
		); err != nil {
			return fmt.Errorf("read config: %w", err)
		}

		if err := loadGlobal(
			v,
			cmd,
			cfg,
		); err != nil {
			return err
		}

		state.Viper = v

		return nil
	}
}

func markSkipConfig(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}

	cmd.Annotations[skipConfigAnnotation] = "true"
}

func shouldSkipConfig(cmd *cobra.Command) bool {
	for c := cmd; c != nil; c = c.Parent() {
		if c.Annotations[skipConfigAnnotation] == "true" {
			return true
		}
	}

	return false
}
