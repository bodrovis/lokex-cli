package uploadmanifest

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterFiles(t *testing.T) {
	t.Parallel()

	files := []DiscoveredFile{
		{
			Path:         filepath.Join("locales", "en", "common.json"),
			RelativePath: filepath.Join("en", "common.json"),
		},
		{
			Path:         filepath.Join("locales", "en", "test.json"),
			RelativePath: filepath.Join("en", "test.json"),
		},
		{
			Path:         filepath.Join("locales", "de", "common.json"),
			RelativePath: filepath.Join("de", "common.json"),
		},
		{
			Path:         filepath.Join("locales", "vendor", "messages.json"),
			RelativePath: filepath.Join("vendor", "messages.json"),
		},
		{
			Path:         filepath.Join("locales", "en", "readme.txt"),
			RelativePath: filepath.Join("en", "readme.txt"),
		},
	}

	tests := []struct {
		name     string
		patterns []string
		want     []DiscoveredFile
	}{
		{
			name:     "no patterns",
			patterns: nil,
			want:     files,
		},
		{
			name: "empty patterns are ignored",
			patterns: []string{
				"",
				"   ",
			},
			want: files,
		},
		{
			name: "exclude single file",
			patterns: []string{
				"en/test.json",
			},
			want: []DiscoveredFile{
				files[0],
				files[2],
				files[3],
				files[4],
			},
		},
		{
			name: "exclude directory recursively",
			patterns: []string{
				"vendor/**",
			},
			want: []DiscoveredFile{
				files[0],
				files[1],
				files[2],
				files[4],
			},
		},
		{
			name: "exclude with globstar",
			patterns: []string{
				"**/test.json",
			},
			want: []DiscoveredFile{
				files[0],
				files[2],
				files[3],
				files[4],
			},
		},
		{
			name: "multiple patterns",
			patterns: []string{
				"vendor/**",
				"**/*.txt",
			},
			want: []DiscoveredFile{
				files[0],
				files[1],
				files[2],
			},
		},
		{
			name: "non matching pattern",
			patterns: []string{
				"**/*.yaml",
			},
			want: files,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := FilterFiles(
				files,
				tt.patterns,
			)

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFilterFiles_InvalidPattern(t *testing.T) {
	t.Parallel()

	files := []DiscoveredFile{
		{
			Path:         filepath.Join("locales", "en", "common.json"),
			RelativePath: filepath.Join("en", "common.json"),
		},
	}

	got, err := FilterFiles(
		files,
		[]string{
			"[",
		},
	)

	require.Nil(t, got)
	require.Error(t, err)
	require.Contains(
		t,
		err.Error(),
		`invalid exclude pattern "["`,
	)
}

func TestFilterFiles_NormalizesPatterns(t *testing.T) {
	t.Parallel()

	files := []DiscoveredFile{
		{
			Path:         filepath.Join("locales", "vendor", "common.json"),
			RelativePath: filepath.Join("vendor", "common.json"),
		},
		{
			Path:         filepath.Join("locales", "en", "common.json"),
			RelativePath: filepath.Join("en", "common.json"),
		},
	}

	got, err := FilterFiles(
		files,
		[]string{
			"  vendor/**  ",
		},
	)

	require.NoError(t, err)
	require.Equal(
		t,
		[]DiscoveredFile{
			files[1],
		},
		got,
	)
}

func TestFilterFiles_MatchesRelativePath(t *testing.T) {
	t.Parallel()

	files := []DiscoveredFile{
		{
			Path: filepath.Join(
				"some",
				"deep",
				"project",
				"locales",
				"en",
				"common.json",
			),
			RelativePath: filepath.Join(
				"en",
				"common.json",
			),
		},
	}

	got, err := FilterFiles(
		files,
		[]string{
			"en/**",
		},
	)

	require.NoError(t, err)
	require.Empty(t, got)
}
