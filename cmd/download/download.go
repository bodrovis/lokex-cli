package download

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/bodrovis/lokex-cli/internal/cli"
	lokexdownload "github.com/bodrovis/lokex/v2/client/download"
)

func NewCommand(cfg *cli.GlobalConfig) *cobra.Command {
	flags := newFlags()

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download translation files from Lokalise",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateCommand(cfg, flags)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(cmd, cfg, flags)
		},
	}

	bindFlags(cmd, flags)

	return cmd
}

func validateCommand(cfg *cli.GlobalConfig, flags *Flags) error {
	if err := cfg.ValidateAPIConfig(); err != nil {
		return err
	}

	if flags.Format == "" {
		return fmt.Errorf("--format is required")
	}

	return nil
}

func runCommand(cmd *cobra.Command, cfg *cli.GlobalConfig, flags *Flags) error {
	downloader, err := newDownloader(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := newCommandContext(flags.ContextTimeout)
	defer cancel()

	params, err := buildParams(cmd, flags)
	if err != nil {
		return err
	}

	url, err := performDownload(ctx, downloader, flags, params)
	if err != nil {
		return err
	}

	printDownloadResult(cmd, url)

	return nil
}

func newDownloader(cfg *cli.GlobalConfig) (*lokexdownload.Downloader, error) {
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
	downloader *lokexdownload.Downloader,
	flags *Flags,
	params lokexdownload.DownloadParams,
) (string, error) {
	if flags.Async {
		return downloader.DownloadAsync(ctx, flags.Out, params)
	}

	return downloader.Download(ctx, flags.Out, params)
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
