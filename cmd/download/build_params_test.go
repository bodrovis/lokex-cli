package download

import (
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestBuildParams_Minimal(t *testing.T) {
	t.Parallel()

	cmd := newTestCommand()
	flags := &Flags{
		Format: "json",
	}

	params, err := buildParams(cmd, flags, nil)
	if err != nil {
		t.Fatalf("buildParams() error = %v", err)
	}

	if got := params["format"]; got != "json" {
		t.Fatalf("expected format %q, got %v", "json", got)
	}

	if len(params) != 1 {
		t.Fatalf("expected exactly 1 param, got %d: %#v", len(params), params)
	}
}

func TestBuildParams_TrimsStringFields(t *testing.T) {
	t.Parallel()

	cmd := newTestCommand()
	flags := &Flags{
		Format:          "json",
		BundleStructure: " %LANG_ISO%.json ",
		DirectoryPrefix: " locale/ ",
	}

	params, err := buildParams(cmd, flags, nil)
	if err != nil {
		t.Fatalf("buildParams() error = %v", err)
	}

	if got := params["bundle_structure"]; got != "%LANG_ISO%.json" {
		t.Fatalf("unexpected bundle_structure: got %v", got)
	}

	if got := params["directory_prefix"]; got != "locale/" {
		t.Fatalf("unexpected directory_prefix: got %v", got)
	}
}

func TestBuildParams_FiltersStringSliceFields(t *testing.T) {
	t.Parallel()

	cmd := newTestCommand()
	flags := &Flags{
		Format:      "json",
		FilterLangs: []string{"en", "", "   ", "fr"},
	}

	params, err := buildParams(cmd, flags, nil)
	if err != nil {
		t.Fatalf("buildParams() error = %v", err)
	}

	got, ok := params["filter_langs"]
	if !ok {
		t.Fatal("expected filter_langs to be set")
	}

	want := []string{"en", "fr"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected filter_langs: got %#v, want %#v", got, want)
	}
}

func TestBuildParams_OmitsEmptyStringSliceFields(t *testing.T) {
	t.Parallel()

	cmd := newTestCommand()
	flags := &Flags{
		Format:      "json",
		FilterLangs: []string{"", "   "},
	}

	params, err := buildParams(cmd, flags, nil)
	if err != nil {
		t.Fatalf("buildParams() error = %v", err)
	}

	if _, ok := params["filter_langs"]; ok {
		t.Fatalf("expected filter_langs to be omitted, got %#v", params["filter_langs"])
	}
}

func TestBuildParams_SetsBoolOnlyWhenFlagChanged(t *testing.T) {
	t.Parallel()

	t.Run("omits bool when flag not changed", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()
		flags := &Flags{
			Format:  "json",
			Compact: true,
		}

		params, err := buildParams(cmd, flags, nil)
		if err != nil {
			t.Fatalf("buildParams() error = %v", err)
		}

		if _, ok := params["compact"]; ok {
			t.Fatalf("expected compact to be omitted, got %#v", params["compact"])
		}
	})

	t.Run("sets bool true when explicitly enabled", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()
		if err := cmd.Flags().Set("compact", "true"); err != nil {
			t.Fatalf("set compact: %v", err)
		}

		flags := &Flags{
			Format:  "json",
			Compact: true,
		}

		params, err := buildParams(cmd, flags, nil)
		if err != nil {
			t.Fatalf("buildParams() error = %v", err)
		}

		got, ok := params["compact"]
		if !ok {
			t.Fatal("expected compact to be set")
		}
		if got != true {
			t.Fatalf("expected compact=true, got %#v", got)
		}
	})

	t.Run("sets bool false when explicitly disabled", func(t *testing.T) {
		t.Parallel()

		cmd := newBoolDefaultTrueCommand()
		if err := cmd.Flags().Set("compact", "false"); err != nil {
			t.Fatalf("set compact: %v", err)
		}

		flags := &Flags{
			Format:  "json",
			Compact: false,
		}

		params, err := buildParams(cmd, flags, nil)
		if err != nil {
			t.Fatalf("buildParams() error = %v", err)
		}

		got, ok := params["compact"]
		if !ok {
			t.Fatal("expected compact to be set")
		}
		if got != false {
			t.Fatalf("expected compact=false, got %#v", got)
		}
	})
}

func TestBuildParams_FilterTaskID(t *testing.T) {
	t.Parallel()

	t.Run("omits filter_task_id when flag not changed", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()
		flags := &Flags{
			Format:       "json",
			FilterTaskID: 123,
		}

		params, err := buildParams(cmd, flags, nil)
		if err != nil {
			t.Fatalf("buildParams() error = %v", err)
		}

		if _, ok := params["filter_task_id"]; ok {
			t.Fatalf("expected filter_task_id to be omitted, got %#v", params["filter_task_id"])
		}
	})

	t.Run("sets filter_task_id when flag changed", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()
		if err := cmd.Flags().Set("filter-task-id", "123"); err != nil {
			t.Fatalf("set filter-task-id: %v", err)
		}

		flags := &Flags{
			Format:       "json",
			FilterTaskID: 123,
		}

		params, err := buildParams(cmd, flags, nil)
		if err != nil {
			t.Fatalf("buildParams() error = %v", err)
		}

		got, ok := params["filter_task_id"]
		if !ok {
			t.Fatal("expected filter_task_id to be set")
		}
		if got != flags.FilterTaskID {
			t.Fatalf("expected filter_task_id=%v, got %#v", flags.FilterTaskID, got)
		}
	})
}

func TestBuildParams_LanguageMapping(t *testing.T) {
	t.Parallel()

	t.Run("sets parsed language_mapping for valid json", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()
		flags := &Flags{
			Format:              "json",
			LanguageMappingJSON: `[{"lang_iso":"en","custom_iso":"en_US"}]`,
		}

		params, err := buildParams(cmd, flags, nil)
		if err != nil {
			t.Fatalf("buildParams() error = %v", err)
		}

		got, ok := params["language_mapping"]
		if !ok {
			t.Fatal("expected language_mapping to be set")
		}

		want := []map[string]any{
			{
				"lang_iso":   "en",
				"custom_iso": "en_US",
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected language_mapping: got %#v, want %#v", got, want)
		}
	})

	t.Run("returns error for invalid json", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()
		flags := &Flags{
			Format:              "json",
			LanguageMappingJSON: `not-json`,
		}

		_, err := buildParams(cmd, flags, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestParseLanguageMapping(t *testing.T) {
	t.Parallel()

	t.Run("parses valid json", func(t *testing.T) {
		t.Parallel()

		raw := `[{"lang_iso":"en","custom_iso":"en_US"}]`

		got, err := parseLanguageMapping(raw)
		if err != nil {
			t.Fatalf("parseLanguageMapping() error = %v", err)
		}

		want := []map[string]any{
			{
				"lang_iso":   "en",
				"custom_iso": "en_US",
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected parsed mapping: got %#v, want %#v", got, want)
		}
	})

	t.Run("returns error for invalid json", func(t *testing.T) {
		t.Parallel()

		_, err := parseLanguageMapping(`{`)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func newTestCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}

	cmd.Flags().Bool("original-filenames", false, "")
	cmd.Flags().Bool("all-platforms", false, "")
	cmd.Flags().Bool("add-newline-eof", false, "")
	cmd.Flags().Bool("include-comments", false, "")
	cmd.Flags().Bool("include-description", false, "")
	cmd.Flags().Bool("replace-breaks", false, "")
	cmd.Flags().Bool("disable-references", false, "")
	cmd.Flags().Bool("icu-numeric", false, "")
	cmd.Flags().Bool("escape-percent", false, "")
	cmd.Flags().Bool("yaml-include-root", false, "")
	cmd.Flags().Bool("json-unescaped-slashes", false, "")
	cmd.Flags().Bool("compact", false, "")

	cmd.Flags().Int("filter-task-id", 0, "")

	return cmd
}

func newBoolDefaultTrueCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}

	cmd.Flags().Bool("original-filenames", false, "")
	cmd.Flags().Bool("all-platforms", false, "")
	cmd.Flags().Bool("add-newline-eof", false, "")
	cmd.Flags().Bool("include-comments", false, "")
	cmd.Flags().Bool("include-description", false, "")
	cmd.Flags().Bool("replace-breaks", false, "")
	cmd.Flags().Bool("disable-references", false, "")
	cmd.Flags().Bool("icu-numeric", false, "")
	cmd.Flags().Bool("escape-percent", false, "")
	cmd.Flags().Bool("yaml-include-root", false, "")
	cmd.Flags().Bool("json-unescaped-slashes", false, "")
	cmd.Flags().Bool("compact", true, "")

	cmd.Flags().Int("filter-task-id", 0, "")

	return cmd
}

func TestBuildParams_OmitsLocalOnlyFields(t *testing.T) {
	t.Parallel()

	cmd := newTestCommand()
	flags := &Flags{
		Out:    "./x",
		Async:  true,
		Format: "json",
	}

	params, err := buildParams(cmd, flags, nil)
	if err != nil {
		t.Fatalf("buildParams() error = %v", err)
	}

	for _, key := range []string{"out", "async"} {
		if _, ok := params[key]; ok {
			t.Fatalf("expected %q to be omitted, got %#v", key, params[key])
		}
	}
}

func TestBuildParams_OmitsWhitespaceOnlyStringFields(t *testing.T) {
	t.Parallel()

	cmd := newTestCommand()
	flags := &Flags{
		Format:          "json",
		BundleStructure: "   ",
		WebhookURL:      "\t",
	}

	params, err := buildParams(cmd, flags, nil)
	if err != nil {
		t.Fatalf("buildParams() error = %v", err)
	}

	if _, ok := params["bundle_structure"]; ok {
		t.Fatalf("expected bundle_structure to be omitted")
	}
	if _, ok := params["webhook_url"]; ok {
		t.Fatalf("expected webhook_url to be omitted")
	}
}

func TestBuildParams_UsesBoolDefaultFromConfig(t *testing.T) {
	t.Parallel()

	cmd := newTestCommand()
	flags := &Flags{
		Format: "json",
	}
	defaults := &DownloadConfig{
		Compact: new(true),
	}

	params, err := buildParams(cmd, flags, defaults)
	if err != nil {
		t.Fatalf("buildParams() error = %v", err)
	}

	got, ok := params["compact"]
	if !ok {
		t.Fatal("expected compact to be set from defaults")
	}
	if got != true {
		t.Fatalf("expected compact=true, got %#v", got)
	}
}

func TestBuildParams_UsesInt64DefaultFromConfig(t *testing.T) {
	t.Parallel()

	cmd := newTestCommand()
	flags := &Flags{
		Format: "json",
	}
	defaults := &DownloadConfig{
		FilterTaskID: new(int64(123)),
	}

	params, err := buildParams(cmd, flags, defaults)
	if err != nil {
		t.Fatalf("buildParams() error = %v", err)
	}

	got, ok := params["filter_task_id"]
	if !ok {
		t.Fatal("expected filter_task_id to be set from defaults")
	}
	if got != int64(123) {
		t.Fatalf("expected 123, got %#v", got)
	}
}

func TestBuildParams_ExplicitFlagOverridesDefault(t *testing.T) {
	t.Parallel()

	cmd := newTestCommand()
	if err := cmd.Flags().Set("compact", "false"); err != nil {
		t.Fatalf("set compact: %v", err)
	}

	flags := &Flags{
		Format:  "json",
		Compact: false,
	}
	defaults := &DownloadConfig{
		Compact: new(true),
	}

	params, err := buildParams(cmd, flags, defaults)
	if err != nil {
		t.Fatalf("buildParams() error = %v", err)
	}

	got, ok := params["compact"]
	if !ok {
		t.Fatal("expected compact to be set")
	}
	if got != false {
		t.Fatalf("expected compact=false, got %#v", got)
	}
}

func TestBuildParams_OmitsLanguageMappingWhenBlank(t *testing.T) {
	t.Parallel()

	cmd := newTestCommand()
	flags := &Flags{
		Format:              "json",
		LanguageMappingJSON: "   ",
	}

	params, err := buildParams(cmd, flags, nil)
	if err != nil {
		t.Fatalf("buildParams() error = %v", err)
	}

	if _, ok := params["language_mapping"]; ok {
		t.Fatalf("expected language_mapping to be omitted")
	}
}

func TestBuildParams_LanguageMappingErrorContainsFlagName(t *testing.T) {
	t.Parallel()

	cmd := newTestCommand()
	flags := &Flags{
		Format:              "json",
		LanguageMappingJSON: `not-json`,
	}

	_, err := buildParams(cmd, flags, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "parse --language-mapping") {
		t.Fatalf("unexpected error: %v", err)
	}
}
