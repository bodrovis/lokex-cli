package upload

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bodrovis/lokex-cli/internal/cli"
	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
)

func NewCommand(cfg *cli.GlobalConfig) *cobra.Command {
	flags := &Flags{}

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload translation files to Lokalise",
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

	if strings.TrimSpace(flags.Filename) == "" {
		return fmt.Errorf("--filename is required")
	}

	if strings.TrimSpace(flags.LangISO) == "" {
		return fmt.Errorf("--lang-iso is required")
	}

	return nil
}

func runCommand(cmd *cobra.Command, cfg *cli.GlobalConfig, flags *Flags) error {
	uploader, err := newUploader(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := newCommandContext(flags.ContextTimeout)
	defer cancel()

	params := buildParams(cmd, flags)

	result, err := performUpload(ctx, uploader, flags, params)
	if err != nil {
		return err
	}

	printUploadResult(cmd, result, flags.Poll)

	return nil
}

func newUploader(cfg *cli.GlobalConfig) (*lokexupload.Uploader, error) {
	client, err := cfg.NewClient()
	if err != nil {
		return nil, err
	}

	return lokexupload.NewUploader(client), nil
}

func newCommandContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		return context.Background(), func() {}
	}

	return context.WithTimeout(context.Background(), timeout)
}

func performUpload(
	ctx context.Context,
	uploader *lokexupload.Uploader,
	flags *Flags,
	params lokexupload.UploadParams,
) (string, error) {
	return uploader.Upload(ctx, params, flags.SrcPath, flags.Poll)
}

func printUploadResult(cmd *cobra.Command, result string, poll bool) {
	if poll {
		cmd.Printf("Upload completed: %s\n", result)
		return
	}

	cmd.Printf("Upload started: %s\n", result)
}
