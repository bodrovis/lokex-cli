package uploadmanifest

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchPathPattern(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		pattern string
		path    string
		want    PathParts
		matched bool
	}{
		{
			name:    "nested language directory",
			pattern: "{lang}/{name}.{ext}",
			path:    filepath.Join("en", "common.json"),
			want: PathParts{
				"lang": "en",
				"name": "common",
				"ext":  "json",
			},
			matched: true,
		},
		{
			name:    "flat filename",
			pattern: "{name}.{lang}.{ext}",
			path:    "common.de.json",
			want: PathParts{
				"lang": "de",
				"name": "common",
				"ext":  "json",
			},
			matched: true,
		},
		{
			name:    "single star matches one directory",
			pattern: "apps/*/locales/{lang}.{ext}",
			path:    "apps/admin/locales/en.json",
			want: PathParts{
				"lang": "en",
				"ext":  "json",
			},
			matched: true,
		},
		{
			name:    "single star does not cross directories",
			pattern: "apps/*/locales/{lang}.{ext}",
			path:    "apps/admin/frontend/locales/en.json",
			matched: false,
		},
		{
			name:    "globstar matches nested directories",
			pattern: "**/locales/{lang}.{ext}",
			path:    "apps/admin/frontend/locales/en.json",
			want: PathParts{
				"lang": "en",
				"ext":  "json",
			},
			matched: true,
		},
		{
			name:    "globstar matches one directory",
			pattern: "**/locales/{lang}.{ext}",
			path:    "apps/locales/de.json",
			want: PathParts{
				"lang": "de",
				"ext":  "json",
			},
			matched: true,
		},
		{
			name:    "globstar matches zero directories",
			pattern: "**/locales/{lang}.{ext}",
			path:    "locales/fr.json",
			want: PathParts{
				"lang": "fr",
				"ext":  "json",
			},
			matched: true,
		},
		{
			name:    "globstar after literal directory",
			pattern: "apps/**/{name}.{lang}.{ext}",
			path:    "apps/admin/locales/common.en.json",
			want: PathParts{
				"name": "common",
				"lang": "en",
				"ext":  "json",
			},
			matched: true,
		},
		{
			name:    "name contains dots",
			pattern: "{lang}/{name}.{ext}",
			path:    filepath.Join("en", "common.messages.json"),
			want: PathParts{
				"lang": "en",
				"name": "common.messages",
				"ext":  "json",
			},
			matched: true,
		},
		{
			name:    "pattern without language",
			pattern: "{name}.{ext}",
			path:    "common.json",
			want: PathParts{
				"name": "common",
				"ext":  "json",
			},
			matched: true,
		},
		{
			name:    "custom capture",
			pattern: "apps/{app}/locales/{lang}.{ext}",
			path:    "apps/mobile/locales/de.json",
			want: PathParts{
				"app":  "mobile",
				"lang": "de",
				"ext":  "json",
			},
			matched: true,
		},
		{
			name:    "custom capture with globstar",
			pattern: "**/{module}/locales/{lang}/{name}.{ext}",
			path:    "apps/frontend/mobile/locales/de/errors.json",
			want: PathParts{
				"module": "mobile",
				"lang":   "de",
				"name":   "errors",
				"ext":    "json",
			},
			matched: true,
		},
		{
			name:    "path does not match",
			pattern: "{lang}/{name}.{ext}",
			path:    "common.json",
			matched: false,
		},
		{
			name:    "literal dots are respected",
			pattern: "{name}.{lang}.{ext}",
			path:    "common-de-json",
			matched: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, matched, err := MatchPathPattern(
				tt.pattern,
				tt.path,
			)

			require.NoError(t, err)
			require.Equal(t, tt.matched, matched)

			if !tt.matched {
				require.Nil(t, got)
				return
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestMatchPathPattern_InvalidPattern(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		pattern string
		wantErr string
	}{
		{
			name:    "empty",
			pattern: "",
			wantErr: "name pattern is empty",
		},
		{
			name:    "whitespace only",
			pattern: "   ",
			wantErr: "name pattern is empty",
		},
		{
			name:    "nested opening brace",
			pattern: "{lang/{name}}.{ext}",
			wantErr: "unclosed pattern placeholder",
		},
		{
			name:    "unclosed placeholder",
			pattern: "{lang/{name}.{ext}",
			wantErr: "unclosed pattern placeholder",
		},
		{
			name:    "unexpected closing brace",
			pattern: "{lang}/{name}.{ext}}",
			wantErr: "unexpected closing brace in pattern",
		},
		{
			name:    "empty placeholder",
			pattern: "{lang}/{}.{ext}",
			wantErr: "invalid pattern placeholder {}",
		},
		{
			name:    "placeholder starts with number",
			pattern: "{lang}/{123name}.{ext}",
			wantErr: "invalid pattern placeholder {123name}",
		},
		{
			name:    "placeholder contains invalid character",
			pattern: "{lang}/{file.name}.{ext}",
			wantErr: "invalid pattern placeholder {file.name}",
		},
		{
			name:    "duplicate placeholder",
			pattern: "{lang}/{name}/{name}.{ext}",
			wantErr: "duplicate name pattern placeholder {name}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, _, err := MatchPathPattern(
				tt.pattern,
				"whatever",
			)

			require.EqualError(
				t,
				err,
				tt.wantErr,
			)
		})
	}
}

func TestRenderPattern(t *testing.T) {
	t.Parallel()

	parts := PathParts{
		"lang": "de",
		"name": "common",
		"ext":  "json",
	}

	tests := []struct {
		name    string
		pattern string
		want    string
	}{
		{
			name:    "default filename",
			pattern: "{name}.{ext}",
			want:    "common.json",
		},
		{
			name:    "language directory",
			pattern: "{lang}/{name}.{ext}",
			want:    "de/common.json",
		},
		{
			name:    "flat language suffix",
			pattern: "{name}.{lang}.{ext}",
			want:    "common.de.json",
		},
		{
			name:    "repeated placeholder",
			pattern: "{lang}/{lang}/{name}.{ext}",
			want:    "de/de/common.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(
				t,
				tt.want,
				RenderPattern(
					tt.pattern,
					parts,
				),
			)
		})
	}
}

func TestMatchPathPattern_CustomCapture(t *testing.T) {
	t.Parallel()

	got, matched, err := MatchPathPattern(
		"apps/{app}/locales/{lang}.{ext}",
		"apps/mobile/locales/de.json",
	)

	require.NoError(t, err)
	require.True(t, matched)

	require.Equal(
		t,
		PathParts{
			"app":  "mobile",
			"lang": "de",
			"ext":  "json",
		},
		got,
	)
}

func TestMatchPathPattern_CustomCaptureWithGlobstar(
	t *testing.T,
) {
	t.Parallel()

	got, matched, err := MatchPathPattern(
		"**/{module}/locales/{lang}/{name}.{ext}",
		"apps/frontend/mobile/locales/de/errors.json",
	)

	require.NoError(t, err)
	require.True(t, matched)

	require.Equal(
		t,
		PathParts{
			"module": "mobile",
			"lang":   "de",
			"name":   "errors",
			"ext":    "json",
		},
		got,
	)
}

func TestMatchAndRenderPatterns(t *testing.T) {
	t.Parallel()

	parts, matched, err := MatchPathPattern(
		"{lang}/{name}.{ext}",
		filepath.Join("fr", "messages.json"),
	)

	require.NoError(t, err)
	require.True(t, matched)

	require.Equal(
		t,
		PathParts{
			"lang": "fr",
			"name": "messages",
			"ext":  "json",
		},
		parts,
	)

	require.Equal(
		t,
		"messages.json",
		RenderPattern(
			"{name}.{ext}",
			parts,
		),
	)
}
