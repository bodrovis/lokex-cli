package uploadmanifest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bodrovis/lokex-cli/internal/languagemapping"
	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
	"github.com/stretchr/testify/require"
)

func TestGenerate_NestedLayout(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	enDir := filepath.Join(root, "locales", "en")
	deDir := filepath.Join(root, "locales", "de")

	require.NoError(t, os.MkdirAll(enDir, 0o755))
	require.NoError(t, os.MkdirAll(deDir, 0o755))

	enFile := filepath.Join(enDir, "common.json")
	deFile := filepath.Join(deDir, "common.json")

	require.NoError(t, os.WriteFile(enFile, []byte(`{}`), 0o644))
	require.NoError(t, os.WriteFile(deFile, []byte(`{}`), 0o644))

	manifestDir := filepath.Join(root, "manifests")

	mf, err := Generate(
		GenerateOptions{
			Paths: []string{
				filepath.Join(root, "locales"),
			},
			NamePattern:     "{lang}/{name}.{ext}",
			FilenamePattern: "{name}.{ext}",
			ManifestDir:     manifestDir,
		},
	)
	require.NoError(t, err)
	require.Len(t, mf.Items, 2)

	require.Equal(
		t,
		"common.json",
		mf.Items[0].Params["filename"],
	)

	require.Equal(
		t,
		"de",
		mf.Items[0].Params["lang_iso"],
	)

	require.Equal(
		t,
		filepath.ToSlash(
			filepath.Join("..", "locales", "de", "common.json"),
		),
		mf.Items[0].SrcPath,
	)

	require.Equal(
		t,
		"en",
		mf.Items[1].Params["lang_iso"],
	)
}

func TestGenerate_FlatLayout(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(root, "common.en.json"),
			[]byte(`{}`),
			0o644,
		),
	)

	mf, err := Generate(
		GenerateOptions{
			Paths: []string{
				root,
			},
			NamePattern:     "{name}.{lang}.{ext}",
			FilenamePattern: "{name}.{ext}",
			ManifestDir:     root,
		},
	)
	require.NoError(t, err)
	require.Len(t, mf.Items, 1)

	require.Equal(
		t,
		"common.json",
		mf.Items[0].Params["filename"],
	)

	require.Equal(
		t,
		"en",
		mf.Items[0].Params["lang_iso"],
	)

	require.Equal(
		t,
		"common.en.json",
		mf.Items[0].SrcPath,
	)
}

func TestGenerate_UsesBaseLanguage(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(root, "common.json"),
			[]byte(`{}`),
			0o644,
		),
	)

	mf, err := Generate(
		GenerateOptions{
			Paths: []string{
				root,
			},
			NamePattern:     "{name}.{ext}",
			FilenamePattern: "{name}.{ext}",
			BaseLang:        "en",
			ManifestDir:     root,
		},
	)
	require.NoError(t, err)
	require.Len(t, mf.Items, 1)

	require.Equal(
		t,
		"en",
		mf.Items[0].Params["lang_iso"],
	)
}

func TestGenerate_AppliesLanguageMapping(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	langDir := filepath.Join(
		root,
		"en-US",
	)

	require.NoError(
		t,
		os.MkdirAll(langDir, 0o755),
	)

	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(langDir, "common.json"),
			[]byte(`{}`),
			0o644,
		),
	)

	mf, err := Generate(
		GenerateOptions{
			Paths: []string{
				root,
			},
			NamePattern:     "{lang}/{name}.{ext}",
			FilenamePattern: "{lang}/{name}.{ext}",
			LanguageMapping: []languagemapping.Mapping{
				{
					OriginalLanguageISO: "en-US",
					CustomLanguageISO:   "en",
				},
			},
			ManifestDir: root,
		},
	)
	require.NoError(t, err)
	require.Len(t, mf.Items, 1)

	require.Equal(
		t,
		"en",
		mf.Items[0].Params["lang_iso"],
	)

	// Mapping changes lang_iso, not the original {lang}
	// captured from the path.
	require.Equal(
		t,
		"en-US/common.json",
		mf.Items[0].Params["filename"],
	)
}

func TestGenerate_ExcludesFiles(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	enDir := filepath.Join(root, "en")
	require.NoError(t, os.MkdirAll(enDir, 0o755))

	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(enDir, "common.json"),
			[]byte(`{}`),
			0o644,
		),
	)

	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(enDir, "common.test.json"),
			[]byte(`{}`),
			0o644,
		),
	)

	mf, err := Generate(
		GenerateOptions{
			Paths: []string{
				root,
			},
			ExcludePatterns: []string{
				"**/*.test.json",
			},
			NamePattern:     "{lang}/{name}.{ext}",
			FilenamePattern: "{name}.{ext}",
			ManifestDir:     root,
		},
	)
	require.NoError(t, err)
	require.Len(t, mf.Items, 1)

	require.Equal(
		t,
		"common.json",
		mf.Items[0].Params["filename"],
	)
}

func TestGenerate_Validation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		opts    GenerateOptions
		wantErr string
	}{
		{
			name:    "no paths",
			opts:    GenerateOptions{},
			wantErr: "no paths provided",
		},
		{
			name: "missing name pattern",
			opts: GenerateOptions{
				Paths:           []string{"whatever"},
				FilenamePattern: "{name}.{ext}",
			},
			wantErr: "name pattern is required",
		},
		{
			name: "missing filename pattern",
			opts: GenerateOptions{
				Paths:       []string{"whatever"},
				NamePattern: "{lang}/{name}.{ext}",
			},
			wantErr: "filename pattern is required",
		},
		{
			name: "base language required",
			opts: GenerateOptions{
				Paths:           []string{"whatever"},
				NamePattern:     "{name}.{ext}",
				FilenamePattern: "{name}.{ext}",
			},
			wantErr: "base language is required when name pattern does not contain {lang}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := Generate(tt.opts)

			require.EqualError(
				t,
				err,
				tt.wantErr,
			)
		})
	}
}

func TestGenerate_NoMatchingFiles(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(root, "common.json"),
			[]byte(`{}`),
			0o644,
		),
	)

	_, err := Generate(
		GenerateOptions{
			Paths: []string{
				root,
			},
			NamePattern:     "{lang}/{name}.{ext}",
			FilenamePattern: "{name}.{ext}",
			ManifestDir:     root,
		},
	)

	require.EqualError(
		t,
		err,
		`no files matched name pattern "{lang}/{name}.{ext}"`,
	)
}

func TestGenerate_InvalidExcludePattern(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(root, "common.en.json"),
			[]byte(`{}`),
			0o644,
		),
	)

	_, err := Generate(
		GenerateOptions{
			Paths: []string{
				root,
			},
			ExcludePatterns: []string{
				"[",
			},
			NamePattern:     "{name}.{lang}.{ext}",
			FilenamePattern: "{name}.{ext}",
			ManifestDir:     root,
		},
	)

	require.Error(t, err)
	require.Contains(
		t,
		err.Error(),
		"invalid exclude pattern",
	)
}

func TestGenerate_CustomCaptureWithGlobstar(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	files := []string{
		filepath.Join(
			root,
			"apps",
			"mobile",
			"locales",
			"de.json",
		),
		filepath.Join(
			root,
			"apps",
			"web",
			"locales",
			"en.json",
		),
		filepath.Join(
			root,
			"archive",
			"apps",
			"admin",
			"locales",
			"fr.json",
		),
	}

	for _, path := range files {
		require.NoError(
			t,
			os.MkdirAll(
				filepath.Dir(path),
				0o755,
			),
		)

		require.NoError(
			t,
			os.WriteFile(
				path,
				[]byte(`{}`),
				0o644,
			),
		)
	}

	mf, err := Generate(
		GenerateOptions{
			Paths: []string{
				root,
			},
			NamePattern:     "**/apps/{app}/locales/{lang}.{ext}",
			FilenamePattern: "{app}/{lang}.{ext}",
			ManifestDir:     root,
		},
	)
	require.NoError(t, err)

	require.Equal(
		t,
		File{
			Items: []Item{
				{
					Params: lokexupload.UploadParams{
						"filename": "mobile/de.json",
						"lang_iso": "de",
					},
					SrcPath: "apps/mobile/locales/de.json",
				},
				{
					Params: lokexupload.UploadParams{
						"filename": "web/en.json",
						"lang_iso": "en",
					},
					SrcPath: "apps/web/locales/en.json",
				},
				{
					Params: lokexupload.UploadParams{
						"filename": "admin/fr.json",
						"lang_iso": "fr",
					},
					SrcPath: "archive/apps/admin/locales/fr.json",
				},
			},
		},
		mf,
	)
}

func TestRelativeSourcePath(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	manifestDir := filepath.Join(
		root,
		"manifests",
	)

	sourcePath := filepath.Join(
		root,
		"locales",
		"en",
		"common.json",
	)

	got, err := relativeSourcePath(
		manifestDir,
		sourcePath,
	)

	require.NoError(t, err)

	require.Equal(
		t,
		"../locales/en/common.json",
		got,
	)
}

func TestSameVolume(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	require.True(
		t,
		sameVolume(
			root,
			filepath.Join(
				root,
				"locales",
			),
		),
	)
}
