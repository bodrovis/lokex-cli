package languagemapping

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		raw     string
		want    []Mapping
		wantErr string
	}{
		{
			name: "empty",
			raw:  "",
			want: nil,
		},
		{
			name: "whitespace only",
			raw:  "   ",
			want: nil,
		},
		{
			name: "valid",
			raw: `[
				{
					"original_language_iso": "en-US",
					"custom_language_iso": "en"
				},
				{
					"original_language_iso": "pt-BR",
					"custom_language_iso": "pt"
				}
			]`,
			want: []Mapping{
				{
					OriginalLanguageISO: "en-US",
					CustomLanguageISO:   "en",
				},
				{
					OriginalLanguageISO: "pt-BR",
					CustomLanguageISO:   "pt",
				},
			},
		},
		{
			name: "trims values",
			raw: `[
				{
					"original_language_iso": "  en-US  ",
					"custom_language_iso": "  en  "
				}
			]`,
			want: []Mapping{
				{
					OriginalLanguageISO: "en-US",
					CustomLanguageISO:   "en",
				},
			},
		},
		{
			name: "missing original language",
			raw: `[
				{
					"custom_language_iso": "en"
				}
			]`,
			wantErr: "original_language_iso is required",
		},
		{
			name: "empty original language",
			raw: `[
				{
					"original_language_iso": "   ",
					"custom_language_iso": "en"
				}
			]`,
			wantErr: "original_language_iso is required",
		},
		{
			name: "missing custom language",
			raw: `[
				{
					"original_language_iso": "en-US"
				}
			]`,
			wantErr: "custom_language_iso is required",
		},
		{
			name: "empty custom language",
			raw: `[
				{
					"original_language_iso": "en-US",
					"custom_language_iso": "   "
				}
			]`,
			wantErr: "custom_language_iso is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Parse(tt.raw)

			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	t.Parallel()

	got, err := Parse(`{`)

	require.Nil(t, got)
	require.Error(t, err)
}

func TestParse_RejectsUnknownMembers(t *testing.T) {
	t.Parallel()

	got, err := Parse(`[
		{
			"original_language_iso": "en-US",
			"custom_language_iso": "en",
			"wat": true
		}
	]`)

	require.Nil(t, got)
	require.Error(t, err)
}

func TestToMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		mappings []Mapping
		want     map[string]string
	}{
		{
			name:     "nil",
			mappings: nil,
			want:     nil,
		},
		{
			name:     "empty",
			mappings: []Mapping{},
			want:     nil,
		},
		{
			name: "converts mappings",
			mappings: []Mapping{
				{
					OriginalLanguageISO: "en-US",
					CustomLanguageISO:   "en",
				},
				{
					OriginalLanguageISO: "pt-BR",
					CustomLanguageISO:   "pt",
				},
			},
			want: map[string]string{
				"en-US": "en",
				"pt-BR": "pt",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ToMap(tt.mappings)

			require.Equal(t, tt.want, got)
		})
	}
}

func TestToMap_DuplicateOriginalLanguageUsesLastValue(t *testing.T) {
	t.Parallel()

	got := ToMap(
		[]Mapping{
			{
				OriginalLanguageISO: "en-US",
				CustomLanguageISO:   "en",
			},
			{
				OriginalLanguageISO: "en-US",
				CustomLanguageISO:   "en_US",
			},
		},
	)

	require.Equal(
		t,
		map[string]string{
			"en-US": "en_US",
		},
		got,
	)
}
