package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	downloadcmd "github.com/bodrovis/lokex-cli/cmd/download"
	uploadcmd "github.com/bodrovis/lokex-cli/cmd/upload"
	"github.com/bodrovis/lokex-cli/internal/global_config"
)

const skipConfigAnnotation = "skipConfig"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type (
	loadGlobalConfigInputFunc func(string, global_config.LoadOptions) (*global_config.GlobalConfigInput, error)
	loadUploadConfigFunc      func(*uploadcmd.UploadConfig, string, string) error
	loadDownloadConfigFunc    func(*downloadcmd.DownloadConfig, string, string) error
)

func RootCmd() *cobra.Command {
	return newRootCmd(
		global_config.LoadGlobalConfigInput,
		uploadcmd.LoadUploadConfig,
		downloadcmd.LoadDownloadConfig,
	)
}

func newRootCmd(
	loadGlobal loadGlobalConfigInputFunc,
	loadUpload loadUploadConfigFunc,
	loadDownload loadDownloadConfigFunc,
) *cobra.Command {
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
		PersistentPreRunE: newPersistentPreRunE(
			cfg,
			uploadCfg,
			downloadCfg,
			&configFile,
			loadGlobal,
			loadUpload,
			loadDownload,
		),
	}

	global_config.BindPersistentFlags(cmd.PersistentFlags(), cfg)
	cmd.PersistentFlags().StringVar(&configFile, "config", "", "Path to YAML config file")

	versionCmd := newVersionCmd()
	markSkipConfig(versionCmd)
	cmd.AddCommand(versionCmd)

	genDocsCmd := newGenDocsCmd(cmd)
	markSkipConfig(genDocsCmd)
	cmd.AddCommand(genDocsCmd)

	cmd.AddCommand(downloadcmd.NewCommand(cfg, downloadCfg))
	cmd.AddCommand(uploadcmd.NewCommand(cfg, uploadCfg))

	return cmd
}

func newPersistentPreRunE(
	cfg *global_config.GlobalConfig,
	uploadCfg *uploadcmd.UploadConfig,
	downloadCfg *downloadcmd.DownloadConfig,
	configFile *string,
	loadGlobal loadGlobalConfigInputFunc,
	loadUpload loadUploadConfigFunc,
	loadDownload loadDownloadConfigFunc,
) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if shouldSkipConfig(cmd) {
			return nil
		}

		loadOpts := global_config.LoadOptions{
			ConfigFile: *configFile,
			EnvPrefix:  "LOKEX",
		}

		globalInput, err := loadGlobal(cfg.UserAgent, loadOpts)
		if err != nil {
			return err
		}

		global_config.ApplyGlobalInput(cmd, cfg, globalInput)

		switch cmd.Name() {
		case "upload":
			if err := loadUpload(uploadCfg, loadOpts.ConfigFile, loadOpts.EnvPrefix); err != nil {
				return err
			}
		case "download":
			if err := loadDownload(downloadCfg, loadOpts.ConfigFile, loadOpts.EnvPrefix); err != nil {
				return err
			}
		}

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
