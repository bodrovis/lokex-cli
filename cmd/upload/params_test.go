package upload

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestBuildParams_Minimal(t *testing.T) {
	t.Parallel()

	cmd := newTestCommand()

	flags := &Flags{
		Filename: "en.json",
		LangISO:  "en",
	}

	params := buildParams(cmd, flags, nil)

	if got := params["filename"]; got != "en.json" {
		t.Fatalf("unexpected filename: got %v, want %q", got, "en.json")
	}

	if got := params["lang_iso"]; got != "en" {
		t.Fatalf("unexpected lang_iso: got %v, want %q", got, "en")
	}

	if len(params) != 2 {
		t.Fatalf("expected exactly 2 params, got %d: %#v", len(params), params)
	}
}

func TestBuildParams_NormalizesFilenamePath(t *testing.T) {
	t.Parallel()

	cmd := newTestCommand()

	flags := &Flags{
		Filename: `admin\main.json`,
		LangISO:  "en",
	}

	params := buildParams(cmd, flags, nil)

	if got := params["filename"]; got != "admin/main.json" {
		t.Fatalf("unexpected filename: got %v, want %q", got, "admin/main.json")
	}
}

func TestBuildParams_SetsTrimmedStringFields(t *testing.T) {
	t.Parallel()

	cmd := newTestCommand()

	flags := &Flags{
		Filename: "en.json",
		LangISO:  "en",
		Data:     " ZGF0YQ== ",
		Format:   " json ",
	}

	params := buildParams(cmd, flags, nil)

	if got := params["data"]; got != "ZGF0YQ==" {
		t.Fatalf("unexpected data: got %v, want %q", got, "ZGF0YQ==")
	}

	if got := params["format"]; got != "json" {
		t.Fatalf("unexpected format: got %v, want %q", got, "json")
	}
}

func TestBuildParams_SetsFilteredStringSliceFields(t *testing.T) {
	t.Parallel()

	cmd := newTestCommand()

	flags := &Flags{
		Filename:                   "en.json",
		LangISO:                    "en",
		Tags:                       []string{"mobile", "", "   ", "web"},
		CustomTranslationStatusIDs: []string{"1", "", "2"},
	}

	params := buildParams(cmd, flags, nil)

	gotTags, ok := params["tags"]
	if !ok {
		t.Fatal("expected tags to be set")
	}

	wantTags := []string{"mobile", "web"}
	if !reflect.DeepEqual(gotTags, wantTags) {
		t.Fatalf("unexpected tags: got %#v, want %#v", gotTags, wantTags)
	}

	gotStatuses, ok := params["custom_translation_status_ids"]
	if !ok {
		t.Fatal("expected custom_translation_status_ids to be set")
	}

	wantStatuses := []string{"1", "2"}
	if !reflect.DeepEqual(gotStatuses, wantStatuses) {
		t.Fatalf("unexpected custom_translation_status_ids: got %#v, want %#v", gotStatuses, wantStatuses)
	}
}

func TestBuildParams_OmitsEmptyOptionalFields(t *testing.T) {
	t.Parallel()

	cmd := newTestCommand()

	flags := &Flags{
		Filename:                   "en.json",
		LangISO:                    "en",
		Data:                       "   ",
		Format:                     "",
		Tags:                       []string{"", "   "},
		CustomTranslationStatusIDs: nil,
	}

	params := buildParams(cmd, flags, nil)

	if _, ok := params["data"]; ok {
		t.Fatalf("expected data to be omitted, got %#v", params["data"])
	}
	if _, ok := params["format"]; ok {
		t.Fatalf("expected format to be omitted, got %#v", params["format"])
	}
	if _, ok := params["tags"]; ok {
		t.Fatalf("expected tags to be omitted, got %#v", params["tags"])
	}
	if _, ok := params["custom_translation_status_ids"]; ok {
		t.Fatalf("expected custom_translation_status_ids to be omitted, got %#v", params["custom_translation_status_ids"])
	}
}

func TestBuildParams_SetsBoolOnlyWhenFlagChanged(t *testing.T) {
	t.Parallel()

	t.Run("omits bool when flag not changed", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()

		flags := &Flags{
			Filename:            "en.json",
			LangISO:             "en",
			ConvertPlaceholders: true,
		}

		params := buildParams(cmd, flags, nil)

		if _, ok := params["convert_placeholders"]; ok {
			t.Fatalf("expected convert_placeholders to be omitted, got %#v", params["convert_placeholders"])
		}
	})

	t.Run("sets bool true when explicitly enabled", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()
		if err := cmd.Flags().Parse([]string{"--convert-placeholders"}); err != nil {
			t.Fatalf("parse flags: %v", err)
		}

		flags := &Flags{
			Filename:            "en.json",
			LangISO:             "en",
			ConvertPlaceholders: true,
		}

		params := buildParams(cmd, flags, nil)

		got, ok := params["convert_placeholders"]
		if !ok {
			t.Fatal("expected convert_placeholders to be set")
		}
		if got != true {
			t.Fatalf("unexpected convert_placeholders: got %#v, want true", got)
		}
	})

	t.Run("sets bool false when explicitly disabled", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()
		if err := cmd.Flags().Parse([]string{"--convert-placeholders=false"}); err != nil {
			t.Fatalf("parse flags: %v", err)
		}

		flags := &Flags{
			Filename:            "en.json",
			LangISO:             "en",
			ConvertPlaceholders: false,
		}

		params := buildParams(cmd, flags, nil)

		got, ok := params["convert_placeholders"]
		if !ok {
			t.Fatal("expected convert_placeholders to be set")
		}
		if got != false {
			t.Fatalf("unexpected convert_placeholders: got %#v, want false", got)
		}
	})
}

func TestBuildParams_FilterTaskID(t *testing.T) {
	t.Parallel()

	t.Run("omits filter_task_id when flag not changed", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()

		flags := &Flags{
			Filename:     "en.json",
			LangISO:      "en",
			FilterTaskID: 42,
		}

		params := buildParams(cmd, flags, nil)

		if _, ok := params["filter_task_id"]; ok {
			t.Fatalf("expected filter_task_id to be omitted, got %#v", params["filter_task_id"])
		}
	})

	t.Run("sets filter_task_id when flag changed", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()
		if err := cmd.Flags().Parse([]string{"--filter-task-id=42"}); err != nil {
			t.Fatalf("parse flags: %v", err)
		}

		flags := &Flags{
			Filename:     "en.json",
			LangISO:      "en",
			FilterTaskID: 42,
		}

		params := buildParams(cmd, flags, nil)

		got, ok := params["filter_task_id"]
		if !ok {
			t.Fatal("expected filter_task_id to be set")
		}
		if got != int64(42) {
			t.Fatalf("unexpected filter_task_id: got %#v, want %d", got, int64(42))
		}
	})
}

func TestBuildParams_SetsMultipleBoolFields(t *testing.T) {
	t.Parallel()

	cmd := newTestCommand()
	if err := cmd.Flags().Parse([]string{
		"--apply-tm",
		"--use-automations",
		"--hidden-from-contributors=false",
	}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	flags := &Flags{
		Filename:               "en.json",
		LangISO:                "en",
		ApplyTM:                true,
		UseAutomations:         true,
		HiddenFromContributors: false,
	}

	params := buildParams(cmd, flags, nil)

	if got := params["apply_tm"]; got != true {
		t.Fatalf("unexpected apply_tm: got %#v, want true", got)
	}
	if got := params["use_automations"]; got != true {
		t.Fatalf("unexpected use_automations: got %#v, want true", got)
	}
	if got := params["hidden_from_contributors"]; got != false {
		t.Fatalf("unexpected hidden_from_contributors: got %#v, want false", got)
	}
}

func newTestCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}

	cmd.Flags().Bool("convert-placeholders", false, "")
	cmd.Flags().Bool("detect-icu-plurals", false, "")
	cmd.Flags().Bool("tag-inserted-keys", false, "")
	cmd.Flags().Bool("tag-updated-keys", false, "")
	cmd.Flags().Bool("tag-skipped-keys", false, "")
	cmd.Flags().Bool("replace-modified", false, "")
	cmd.Flags().Bool("slashn-to-linebreak", false, "")
	cmd.Flags().Bool("keys-to-values", false, "")
	cmd.Flags().Bool("distinguish-by-file", false, "")
	cmd.Flags().Bool("apply-tm", false, "")
	cmd.Flags().Bool("use-automations", false, "")
	cmd.Flags().Bool("hidden-from-contributors", true, "")
	cmd.Flags().Bool("cleanup-mode", false, "")
	cmd.Flags().Bool("custom-translation-status-inserted-keys", false, "")
	cmd.Flags().Bool("custom-translation-status-updated-keys", false, "")
	cmd.Flags().Bool("custom-translation-status-skipped-keys", false, "")
	cmd.Flags().Bool("skip-detect-lang-iso", false, "")

	cmd.Flags().Int64("filter-task-id", 0, "")

	return cmd
}
