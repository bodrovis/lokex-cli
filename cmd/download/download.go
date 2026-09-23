package download

import (
	"context"
	"errors"

	"github.com/spf13/cobra"

	"github.com/bodrovis/lokex-cli/internal/appstate"
	commandctx "github.com/bodrovis/lokex-cli/internal/commandctx"
	globalCfg "github.com/bodrovis/lokex-cli/internal/global_config"
	"github.com/bodrovis/lokex-cli/internal/params"
	"github.com/bodrovis/lokex-cli/internal/ptrutil"
	lokexdownload "github.com/bodrovis/lokex/v2/client/download"
)

type downloader interface {
	Download(
		ctx context.Context,
		out string,
		params lokexdownload.DownloadParams,
	) (string, error)

	DownloadAsync(
		ctx context.Context,
		out string,
		params lokexdownload.DownloadParams,
	) (string, error)
}

var newDownloaderFunc = newDownloader

func NewCommand(
	cfg *globalCfg.GlobalConfig,
	state *appstate.State,
) *cobra.Command {
	resolved := &DownloadConfig{}

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download translation files from Lokalise",

		PreRunE: func(cmd *cobra.Command, _ []string) error {
			if state == nil || state.Viper == nil {
				return errors.New("application config is not initialized")
			}

			if err := LoadDownloadConfig(
				state.Viper,
				cmd,
				resolved,
			); err != nil {
				return err
			}

			return validateCommand(cfg, resolved)
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCommand(cmd, cfg, resolved)
		},
	}

	params.BindFlags(cmd, downloadParamSpecs)

	return cmd
}

func validateCommand(
	cfg *globalCfg.GlobalConfig,
	downloadCfg *DownloadConfig,
) error {
	if cfg == nil {
		return errors.New("global config is nil")
	}

	if downloadCfg == nil {
		return errors.New("download config is nil")
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	if ptrutil.TrimmedString(downloadCfg.Format) == "" {
		return errors.New("--format is required")
	}

	return nil
}

func runCommand(
	cmd *cobra.Command,
	cfg *globalCfg.GlobalConfig,
	downloadCfg *DownloadConfig,
) error {
	dl, err := newDownloaderFunc(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := commandctx.NewCommandContext(cfg.ContextTimeout)
	defer cancel()

	requestParams, err := buildParams(downloadCfg)
	if err != nil {
		return err
	}

	url, err := performDownload(
		ctx,
		dl,
		requestParams,
		ptrutil.Value(downloadCfg.Out),
		ptrutil.Value(downloadCfg.Async),
	)
	if err != nil {
		return err
	}

	printDownloadResult(cmd, url)

	return nil
}

func newDownloader(
	cfg *globalCfg.GlobalConfig,
) (downloader, error) {
	client, err := cfg.NewClient()
	if err != nil {
		return nil, err
	}

	return lokexdownload.NewDownloader(client), nil
}

func performDownload(
	ctx context.Context,
	dl downloader,
	requestParams lokexdownload.DownloadParams,
	out string,
	async bool,
) (string, error) {
	if async {
		return dl.DownloadAsync(ctx, out, requestParams)
	}

	return dl.Download(ctx, out, requestParams)
}

func printDownloadResult(
	cmd *cobra.Command,
	url string,
) {
	cmd.Printf(
		"Bundle downloaded from: %s\n",
		truncateURLForOutput(url, 150),
	)
}

func truncateURLForOutput(
	url string,
	max int,
) string {
	if max <= 0 {
		return ""
	}

	if len(url) <= max {
		return url
	}

	if max <= 3 {
		return url[:max]
	}

	return url[:max-3] + "..."
}
