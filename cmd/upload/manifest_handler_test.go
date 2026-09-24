package upload

import (
	"bytes"
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/bodrovis/lokex-cli/internal/global_config"
	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		got, err := performBatchUpload(
			context.Background(),
			mu,
			true,
			items,
		)

		require.NoError(t, err)

		require.True(t, mu.batchCalled)
		require.NotNil(t, mu.gotBatchCtx)
		require.True(t, mu.gotBatchPoll)

		require.Len(t, mu.gotBatchItems, 2)
		assert.Equal(t, "./locales/en.json", mu.gotBatchItems[0].SrcPath)
		assert.Equal(t, "de", mu.gotBatchItems[1].Params["lang_iso"])

		require.Len(t, got.Items, 2)
		assert.Equal(t, "process-1", got.Items[0].ProcessID)
	})

	t.Run("error", func(t *testing.T) {
		wantErr := errors.New("batch upload failed")

		mu := &mockUploader{
			batchErr: wantErr,
		}

		_, err := performBatchUpload(
			context.Background(),
			mu,
			false,
			items,
		)

		require.ErrorIs(t, err, wantErr)

		require.True(t, mu.batchCalled)
		require.False(t, mu.gotBatchPoll)
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

		newUploaderFunc = func(
			*global_config.GlobalConfig,
		) (uploader, error) {
			return mu, nil
		}

		loadManifestFileFunc = func(path string) (manifestFile, error) {
			require.Equal(t, "./manifest.json", path)

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

		buildBatchUploadItemsFunc = func(
			manifestPath string,
			mf manifestFile,
		) ([]lokexupload.BatchUploadItem, error) {
			require.Equal(t, "./manifest.json", manifestPath)
			require.Len(t, mf.Items, 1)

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
			poll bool,
			items []lokexupload.BatchUploadItem,
		) (lokexupload.BatchUploadResult, error) {
			require.NotNil(t, ctx)
			require.Equal(t, mu, up)
			require.False(t, poll)
			require.Len(t, items, 1)

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

		uploadCfg := newManifestUploadConfig(
			"./manifest.json",
			false,
		)

		cmd := &cobra.Command{Use: "upload"}

		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		err := runCommand(
			cmd,
			cfg,
			uploadCfg,
		)
		require.NoError(t, err)

		require.False(t, mu.uploadCalled)

		assert.Contains(
			t,
			out.String(),
			`Upload started: index=0 src="./locales/en.json" process_id=batch-123`,
		)

		assert.Contains(
			t,
			out.String(),
			"Batch summary: total=1 success=1 failed=0",
		)
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

		newUploaderFunc = func(
			*global_config.GlobalConfig,
		) (uploader, error) {
			return mu, nil
		}

		loadManifestFileFunc = func(string) (manifestFile, error) {
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

		buildBatchUploadItemsFunc = func(
			string,
			manifestFile,
		) ([]lokexupload.BatchUploadItem, error) {
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
			_ context.Context,
			_ uploader,
			poll bool,
			_ []lokexupload.BatchUploadItem,
		) (lokexupload.BatchUploadResult, error) {
			require.True(t, poll)

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

		uploadCfg := newManifestUploadConfig(
			"./manifest.json",
			true,
		)

		cmd := &cobra.Command{Use: "upload"}

		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		err := runCommand(
			cmd,
			cfg,
			uploadCfg,
		)
		require.NoError(t, err)

		require.False(t, mu.uploadCalled)

		assert.Contains(
			t,
			out.String(),
			`Upload completed: index=0 src="./locales/de.json" process_id=batch-456`,
		)
	})

	t.Run("manifest load error", func(t *testing.T) {
		oldUploader := newUploaderFunc
		oldLoad := loadManifestFileFunc

		t.Cleanup(func() {
			newUploaderFunc = oldUploader
			loadManifestFileFunc = oldLoad
		})

		mu := &mockUploader{}

		newUploaderFunc = func(
			*global_config.GlobalConfig,
		) (uploader, error) {
			return mu, nil
		}

		wantErr := errors.New("bad manifest")

		loadManifestFileFunc = func(string) (manifestFile, error) {
			return manifestFile{}, wantErr
		}

		cfg := &global_config.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}

		uploadCfg := newManifestUploadConfig(
			"./manifest.json",
			false,
		)

		cmd := &cobra.Command{Use: "upload"}

		err := runCommand(
			cmd,
			cfg,
			uploadCfg,
		)

		require.ErrorIs(t, err, wantErr)
		require.False(t, mu.uploadCalled)
		require.False(t, mu.batchCalled)
	})

	t.Run("build batch items error", func(t *testing.T) {
		oldUploader := newUploaderFunc
		oldLoad := loadManifestFileFunc
		oldBuild := buildBatchUploadItemsFunc

		t.Cleanup(func() {
			newUploaderFunc = oldUploader
			loadManifestFileFunc = oldLoad
			buildBatchUploadItemsFunc = oldBuild
		})

		mu := &mockUploader{}

		newUploaderFunc = func(
			*global_config.GlobalConfig,
		) (uploader, error) {
			return mu, nil
		}

		loadManifestFileFunc = func(string) (manifestFile, error) {
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

		wantErr := errors.New("invalid manifest item")

		buildBatchUploadItemsFunc = func(
			string,
			manifestFile,
		) ([]lokexupload.BatchUploadItem, error) {
			return nil, wantErr
		}

		cfg := &global_config.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}

		uploadCfg := newManifestUploadConfig(
			"./manifest.json",
			false,
		)

		cmd := &cobra.Command{Use: "upload"}

		err := runCommand(
			cmd,
			cfg,
			uploadCfg,
		)

		require.ErrorIs(t, err, wantErr)
		require.False(t, mu.uploadCalled)
		require.False(t, mu.batchCalled)
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

		newUploaderFunc = func(
			*global_config.GlobalConfig,
		) (uploader, error) {
			return mu, nil
		}

		loadManifestFileFunc = func(string) (manifestFile, error) {
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

		buildBatchUploadItemsFunc = func(
			string,
			manifestFile,
		) ([]lokexupload.BatchUploadItem, error) {
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

		wantErr := errors.New("batch upload failed")

		performBatchUploadFunc = func(
			context.Context,
			uploader,
			bool,
			[]lokexupload.BatchUploadItem,
		) (lokexupload.BatchUploadResult, error) {
			return lokexupload.BatchUploadResult{}, wantErr
		}

		cfg := &global_config.GlobalConfig{
			Token:     "token",
			ProjectID: "project-id",
		}

		uploadCfg := newManifestUploadConfig(
			"./manifest.json",
			false,
		)

		cmd := &cobra.Command{Use: "upload"}

		err := runCommand(
			cmd,
			cfg,
			uploadCfg,
		)

		require.EqualError(
			t,
			err,
			`upload batch from manifest "./manifest.json": batch upload failed`,
		)

		require.False(t, mu.uploadCalled)
	})
}

func TestBuildBatchUploadItems(t *testing.T) {
	t.Run("relative src_path is resolved against manifest directory", func(t *testing.T) {
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

		items, err := buildBatchUploadItems(
			"/configs/manifest.json",
			mf,
		)
		require.NoError(t, err)
		require.Len(t, items, 1)

		assert.Equal(
			t,
			filepath.Join("/configs", "en.json"),
			items[0].SrcPath,
		)
	})

	t.Run("absolute src_path is preserved", func(t *testing.T) {
		absPath := filepath.Join(
			t.TempDir(),
			"en.json",
		)

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

		items, err := buildBatchUploadItems(
			filepath.Join("configs", "manifest.json"),
			mf,
		)
		require.NoError(t, err)
		require.Len(t, items, 1)

		assert.Equal(t, absPath, items[0].SrcPath)
	})

	t.Run("empty src_path is preserved", func(t *testing.T) {
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

		items, err := buildBatchUploadItems(
			"/configs/manifest.json",
			mf,
		)
		require.NoError(t, err)
		require.Len(t, items, 1)

		assert.Empty(t, items[0].SrcPath)
	})

	tests := []struct {
		name   string
		params lokexupload.UploadParams
		want   string
	}{
		{
			name: "missing filename",
			params: lokexupload.UploadParams{
				"lang_iso": "en",
			},
			want: "params.filename is required",
		},
		{
			name: "missing lang_iso",
			params: lokexupload.UploadParams{
				"filename": "locales/en.json",
			},
			want: "params.lang_iso is required",
		},
		{
			name: "whitespace filename",
			params: lokexupload.UploadParams{
				"filename": "   ",
				"lang_iso": "en",
			},
			want: "params.filename is required",
		},
		{
			name: "whitespace lang_iso",
			params: lokexupload.UploadParams{
				"filename": "locales/en.json",
				"lang_iso": "   ",
			},
			want: "params.lang_iso is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mf := manifestFile{
				Items: []manifestItem{
					{
						Params: tt.params,
					},
				},
			}

			_, err := buildBatchUploadItems(
				"/configs/manifest.json",
				mf,
			)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.want)
		})
	}
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

func newManifestUploadConfig(
	path string,
	poll bool,
) *UploadConfig {
	return &UploadConfig{
		Manifest: &path,
		Poll:     &poll,
	}
}
