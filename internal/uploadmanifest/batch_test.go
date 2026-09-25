package uploadmanifest

import (
	"path/filepath"
	"testing"

	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
	"github.com/stretchr/testify/require"
)

func TestBuildBatchUploadItems(t *testing.T) {
	t.Parallel()

	manifestPath := filepath.Join(
		"project",
		"manifests",
		"upload.json",
	)

	mf := File{
		Items: []Item{
			{
				SrcPath: filepath.Join(
					"..",
					"locales",
					"en",
					"common.json",
				),
				Params: lokexupload.UploadParams{
					"filename": "common.json",
					"lang_iso": "en",
				},
			},
		},
	}

	got, err := BuildBatchUploadItems(
		manifestPath,
		mf,
	)
	require.NoError(t, err)

	require.Equal(
		t,
		[]lokexupload.BatchUploadItem{
			{
				SrcPath: filepath.Join(
					"project",
					"manifests",
					"..",
					"locales",
					"en",
					"common.json",
				),
				Params: lokexupload.UploadParams{
					"filename": "common.json",
					"lang_iso": "en",
				},
			},
		},
		got,
	)
}

func TestBuildBatchUploadItems_PreservesAbsoluteSourcePath(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	sourcePath := filepath.Join(
		root,
		"locales",
		"en.json",
	)

	manifestPath := filepath.Join(
		root,
		"manifests",
		"upload.json",
	)

	mf := File{
		Items: []Item{
			{
				SrcPath: sourcePath,
				Params: lokexupload.UploadParams{
					"filename": "messages.json",
					"lang_iso": "en",
				},
			},
		},
	}

	got, err := BuildBatchUploadItems(
		manifestPath,
		mf,
	)
	require.NoError(t, err)
	require.Len(t, got, 1)

	require.Equal(
		t,
		sourcePath,
		got[0].SrcPath,
	)
}

func TestBuildBatchUploadItems_PreservesEmptySourcePath(t *testing.T) {
	t.Parallel()

	mf := File{
		Items: []Item{
			{
				Params: lokexupload.UploadParams{
					"filename": "messages.json",
					"lang_iso": "en",
				},
			},
		},
	}

	got, err := BuildBatchUploadItems(
		filepath.Join("manifests", "upload.json"),
		mf,
	)
	require.NoError(t, err)
	require.Len(t, got, 1)

	require.Empty(
		t,
		got[0].SrcPath,
	)
}

func TestBuildBatchUploadItems_ValidatesItems(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		params  lokexupload.UploadParams
		wantErr string
	}{
		{
			name:    "missing params",
			params:  nil,
			wantErr: "manifest item 1: params is required",
		},
		{
			name: "missing filename",
			params: lokexupload.UploadParams{
				"lang_iso": "en",
			},
			wantErr: "manifest item 1: params.filename is required",
		},
		{
			name: "empty filename",
			params: lokexupload.UploadParams{
				"filename": "",
				"lang_iso": "en",
			},
			wantErr: "manifest item 1: params.filename is required",
		},
		{
			name: "whitespace filename",
			params: lokexupload.UploadParams{
				"filename": "   ",
				"lang_iso": "en",
			},
			wantErr: "manifest item 1: params.filename is required",
		},
		{
			name: "filename has wrong type",
			params: lokexupload.UploadParams{
				"filename": 123,
				"lang_iso": "en",
			},
			wantErr: "manifest item 1: params.filename is required",
		},
		{
			name: "missing language",
			params: lokexupload.UploadParams{
				"filename": "messages.json",
			},
			wantErr: "manifest item 1: params.lang_iso is required",
		},
		{
			name: "empty language",
			params: lokexupload.UploadParams{
				"filename": "messages.json",
				"lang_iso": "",
			},
			wantErr: "manifest item 1: params.lang_iso is required",
		},
		{
			name: "whitespace language",
			params: lokexupload.UploadParams{
				"filename": "messages.json",
				"lang_iso": "   ",
			},
			wantErr: "manifest item 1: params.lang_iso is required",
		},
		{
			name: "language has wrong type",
			params: lokexupload.UploadParams{
				"filename": "messages.json",
				"lang_iso": 123,
			},
			wantErr: "manifest item 1: params.lang_iso is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mf := File{
				Items: []Item{
					{
						Params: tt.params,
					},
				},
			}

			got, err := BuildBatchUploadItems(
				"upload.json",
				mf,
			)

			require.Nil(t, got)
			require.EqualError(
				t,
				err,
				tt.wantErr,
			)
		})
	}
}

func TestBuildBatchUploadItems_ReportsItemIndex(t *testing.T) {
	t.Parallel()

	mf := File{
		Items: []Item{
			{
				Params: lokexupload.UploadParams{
					"filename": "first.json",
					"lang_iso": "en",
				},
			},
			{
				Params: lokexupload.UploadParams{
					"filename": "second.json",
				},
			},
		},
	}

	got, err := BuildBatchUploadItems(
		"upload.json",
		mf,
	)

	require.Nil(t, got)

	require.EqualError(
		t,
		err,
		"manifest item 2: params.lang_iso is required",
	)
}

func TestBuildBatchUploadItems_EmptyManifest(t *testing.T) {
	t.Parallel()

	got, err := BuildBatchUploadItems(
		"upload.json",
		File{},
	)

	require.NoError(t, err)
	require.Empty(t, got)
}
