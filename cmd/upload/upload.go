package upload

import (
	"context"
	"errors"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bodrovis/lokex-cli/internal/appstate"
	commandctx "github.com/bodrovis/lokex-cli/internal/commandctx"
	globalCfg "github.com/bodrovis/lokex-cli/internal/global_config"
	"github.com/bodrovis/lokex-cli/internal/params"
	"github.com/bodrovis/lokex-cli/internal/ptrutil"
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

func NewCommand(
	cfg *globalCfg.GlobalConfig,
	state *appstate.State,
) *cobra.Command {
	resolved := &UploadConfig{}

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload translation files to Lokalise",

		PreRunE: func(cmd *cobra.Command, _ []string) error {
			if state == nil || state.Viper == nil {
				return errors.New("application config is not initialized")
			}

			if err := LoadUploadConfig(
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

	params.BindFlags(cmd, uploadParamSpecs)

	return cmd
}

func validateCommand(
	cfg *globalCfg.GlobalConfig,
	uploadCfg *UploadConfig,
) error {
	if cfg == nil {
		return errors.New("global config is nil")
	}

	if uploadCfg == nil {
		return errors.New("upload config is nil")
	}

	if err := cfg.ValidateProjectAccess(); err != nil {
		return err
	}

	if ptrutil.TrimmedString(uploadCfg.Manifest) != "" {
		return nil
	}

	if ptrutil.TrimmedString(uploadCfg.Filename) == "" {
		return errors.New("filename is required")
	}

	if ptrutil.TrimmedString(uploadCfg.LangISO) == "" {
		return errors.New("lang-iso is required")
	}

	return nil
}

func runCommand(
	cmd *cobra.Command,
	cfg *globalCfg.GlobalConfig,
	uploadCfg *UploadConfig,
) error {
	up, err := newUploaderFunc(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := commandctx.NewCommandContext(cfg.ContextTimeout)
	defer cancel()

	if ptrutil.TrimmedString(uploadCfg.Manifest) != "" {
		return runManifestCommand(
			cmd,
			up,
			uploadCfg,
			ctx,
		)
	}

	requestParams, err := buildParamsFunc(uploadCfg)
	if err != nil {
		return err
	}

	poll := ptrutil.Value(uploadCfg.Poll)
	srcPath := ptrutil.Value(uploadCfg.SrcPath)

	result, err := performUpload(
		ctx,
		up,
		requestParams,
		srcPath,
		poll,
	)
	if err != nil {
		return err
	}

	printUploadResult(cmd, result, poll)

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
	params lokexupload.UploadParams,
	srcPath string,
	poll bool,
) (string, error) {
	return up.Upload(ctx, params, srcPath, poll)
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
