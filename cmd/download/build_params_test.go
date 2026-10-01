package download

import (
	"slices"
	"strings"
	"testing"

	"github.com/bodrovis/lokex-cli/internal/languagemapping"
	lokexdownload "github.com/bodrovis/lokex/v2/client/download"
)

func TestApplyLanguageMapping(t *testing.T) {
	t.Parallel()

	t.Run("nil value is ignored", func(t *testing.T) {
		t.Parallel()

		req := lokexdownload.DownloadParams{
			"format": "json",
		}

		err := applyLanguageMapping(req, nil)
		if err != nil {
			t.Fatalf("applyLanguageMapping() error = %v", err)
		}

		if _, ok := req["language_mapping"]; ok {
			t.Fatal("language_mapping was added, want it to be absent")
		}

		if got := req["format"]; got != "json" {
			t.Fatalf("format = %v, want json", got)
		}
	})

	t.Run("blank value is ignored", func(t *testing.T) {
		t.Parallel()

		raw := "   \n\t  "
		req := lokexdownload.DownloadParams{}

		err := applyLanguageMapping(req, &raw)
		if err != nil {
			t.Fatalf("applyLanguageMapping() error = %v", err)
		}

		if _, ok := req["language_mapping"]; ok {
			t.Fatal("language_mapping was added, want it to be absent")
		}
	})

	t.Run("valid mapping is added", func(t *testing.T) {
		t.Parallel()

		raw := `[
		{
			"original_language_iso": " en ",
			"custom_language_iso": " en_US "
		},
		{
			"original_language_iso": "pt",
			"custom_language_iso": "pt_BR"
		}
	]`

		req := lokexdownload.DownloadParams{
			"format": "json",
		}

		err := applyLanguageMapping(req, &raw)
		if err != nil {
			t.Fatalf("applyLanguageMapping() error = %v", err)
		}

		got, ok := req["language_mapping"].([]languagemapping.Mapping)
		if !ok {
			t.Fatalf(
				"language_mapping type = %T, want []languagemapping.Mapping",
				req["language_mapping"],
			)
		}

		want := []languagemapping.Mapping{
			{
				OriginalLanguageISO: "en",
				CustomLanguageISO:   "en_US",
			},
			{
				OriginalLanguageISO: "pt",
				CustomLanguageISO:   "pt_BR",
			},
		}

		if !slices.Equal(got, want) {
			t.Fatalf("language_mapping = %#v, want %#v", got, want)
		}

		if got := req["format"]; got != "json" {
			t.Fatalf("format = %v, want json", got)
		}
	})

	t.Run("invalid mapping returns wrapped error", func(t *testing.T) {
		t.Parallel()

		raw := `[{`
		req := lokexdownload.DownloadParams{}

		err := applyLanguageMapping(req, &raw)
		if err == nil {
			t.Fatal("applyLanguageMapping() error = nil, want error")
		}

		if !strings.Contains(err.Error(), "parse language-mapping") {
			t.Fatalf(
				"error = %q, want it to contain %q",
				err,
				"parse language-mapping",
			)
		}

		if _, ok := req["language_mapping"]; ok {
			t.Fatal("language_mapping was added after parse failure")
		}
	})
}
