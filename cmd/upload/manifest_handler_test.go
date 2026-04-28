package upload

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/bodrovis/lokex-cli/internal/global_config"
	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
)

func TestPerformBatchUpload(t *testing.T) {
	items := []lokexupload.BatchUploadItem{
		{
			Params: lokexupload.UploadParams{
				"filename": "locales/en.json",
				"lang_iso": "en",
			},
			SrcPath: "./locales/en.json",
		},
		{
			Params: lokexupload.UploadParams{
				"filename": "locales/de.json",
				"lang_iso": "de",
			},
			SrcPath: "./locales/de.json",
		},
	}

	t.Run("success", func(t *testing.T) {
		mu := &mockUploader{
			batchResult: lokexupload.BatchUploadResult{
				Items: []lokexupload.BatchUploadResultItem{
					{
						Index:     0,
						SrcPath:   "./locales/en.json",
						ProcessID: "process-1",
					},
					{
						Index:     1,
						SrcPath:   "./locales/de.json",
						ProcessID: "process-2",
					},
				},
			},
		}
		flags := &Flags{
			Poll: true,
		}

		got, err := performBatchUpload(context.Background(), mu, flags, items)
		if err != nil {
			t.Fatalf("performBatchUpload() error = %v", err)
		}
		if !mu.batchCalled {
			t.Fatal("expected UploadBatch to be called")
		}
		if mu.gotBatchCtx == nil {
			t.Fatal("expected context to be passed")
		}
		if !mu.gotBatchPoll {
			t.Fatal("expected poll=true to be passed")
		}
		if len(mu.gotBatchItems) != 2 {
			t.Fatalf("unexpected batch items len: got %d, want 2", len(mu.gotBatchItems))
		}
		if mu.gotBatchItems[0].SrcPath != "./locales/en.json" {
			t.Fatalf("unexpected first src path: got %q", mu.gotBatchItems[0].SrcPath)
		}
		if mu.gotBatchItems[1].Params["lang_iso"] != "de" {
			t.Fatalf("unexpected second params: %#v", mu.gotBatchItems[1].Params)
		}
		if len(got.Items) != 2 {
			t.Fatalf("unexpected result items len: got %d, want 2", len(got.Items))
		}
		if got.Items[0].ProcessID != "process-1" {
			t.Fatalf("unexpected first process id: got %q", got.Items[0].ProcessID)
		}
	})

	t.Run("error", func(t *testing.T) {
		mu := &mockUploader{
			batchErr: errors.New("batch upload failed"),
		}
		flags := &Flags{
			Poll: false,
		}

		_, err := performBatchUpload(context.Background(), mu, flags, items)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "batch upload failed" {
			t.Fatalf("unexpected error: %q", err.Error())
		}
		if !mu.batchCalled {
			t.Fatal("expected UploadBatch to be called")
		}
		if mu.gotBatchPoll {
			t.Fatal("expected poll=false to be passed")
		}
	})
}

func TestRunCommand_WithManifest(t *testing.T) {
	t.Run("happy path without poll", func(t *testing.T) {
		oldUploader := newUploaderFunc
		oldLoad := loadManifestFileFunc
		oldBuild := buildBatchUploadItemsFunc
		oldPerform := performBatchUploadFunc
		t.Cleanup(func() {
			newUploaderFunc = oldUploader
			loadManifestFileFunc = oldLoad
			buildBatchUploadItemsFunc = oldBuild
			performBatchUploadFunc = oldPerform
		})

		mu := &mockUploader{}
		newUploaderFunc = func(cfg *global_config.GlobalConfig) (uploader, error) {
			return mu, nil
		}

		loadManifestFileFunc = func(path string) (manifestFile, error) {
			if path != "./manifest.json" {
				t.Fatalf("unexpected manifest path: %q", path)
			}
			return manifestFile{
				Items: []manifestItem{
					{
						Params: lokexupload.UploadParams{
							"filename": "locales/en.json",
							"lang_iso": "en",
						},
						SrcPath: "./locales/en.json",
					},
				},
			}, nil
		}

		buildBatchUploadItemsFunc = func(manifestPath string, mf manifestFile) ([]lokexupload.BatchUploadItem, error) {
			if manifestPath != "./manifest.json" {
				t.Fatalf("unexpected manifest path: %q", manifestPath)
			}
			if len(mf.Items) != 1 {
				t.Fatalf("unexpected manifest items len: %d", len(mf.Items))
			}
			return []lokexupload.BatchUploadItem{
				{
					Params: lokexupload.UploadParams{
						"filename": "locales/en.json",
						"lang_iso": "en",
					},
					SrcPath: "./locales/en.json",
				},
			}, nil
		}

		performBatchUploadFunc = func(
			ctx context.Context,
			up uploader,
			flags *Flags,
			items []lokexupload.BatchUploadItem,
		) (lokexupload.BatchUploadResult, error) {
			if up != mu {
				t.Fatal("expected mock uploader to be passed")
			}
			if flags.Manifest != "./manifest.json" {
				t.Fatalf("unexpected manifest flag: %q", flags.Manifest)
			}
			if flags.Poll {
				t.Fatal("expected poll=false")
			}
			if len(items) != 1 {
				t.Fatalf("unexpected items len: %d", len(items))
			}

			return lokexupload.BatchUploadResult{
				Items: []lokexupload.BatchUploadResultItem{
					{
						Index:     0,
						SrcPath:   "./locales/en.json",
						ProcessID: "batch-123",
					},
				},
			}, nil
		}

		cfg := &global_config.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}
		flags := newFlags()

		cmd := newBoundTestCommand(flags)
		if err := cmd.Flags().Parse([]string{
			"--manifest=./manifest.json",
		}); err != nil {
			t.Fatalf("parse flags: %v", err)
		}

		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		err := runCommand(cmd, cfg, flags, nil)
		if err != nil {
			t.Fatalf("runCommand() error = %v", err)
		}

		if mu.uploadCalled {
			t.Fatal("expected Upload not to be called in manifest mode")
		}

		gotOutput := out.String()
		if !strings.Contains(gotOutput, `Upload started: index=0 src="./locales/en.json" process_id=batch-123`) {
			t.Fatalf("unexpected output: %q", gotOutput)
		}
		if !strings.Contains(gotOutput, "Batch summary: total=1 success=1 failed=0") {
			t.Fatalf("unexpected output: %q", gotOutput)
		}
	})

	t.Run("happy path with poll", func(t *testing.T) {
		oldUploader := newUploaderFunc
		oldLoad := loadManifestFileFunc
		oldBuild := buildBatchUploadItemsFunc
		oldPerform := performBatchUploadFunc
		t.Cleanup(func() {
			newUploaderFunc = oldUploader
			loadManifestFileFunc = oldLoad
			buildBatchUploadItemsFunc = oldBuild
			performBatchUploadFunc = oldPerform
		})

		mu := &mockUploader{}
		newUploaderFunc = func(cfg *global_config.GlobalConfig) (uploader, error) {
			return mu, nil
		}

		loadManifestFileFunc = func(path string) (manifestFile, error) {
			return manifestFile{
				Items: []manifestItem{
					{
						Params: lokexupload.UploadParams{
							"filename": "locales/de.json",
							"lang_iso": "de",
						},
						SrcPath: "./locales/de.json",
					},
				},
			}, nil
		}

		buildBatchUploadItemsFunc = func(manifestPath string, mf manifestFile) ([]lokexupload.BatchUploadItem, error) {
			return []lokexupload.BatchUploadItem{
				{
					Params: lokexupload.UploadParams{
						"filename": "locales/de.json",
						"lang_iso": "de",
					},
					SrcPath: "./locales/de.json",
				},
			}, nil
		}

		performBatchUploadFunc = func(
			ctx context.Context,
			up uploader,
			flags *Flags,
			items []lokexupload.BatchUploadItem,
		) (lokexupload.BatchUploadResult, error) {
			if !flags.Poll {
				t.Fatal("expected poll=true")
			}

			return lokexupload.BatchUploadResult{
				Items: []lokexupload.BatchUploadResultItem{
					{
						Index:     0,
						SrcPath:   "./locales/de.json",
						ProcessID: "batch-456",
					},
				},
			}, nil
		}

		cfg := &global_config.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}
		flags := newFlags()

		cmd := newBoundTestCommand(flags)
		if err := cmd.Flags().Parse([]string{
			"--manifest=./manifest.json",
			"--poll",
		}); err != nil {
			t.Fatalf("parse flags: %v", err)
		}

		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		err := runCommand(cmd, cfg, flags, nil)
		if err != nil {
			t.Fatalf("runCommand() error = %v", err)
		}

		if mu.uploadCalled {
			t.Fatal("expected Upload not to be called in manifest mode")
		}

		gotOutput := out.String()
		if !strings.Contains(gotOutput, `Upload completed: index=0 src="./locales/de.json" process_id=batch-456`) {
			t.Fatalf("unexpected output: %q", gotOutput)
		}
	})

	t.Run("manifest load error", func(t *testing.T) {
		oldUploader := newUploaderFunc
		oldLoad := loadManifestFileFunc
		oldBuild := buildBatchUploadItemsFunc
		oldPerform := performBatchUploadFunc
		t.Cleanup(func() {
			newUploaderFunc = oldUploader
			loadManifestFileFunc = oldLoad
			buildBatchUploadItemsFunc = oldBuild
			performBatchUploadFunc = oldPerform
		})

		mu := &mockUploader{}
		newUploaderFunc = func(cfg *global_config.GlobalConfig) (uploader, error) {
			return mu, nil
		}

		loadManifestFileFunc = func(path string) (manifestFile, error) {
			return manifestFile{}, errors.New("bad manifest")
		}

		cfg := &global_config.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}
		flags := newFlags()

		cmd := newBoundTestCommand(flags)
		if err := cmd.Flags().Parse([]string{
			"--manifest=./manifest.json",
		}); err != nil {
			t.Fatalf("parse flags: %v", err)
		}

		err := runCommand(cmd, cfg, flags, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "bad manifest" {
			t.Fatalf("unexpected error: %q", err.Error())
		}
		if mu.uploadCalled {
			t.Fatal("expected Upload not to be called")
		}
		if mu.batchCalled {
			t.Fatal("expected UploadBatch not to be called")
		}
	})

	t.Run("build batch items error", func(t *testing.T) {
		oldUploader := newUploaderFunc
		oldLoad := loadManifestFileFunc
		oldBuild := buildBatchUploadItemsFunc
		oldPerform := performBatchUploadFunc
		t.Cleanup(func() {
			newUploaderFunc = oldUploader
			loadManifestFileFunc = oldLoad
			buildBatchUploadItemsFunc = oldBuild
			performBatchUploadFunc = oldPerform
		})

		mu := &mockUploader{}
		newUploaderFunc = func(cfg *global_config.GlobalConfig) (uploader, error) {
			return mu, nil
		}

		loadManifestFileFunc = func(path string) (manifestFile, error) {
			return manifestFile{
				Items: []manifestItem{
					{
						Params: lokexupload.UploadParams{
							"filename": "locales/en.json",
							"lang_iso": "en",
						},
					},
				},
			}, nil
		}

		buildBatchUploadItemsFunc = func(manifestPath string, mf manifestFile) ([]lokexupload.BatchUploadItem, error) {
			return nil, errors.New("invalid manifest item")
		}

		cfg := &global_config.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}
		flags := newFlags()

		cmd := newBoundTestCommand(flags)
		if err := cmd.Flags().Parse([]string{
			"--manifest=./manifest.json",
		}); err != nil {
			t.Fatalf("parse flags: %v", err)
		}

		err := runCommand(cmd, cfg, flags, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "invalid manifest item" {
			t.Fatalf("unexpected error: %q", err.Error())
		}
		if mu.uploadCalled {
			t.Fatal("expected Upload not to be called")
		}
		if mu.batchCalled {
			t.Fatal("expected UploadBatch not to be called")
		}
	})

	t.Run("batch upload error", func(t *testing.T) {
		oldUploader := newUploaderFunc
		oldLoad := loadManifestFileFunc
		oldBuild := buildBatchUploadItemsFunc
		oldPerform := performBatchUploadFunc
		t.Cleanup(func() {
			newUploaderFunc = oldUploader
			loadManifestFileFunc = oldLoad
			buildBatchUploadItemsFunc = oldBuild
			performBatchUploadFunc = oldPerform
		})

		mu := &mockUploader{}
		newUploaderFunc = func(cfg *global_config.GlobalConfig) (uploader, error) {
			return mu, nil
		}

		loadManifestFileFunc = func(path string) (manifestFile, error) {
			return manifestFile{
				Items: []manifestItem{
					{
						Params: lokexupload.UploadParams{
							"filename": "locales/en.json",
							"lang_iso": "en",
						},
						SrcPath: "./locales/en.json",
					},
				},
			}, nil
		}

		buildBatchUploadItemsFunc = func(manifestPath string, mf manifestFile) ([]lokexupload.BatchUploadItem, error) {
			return []lokexupload.BatchUploadItem{
				{
					Params: lokexupload.UploadParams{
						"filename": "locales/en.json",
						"lang_iso": "en",
					},
					SrcPath: "./locales/en.json",
				},
			}, nil
		}

		performBatchUploadFunc = func(
			ctx context.Context,
			up uploader,
			flags *Flags,
			items []lokexupload.BatchUploadItem,
		) (lokexupload.BatchUploadResult, error) {
			return lokexupload.BatchUploadResult{}, errors.New("batch upload failed")
		}

		cfg := &global_config.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}
		flags := newFlags()

		cmd := newBoundTestCommand(flags)
		if err := cmd.Flags().Parse([]string{
			"--manifest=./manifest.json",
		}); err != nil {
			t.Fatalf("parse flags: %v", err)
		}

		err := runCommand(cmd, cfg, flags, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != `upload batch from manifest "./manifest.json": batch upload failed` {
			t.Fatalf("unexpected error: %q", err.Error())
		}
		if mu.uploadCalled {
			t.Fatal("expected Upload not to be called")
		}
	})
}

func TestLoadManifestFile(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "manifest.json")

		content := `{
			"items": [
				{
					"params": {
						"filename": "locales/en.json",
						"lang_iso": "en"
					}
				},
				{
					"params": {
						"filename": "locales/de.json",
						"lang_iso": "de"
					},
					"src_path": "./de.json"
				}
			]
		}`

		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write manifest: %v", err)
		}

		got, err := loadManifestFile(path)
		if err != nil {
			t.Fatalf("loadManifestFile() error = %v", err)
		}

		if len(got.Items) != 2 {
			t.Fatalf("unexpected items len: got %d, want 2", len(got.Items))
		}
		if got.Items[0].Params["filename"] != "locales/en.json" {
			t.Fatalf("unexpected first filename: %#v", got.Items[0].Params)
		}
		if got.Items[1].Params["lang_iso"] != "de" {
			t.Fatalf("unexpected second lang_iso: %#v", got.Items[1].Params)
		}
		if got.Items[1].SrcPath != "./de.json" {
			t.Fatalf("unexpected second src_path: got %q", got.Items[1].SrcPath)
		}
	})

	t.Run("preserves numeric params as json.Number", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "manifest.json")

		content := `{
		"items": [
			{
				"params": {
					"filename": "locales/en.json",
					"lang_iso": "en",
					"filter_task_id": 1234567890123456789
				},
				"src_path": "./en.json"
			}
		]
	}`

		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write manifest: %v", err)
		}

		got, err := loadManifestFile(path)
		if err != nil {
			t.Fatalf("loadManifestFile() error = %v", err)
		}

		if len(got.Items) != 1 {
			t.Fatalf("unexpected items len: got %d, want 1", len(got.Items))
		}

		raw := got.Items[0].Params["filter_task_id"]

		num, ok := raw.(json.Number)
		if !ok {
			t.Fatalf("expected filter_task_id to be json.Number, got %T (%#v)", raw, raw)
		}

		if num.String() != "1234567890123456789" {
			t.Fatalf("unexpected filter_task_id: got %q", num.String())
		}
	})

	t.Run("file does not exist", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "missing.json")

		_, err := loadManifestFile(path)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), `read manifest file "`) {
			t.Fatalf("unexpected error: %q", err.Error())
		}
		if !strings.Contains(err.Error(), `missing.json`) {
			t.Fatalf("unexpected error: %q", err.Error())
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "manifest.json")

		content := `{
			"items": [
				{
					"params": {
						"filename": "locales/en.json",
						"lang_iso": "en"
					}
				}
			`
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write manifest: %v", err)
		}

		_, err := loadManifestFile(path)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), `parse manifest file "`) {
			t.Fatalf("unexpected error: %q", err.Error())
		}
	})

	t.Run("empty items", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "manifest.json")

		content := `{
			"items": []
		}`
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write manifest: %v", err)
		}

		_, err := loadManifestFile(path)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		want := fmt.Sprintf("manifest file %q contains no items", path)
		if err.Error() != want {
			t.Fatalf("unexpected error: got %q, want %q", err.Error(), want)
		}
	})

	t.Run("missing items field", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "manifest.json")

		content := `{
			"foo": "bar"
		}`
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write manifest: %v", err)
		}

		_, err := loadManifestFile(path)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		want := fmt.Sprintf("manifest file %q contains no items", path)
		if err.Error() != want {
			t.Fatalf("unexpected error: got %q, want %q", err.Error(), want)
		}
	})
}

func TestBuildBatchUploadItems(t *testing.T) {
	t.Run("ok with relative src_path", func(t *testing.T) {
		mf := manifestFile{
			Items: []manifestItem{
				{
					Params: lokexupload.UploadParams{
						"filename": "locales/en.json",
						"lang_iso": "en",
					},
					SrcPath: "en.json",
				},
			},
		}

		items, err := buildBatchUploadItems("/configs/manifest.json", mf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(items) != 1 {
			t.Fatalf("unexpected items len: %d", len(items))
		}

		wantPath := filepath.Join("/configs", "en.json")
		if items[0].SrcPath != wantPath {
			t.Fatalf("unexpected src path: got %q, want %q", items[0].SrcPath, wantPath)
		}
	})

	t.Run("ok with absolute src_path", func(t *testing.T) {
		absPath := filepath.Join(t.TempDir(), "en.json")

		mf := manifestFile{
			Items: []manifestItem{
				{
					Params: lokexupload.UploadParams{
						"filename": "locales/en.json",
						"lang_iso": "en",
					},
					SrcPath: absPath,
				},
			},
		}

		items, err := buildBatchUploadItems(filepath.Join("configs", "manifest.json"), mf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if items[0].SrcPath != absPath {
			t.Fatalf("unexpected src path: got %q, want %q", items[0].SrcPath, absPath)
		}
	})

	t.Run("ok without src_path", func(t *testing.T) {
		mf := manifestFile{
			Items: []manifestItem{
				{
					Params: lokexupload.UploadParams{
						"filename": "locales/en.json",
						"lang_iso": "en",
					},
				},
			},
		}

		items, err := buildBatchUploadItems("/configs/manifest.json", mf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if items[0].SrcPath != "" {
			t.Fatalf("expected empty src path, got %q", items[0].SrcPath)
		}
	})

	t.Run("error missing filename", func(t *testing.T) {
		mf := manifestFile{
			Items: []manifestItem{
				{
					Params: lokexupload.UploadParams{
						"lang_iso": "en",
					},
				},
			},
		}

		_, err := buildBatchUploadItems("/configs/manifest.json", mf)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "params.filename is required") {
			t.Fatalf("unexpected error: %q", err.Error())
		}
	})

	t.Run("error missing lang_iso", func(t *testing.T) {
		mf := manifestFile{
			Items: []manifestItem{
				{
					Params: lokexupload.UploadParams{
						"filename": "locales/en.json",
					},
				},
			},
		}

		_, err := buildBatchUploadItems("/configs/manifest.json", mf)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "params.lang_iso is required") {
			t.Fatalf("unexpected error: %q", err.Error())
		}
	})

	t.Run("error whitespace values", func(t *testing.T) {
		mf := manifestFile{
			Items: []manifestItem{
				{
					Params: lokexupload.UploadParams{
						"filename": "   ",
						"lang_iso": "   ",
					},
				},
			},
		}

		_, err := buildBatchUploadItems("/configs/manifest.json", mf)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "params.filename is required") &&
			!strings.Contains(err.Error(), "params.lang_iso is required") {
			t.Fatalf("unexpected error: %q", err.Error())
		}
	})
}

func TestPrintBatchItemResult(t *testing.T) {
	t.Run("prints failed item", func(t *testing.T) {
		cmd := &cobra.Command{}
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		item := lokexupload.BatchUploadResultItem{
			Index:   2,
			SrcPath: "./locales/de.json",
			Err:     errors.New("boom"),
		}

		printBatchItemResult(cmd, item, false)

		got := out.String()
		want := "Upload failed: index=2 src=\"./locales/de.json\" err=boom\n"
		if got != want {
			t.Fatalf("unexpected output:\n got: %q\nwant: %q", got, want)
		}
	})

	t.Run("prints started item when poll is false", func(t *testing.T) {
		cmd := &cobra.Command{}
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		item := lokexupload.BatchUploadResultItem{
			Index:     0,
			SrcPath:   "./locales/en.json",
			ProcessID: "process-123",
		}

		printBatchItemResult(cmd, item, false)

		got := out.String()
		want := "Upload started: index=0 src=\"./locales/en.json\" process_id=process-123\n"
		if got != want {
			t.Fatalf("unexpected output:\n got: %q\nwant: %q", got, want)
		}
	})

	t.Run("prints completed item when poll is true", func(t *testing.T) {
		cmd := &cobra.Command{}
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		item := lokexupload.BatchUploadResultItem{
			Index:     1,
			SrcPath:   "./locales/fr.json",
			ProcessID: "process-456",
		}

		printBatchItemResult(cmd, item, true)

		got := out.String()
		want := "Upload completed: index=1 src=\"./locales/fr.json\" process_id=process-456\n"
		if got != want {
			t.Fatalf("unexpected output:\n got: %q\nwant: %q", got, want)
		}
	})
}

func TestPrintBatchSummary(t *testing.T) {
	cmd := &cobra.Command{}
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	printBatchSummary(cmd, 3, 2, 1)

	got := out.String()
	want := "Batch summary: total=3 success=2 failed=1\n"
	if got != want {
		t.Fatalf("unexpected output:\n got: %q\nwant: %q", got, want)
	}
}

func TestPrintBatchUploadResult(t *testing.T) {
	t.Run("prints started items and summary when poll is false", func(t *testing.T) {
		cmd := &cobra.Command{}
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		result := lokexupload.BatchUploadResult{
			Items: []lokexupload.BatchUploadResultItem{
				{
					Index:     0,
					SrcPath:   "./locales/en.json",
					ProcessID: "process-1",
				},
				{
					Index:   1,
					SrcPath: "./locales/de.json",
					Err:     errors.New("failed upload"),
				},
			},
		}

		printBatchUploadResult(cmd, result, false)

		got := out.String()

		if !strings.Contains(got, `Upload started: index=0 src="./locales/en.json" process_id=process-1`) {
			t.Fatalf("unexpected output: %q", got)
		}
		if !strings.Contains(got, `Upload failed: index=1 src="./locales/de.json" err=failed upload`) {
			t.Fatalf("unexpected output: %q", got)
		}
		if !strings.Contains(got, "Batch summary: total=2 success=1 failed=1") {
			t.Fatalf("unexpected output: %q", got)
		}
	})

	t.Run("prints completed items and summary when poll is true", func(t *testing.T) {
		cmd := &cobra.Command{}
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		result := lokexupload.BatchUploadResult{
			Items: []lokexupload.BatchUploadResultItem{
				{
					Index:     0,
					SrcPath:   "./locales/en.json",
					ProcessID: "process-1",
				},
				{
					Index:     1,
					SrcPath:   "./locales/fr.json",
					ProcessID: "process-2",
				},
			},
		}

		printBatchUploadResult(cmd, result, true)

		got := out.String()

		if !strings.Contains(got, `Upload completed: index=0 src="./locales/en.json" process_id=process-1`) {
			t.Fatalf("unexpected output: %q", got)
		}
		if !strings.Contains(got, `Upload completed: index=1 src="./locales/fr.json" process_id=process-2`) {
			t.Fatalf("unexpected output: %q", got)
		}
		if !strings.Contains(got, "Batch summary: total=2 success=2 failed=0") {
			t.Fatalf("unexpected output: %q", got)
		}
	})

	t.Run("prints summary for empty result", func(t *testing.T) {
		cmd := &cobra.Command{}
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		result := lokexupload.BatchUploadResult{}

		printBatchUploadResult(cmd, result, false)

		got := out.String()
		want := "Batch summary: total=0 success=0 failed=0\n"
		if got != want {
			t.Fatalf("unexpected output:\n got: %q\nwant: %q", got, want)
		}
	})
}

func TestRequiredManifestString(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		params := lokexupload.UploadParams{
			"filename": "  en.json  ",
		}

		got, ok := requiredManifestString(params, "filename")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if got != "en.json" {
			t.Fatalf("unexpected value: got %q, want %q", got, "en.json")
		}
	})

	t.Run("missing key", func(t *testing.T) {
		params := lokexupload.UploadParams{}

		got, ok := requiredManifestString(params, "filename")
		if ok {
			t.Fatal("expected ok=false")
		}
		if got != "" {
			t.Fatalf("unexpected value: %q", got)
		}
	})

	t.Run("non string value", func(t *testing.T) {
		params := lokexupload.UploadParams{
			"filename": 123,
		}

		got, ok := requiredManifestString(params, "filename")
		if ok {
			t.Fatal("expected ok=false")
		}
		if got != "" {
			t.Fatalf("unexpected value: %q", got)
		}
	})

	t.Run("whitespace string", func(t *testing.T) {
		params := lokexupload.UploadParams{
			"filename": "   ",
		}

		got, ok := requiredManifestString(params, "filename")
		if ok {
			t.Fatal("expected ok=false")
		}
		if got != "" {
			t.Fatalf("unexpected value: %q", got)
		}
	})
}

func TestValidateManifestItem(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		item := manifestItem{
			Params: lokexupload.UploadParams{
				"filename": "en.json",
				"lang_iso": "en",
			},
		}

		err := validateManifestItem(item, 3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("missing params", func(t *testing.T) {
		item := manifestItem{}

		err := validateManifestItem(item, 1)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		want := "manifest item 1: params is required"
		if err.Error() != want {
			t.Fatalf("unexpected error: got %q, want %q", err.Error(), want)
		}
	})

	t.Run("missing filename", func(t *testing.T) {
		item := manifestItem{
			Params: lokexupload.UploadParams{
				"lang_iso": "en",
			},
		}

		err := validateManifestItem(item, 2)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		want := "manifest item 2: params.filename is required"
		if err.Error() != want {
			t.Fatalf("unexpected error: got %q, want %q", err.Error(), want)
		}
	})

	t.Run("missing lang_iso", func(t *testing.T) {
		item := manifestItem{
			Params: lokexupload.UploadParams{
				"filename": "en.json",
			},
		}

		err := validateManifestItem(item, 4)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		want := "manifest item 4: params.lang_iso is required"
		if err.Error() != want {
			t.Fatalf("unexpected error: got %q, want %q", err.Error(), want)
		}
	})

	t.Run("whitespace filename", func(t *testing.T) {
		item := manifestItem{
			Params: lokexupload.UploadParams{
				"filename": "   ",
				"lang_iso": "en",
			},
		}

		err := validateManifestItem(item, 5)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		want := "manifest item 5: params.filename is required"
		if err.Error() != want {
			t.Fatalf("unexpected error: got %q, want %q", err.Error(), want)
		}
	})

	t.Run("whitespace lang_iso", func(t *testing.T) {
		item := manifestItem{
			Params: lokexupload.UploadParams{
				"filename": "en.json",
				"lang_iso": "   ",
			},
		}

		err := validateManifestItem(item, 6)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		want := "manifest item 6: params.lang_iso is required"
		if err.Error() != want {
			t.Fatalf("unexpected error: got %q, want %q", err.Error(), want)
		}
	})

	t.Run("non string filename", func(t *testing.T) {
		item := manifestItem{
			Params: lokexupload.UploadParams{
				"filename": 123,
				"lang_iso": "en",
			},
		}

		err := validateManifestItem(item, 7)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		want := "manifest item 7: params.filename is required"
		if err.Error() != want {
			t.Fatalf("unexpected error: got %q, want %q", err.Error(), want)
		}
	})

	t.Run("non string lang_iso", func(t *testing.T) {
		item := manifestItem{
			Params: lokexupload.UploadParams{
				"filename": "en.json",
				"lang_iso": false,
			},
		}

		err := validateManifestItem(item, 8)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		want := "manifest item 8: params.lang_iso is required"
		if err.Error() != want {
			t.Fatalf("unexpected error: got %q, want %q", err.Error(), want)
		}
	})
}
