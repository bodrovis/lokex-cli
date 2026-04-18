package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	downloadcmd "github.com/bodrovis/lokex-cli/cmd/download"
	uploadcmd "github.com/bodrovis/lokex-cli/cmd/upload"
	"github.com/bodrovis/lokex-cli/internal/global_config"
)

var version = "dev"

func RootCmd() *cobra.Command {
	cfg := &global_config.GlobalConfig{
		UserAgent: fmt.Sprintf("lokex-cli/%s", version),
	}
	uploadCfg := &uploadcmd.UploadConfig{}
	downloadCfg := &downloadcmd.DownloadConfig{}

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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			loadOpts := global_config.LoadOptions{
				ConfigFile: configFile,
				EnvPrefix:  "LOKEX",
			}

			globalInput, err := global_config.LoadGlobalConfigInput(cfg.UserAgent, loadOpts)
			if err != nil {
				return err
			}

			global_config.ApplyGlobalDefaults(cmd, cfg, globalInput)

			switch cmd.Name() {
			case "upload":
				if err := uploadcmd.LoadUploadConfig(uploadCfg, loadOpts.ConfigFile, loadOpts.EnvPrefix); err != nil {
					return err
				}
			case "download":
				if err := downloadcmd.LoadDownloadConfig(downloadCfg, loadOpts.ConfigFile, loadOpts.EnvPrefix); err != nil {
					return err
				}
			}

			return nil
		},
	}

	global_config.BindPersistentFlags(cmd.PersistentFlags(), cfg)
	cmd.PersistentFlags().StringVar(&configFile, "config", "", "Path to YAML config file")

	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newGenDocsCmd(cmd))
	cmd.AddCommand(downloadcmd.NewCommand(cfg, downloadCfg))
	cmd.AddCommand(uploadcmd.NewCommand(cfg, uploadCfg))

	return cmd
}
