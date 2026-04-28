package upload

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
)

type manifestFile struct {
	Items []manifestItem `json:"items"`
}

type manifestItem struct {
	Params  lokexupload.UploadParams `json:"params"`
	SrcPath string                   `json:"src_path"`
}

var (
	loadManifestFileFunc      = loadManifestFile
	buildBatchUploadItemsFunc = buildBatchUploadItems
	performBatchUploadFunc    = performBatchUpload
)

func performBatchUpload(
	ctx context.Context,
	up uploader,
	flags *Flags,
	items []lokexupload.BatchUploadItem,
) (lokexupload.BatchUploadResult, error) {
	return up.UploadBatch(ctx, items, flags.Poll)
}

func runManifestCommand(
	cmd *cobra.Command,
	up uploader,
	flags *Flags,
	ctx context.Context,
) error {
	manifestPath := strings.TrimSpace(flags.Manifest)

	mf, err := loadManifestFileFunc(manifestPath)
	if err != nil {
		return err
	}

	items, err := buildBatchUploadItemsFunc(manifestPath, mf)
	if err != nil {
		return err
	}

	result, err := performBatchUploadFunc(ctx, up, flags, items)
	if err != nil {
		return fmt.Errorf("upload batch from manifest %q: %w", manifestPath, err)
	}

	printBatchUploadResult(cmd, result, flags.Poll)
	return nil
}

func loadManifestFile(path string) (manifestFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return manifestFile{}, fmt.Errorf("read manifest file %q: %w", path, err)
	}

	var mf manifestFile

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()

	if err := dec.Decode(&mf); err != nil {
		return manifestFile{}, fmt.Errorf("parse manifest file %q: %w", path, err)
	}

	if len(mf.Items) == 0 {
		return manifestFile{}, fmt.Errorf("manifest file %q contains no items", path)
	}

	return mf, nil
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

func printBatchUploadResult(cmd *cobra.Command, result lokexupload.BatchUploadResult, poll bool) {
	var successCount int
	var failedCount int

	for _, item := range result.Items {
		if item.Err != nil {
			failedCount++
		} else {
			successCount++
		}

		printBatchItemResult(cmd, item, poll)
	}

	printBatchSummary(cmd, len(result.Items), successCount, failedCount)
}

func printBatchItemResult(cmd *cobra.Command, item lokexupload.BatchUploadResultItem, poll bool) {
	if item.Err != nil {
		cmd.Printf(
			"Upload failed: index=%d src=%q err=%v\n",
			item.Index,
			item.SrcPath,
			item.Err,
		)
		return
	}

	if poll {
		cmd.Printf(
			"Upload completed: index=%d src=%q process_id=%s\n",
			item.Index,
			item.SrcPath,
			item.ProcessID,
		)
		return
	}

	cmd.Printf(
		"Upload started: index=%d src=%q process_id=%s\n",
		item.Index,
		item.SrcPath,
		item.ProcessID,
	)
}

func printBatchSummary(cmd *cobra.Command, total, success, failed int) {
	cmd.Printf(
		"Batch summary: total=%d success=%d failed=%d\n",
		total,
		success,
		failed,
	)
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
