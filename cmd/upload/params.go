package upload

import (
	"path/filepath"

	"github.com/spf13/cobra"

	lokexupload "github.com/bodrovis/lokex/v2/client/upload"

	cliutils "github.com/bodrovis/lokex-cli/internal/cli"
)

func buildParams(cmd *cobra.Command, flags *Flags) lokexupload.UploadParams {
	params := lokexupload.UploadParams{
		"filename": filepath.ToSlash(flags.Filename),
		"lang_iso": flags.LangISO,
	}

	cliutils.SetString(params, "data", flags.Data)
	cliutils.SetString(params, "format", flags.Format)

	cliutils.SetStringSlice(params, "tags", flags.Tags)
	cliutils.SetStringSlice(params, "custom_translation_status_ids", flags.CustomTranslationStatusIDs)

	cliutils.SetChangedBool(cmd, params, "convert-placeholders", "convert_placeholders", flags.ConvertPlaceholders)
	cliutils.SetChangedBool(cmd, params, "detect-icu-plurals", "detect_icu_plurals", flags.DetectICUPlurals)
	cliutils.SetChangedBool(cmd, params, "tag-inserted-keys", "tag_inserted_keys", flags.TagInsertedKeys)
	cliutils.SetChangedBool(cmd, params, "tag-updated-keys", "tag_updated_keys", flags.TagUpdatedKeys)
	cliutils.SetChangedBool(cmd, params, "tag-skipped-keys", "tag_skipped_keys", flags.TagSkippedKeys)
	cliutils.SetChangedBool(cmd, params, "replace-modified", "replace_modified", flags.ReplaceModified)
	cliutils.SetChangedBool(cmd, params, "slashn-to-linebreak", "slashn_to_linebreak", flags.SlashNToLinebreak)
	cliutils.SetChangedBool(cmd, params, "keys-to-values", "keys_to_values", flags.KeysToValues)
	cliutils.SetChangedBool(cmd, params, "distinguish-by-file", "distinguish_by_file", flags.DistinguishByFile)
	cliutils.SetChangedBool(cmd, params, "apply-tm", "apply_tm", flags.ApplyTM)
	cliutils.SetChangedBool(cmd, params, "use-automations", "use_automations", flags.UseAutomations)
	cliutils.SetChangedBool(cmd, params, "hidden-from-contributors", "hidden_from_contributors", flags.HiddenFromContributors)
	cliutils.SetChangedBool(cmd, params, "cleanup-mode", "cleanup_mode", flags.CleanupMode)
	cliutils.SetChangedBool(cmd, params, "custom-translation-status-inserted-keys", "custom_translation_status_inserted_keys", flags.CustomTranslationStatusInsertedKeys)
	cliutils.SetChangedBool(cmd, params, "custom-translation-status-updated-keys", "custom_translation_status_updated_keys", flags.CustomTranslationStatusUpdatedKeys)
	cliutils.SetChangedBool(cmd, params, "custom-translation-status-skipped-keys", "custom_translation_status_skipped_keys", flags.CustomTranslationStatusSkippedKeys)
	cliutils.SetChangedBool(cmd, params, "skip-detect-lang-iso", "skip_detect_lang_iso", flags.SkipDetectLangISO)

	if cmd.Flags().Changed("filter-task-id") {
		params["filter_task_id"] = flags.FilterTaskID
	}

	return params
}
