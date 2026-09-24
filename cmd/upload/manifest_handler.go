package upload

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bodrovis/lokex-cli/internal/ptrutil"
	"github.com/bodrovis/lokex-cli/internal/uploadmanifest"
	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
)

var (
	loadManifestFileFunc      = uploadmanifest.Load
	buildBatchUploadItemsFunc = uploadmanifest.BuildBatchUploadItems
	performBatchUploadFunc    = performBatchUpload
)

func performBatchUpload(
	ctx context.Context,
	up uploader,
	poll bool,
	items []lokexupload.BatchUploadItem,
) (lokexupload.BatchUploadResult, error) {
	return up.UploadBatch(ctx, items, poll)
}

func runManifestCommand(
	cmd *cobra.Command,
	up uploader,
	cfg *UploadConfig,
	ctx context.Context,
) error {
	manifestPath := ptrutil.TrimmedString(cfg.Manifest)
	poll := ptrutil.Value(cfg.Poll)

	mf, err := loadManifestFileFunc(manifestPath)
	if err != nil {
		return err
	}

	items, err := buildBatchUploadItemsFunc(manifestPath, mf)
	if err != nil {
		return err
	}

	result, err := performBatchUploadFunc(
		ctx,
		up,
		poll,
		items,
	)
	if err != nil {
		return fmt.Errorf(
			"upload batch from manifest %q: %w",
			manifestPath,
			err,
		)
	}

	printBatchUploadResult(cmd, result, poll)

	return nil
}
