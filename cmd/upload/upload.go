package upload

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	globalCfg "github.com/bodrovis/lokex-cli/internal/global_config"
	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
)

type uploader interface {
	Upload(ctx context.Context, params lokexupload.UploadParams, srcPath string, poll bool) (string, error)
}

var newUploaderFunc = newUploader

func NewCommand(cfg *globalCfg.GlobalConfig, defaults *UploadConfig) *cobra.Command {
	flags := newFlags()

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload translation files to Lokalise",
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

	if strings.TrimSpace(flags.Filename) == "" {
		return fmt.Errorf("--filename is required")
	}

	if strings.TrimSpace(flags.LangISO) == "" {
		return fmt.Errorf("--lang-iso is required")
	}

	return nil
}

func runCommand(cmd *cobra.Command, cfg *globalCfg.GlobalConfig, flags *Flags, defaults *UploadConfig) error {
	up, err := newUploaderFunc(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := newCommandContext(cfg.ContextTimeout)
	defer cancel()

	params, err := buildParams(cmd, flags, defaults)
	if err != nil {
		return err
	}

	result, err := performUpload(ctx, up, flags, params)
	if err != nil {
		return err
	}

	printUploadResult(cmd, result, flags.Poll)

	return nil
}

func newUploader(cfg *globalCfg.GlobalConfig) (uploader, error) {
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
	up uploader,
	flags *Flags,
	params lokexupload.UploadParams,
) (string, error) {
	return up.Upload(ctx, params, flags.SrcPath, flags.Poll)
}

func printUploadResult(cmd *cobra.Command, result string, poll bool) {
	result = strings.TrimSpace(result)

	if result == "" {
		if poll {
			cmd.Println("Upload completed (process ID unknown)")
			return
		}
		cmd.Println("Upload started (process ID unknown)")
		return
	}

	if poll {
		cmd.Printf("Upload completed: %s\n", result)
		return
	}

	cmd.Printf("Upload started: %s\n", result)
}
