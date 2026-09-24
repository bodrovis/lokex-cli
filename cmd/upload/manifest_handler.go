package upload

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bodrovis/lokex-cli/internal/ptrutil"
	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
)

var (
	loadManifestFileFunc      = loadManifestFile
	buildBatchUploadItemsFunc = buildBatchUploadItems
	performBatchUploadFunc    = performBatchUpload
)

var preserveJSONNumbers = json.WithUnmarshalers(
	json.UnmarshalFromFunc(func(dec *jsontext.Decoder, val *any) error {
		if dec.PeekKind() == jsontext.KindNumber {
			*val = jsontext.Value(nil)
		}

		return errors.ErrUnsupported
	}),
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

func buildBatchUploadItems(manifestPath string, mf manifestFile) ([]lokexupload.BatchUploadItem, error) {
	manifestDir := filepath.Dir(manifestPath)

	items := make([]lokexupload.BatchUploadItem, 0, len(mf.Items))
	for i, item := range mf.Items {
		if err := validateManifestItem(item, i+1); err != nil {
			return nil, err
		}

		srcPath := item.SrcPath
		if srcPath != "" && !filepath.IsAbs(srcPath) {
			srcPath = filepath.Join(manifestDir, srcPath)
		}

		items = append(items, lokexupload.BatchUploadItem{
			Params:  item.Params,
			SrcPath: srcPath,
		})
	}

	return items, nil
}

func validateManifestItem(item manifestItem, index int) error {
	if item.Params == nil {
		return fmt.Errorf("manifest item %d: params is required", index)
	}

	if _, ok := requiredManifestString(item.Params, "filename"); !ok {
		return fmt.Errorf("manifest item %d: params.filename is required", index)
	}

	if _, ok := requiredManifestString(item.Params, "lang_iso"); !ok {
		return fmt.Errorf("manifest item %d: params.lang_iso is required", index)
	}

	return nil
}

func requiredManifestString(params lokexupload.UploadParams, key string) (string, bool) {
	val, ok := params[key]
	if !ok {
		return "", false
	}

	s, ok := val.(string)
	if !ok {
		return "", false
	}

	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}

	return s, true
}
