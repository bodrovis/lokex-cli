package upload

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	commandctx "github.com/bodrovis/lokex-cli/internal/commandctx"
	globalCfg "github.com/bodrovis/lokex-cli/internal/global_config"
	params "github.com/bodrovis/lokex-cli/internal/params"
	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
)

type uploader interface {
	Upload(ctx context.Context, params lokexupload.UploadParams, srcPath string, poll bool) (string, error)
	UploadBatch(ctx context.Context, items []lokexupload.BatchUploadItem, poll bool) (lokexupload.BatchUploadResult, error)
}

var (
	newUploaderFunc = newUploader
	buildParamsFunc = buildParams
)

func NewCommand(cfg *globalCfg.GlobalConfig, defaults *UploadConfig) *cobra.Command {
	flags := newFlags()

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload translation files to Lokalise",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			params.ApplyDefaults(cmd, flags, defaults, uploadParamSpecs)
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
	if cfg == nil {
		return fmt.Errorf("global config is nil")
	}

	if flags == nil {
		return fmt.Errorf("upload flags are nil")
	}

	if err := cfg.ValidateClientConfig(); err != nil {
		return err
	}

	if strings.TrimSpace(flags.Manifest) != "" {
		return nil
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

	ctx, cancel := commandctx.NewCommandContext(cfg.ContextTimeout)
	defer cancel()

	if strings.TrimSpace(flags.Manifest) != "" {
		return runManifestCommand(cmd, up, flags, ctx)
	}

	params, err := buildParamsFunc(cmd, flags, defaults)
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
