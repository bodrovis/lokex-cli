package download

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	downloadCfg "github.com/bodrovis/lokex-cli/internal/download_config"
	globalCfg "github.com/bodrovis/lokex-cli/internal/global_config"
	lokexdownload "github.com/bodrovis/lokex/v2/client/download"
)

type downloader interface {
	Download(ctx context.Context, out string, params lokexdownload.DownloadParams) (string, error)
	DownloadAsync(ctx context.Context, out string, params lokexdownload.DownloadParams) (string, error)
}

var newDownloaderFunc = newDownloader

func NewCommand(cfg *globalCfg.GlobalConfig, defaults *downloadCfg.DownloadConfig) *cobra.Command {
	flags := newFlags()

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download translation files from Lokalise",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			applyDefaults(cmd, flags, defaults)
			return validateCommand(cfg, flags)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(cmd, cfg, flags, defaults)
		},
	}

	bindFlags(cmd, flags)

	return cmd
}

func validateCommand(cfg *globalCfg.GlobalConfig, flags *Flags) error {
	if err := cfg.ValidateClientConfig(); err != nil {
		return err
	}

	if strings.TrimSpace(flags.Format) == "" {
		return fmt.Errorf("--format is required")
	}

	return nil
}

func runCommand(cmd *cobra.Command, cfg *globalCfg.GlobalConfig, flags *Flags, defaults *downloadCfg.DownloadConfig) error {
	dl, err := newDownloaderFunc(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := newCommandContext(flags.ContextTimeout)
	defer cancel()

	params, err := buildParams(cmd, flags, defaults)
	if err != nil {
		return err
	}

	url, err := performDownload(ctx, dl, flags, params)
	if err != nil {
		return err
	}

	printDownloadResult(cmd, url)

	return nil
}

func newDownloader(cfg *globalCfg.GlobalConfig) (downloader, error) {
	client, err := cfg.NewClient()
	if err != nil {
		return nil, err
	}

	return lokexdownload.NewDownloader(client), nil
}

func newCommandContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		return context.Background(), func() {}
	}

	return context.WithTimeout(context.Background(), timeout)
}

func printDownloadResult(cmd *cobra.Command, url string) {
	cmd.Printf("Bundle downloaded from: %s\n", truncateURLForOutput(url, 150))
}

func performDownload(
	ctx context.Context,
	dl downloader,
	flags *Flags,
	params lokexdownload.DownloadParams,
) (string, error) {
	if flags.Async {
		return dl.DownloadAsync(ctx, flags.Out, params)
	}

	return dl.Download(ctx, flags.Out, params)
}

func truncateURLForOutput(url string, max int) string {
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
