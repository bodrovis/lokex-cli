package upload

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	lokexupload "github.com/bodrovis/lokex/v2/client/upload"

	setters "github.com/bodrovis/lokex-cli/internal/params"
	uploadCfg "github.com/bodrovis/lokex-cli/internal/upload_config"
)

func buildParams(cmd *cobra.Command, flags *Flags, defaults *uploadCfg.UploadConfig) lokexupload.UploadParams {
	if defaults == nil {
		defaults = &uploadCfg.UploadConfig{}
	}

	params := lokexupload.UploadParams{
		"filename": filepath.ToSlash(strings.TrimSpace(flags.Filename)),
		"lang_iso": strings.TrimSpace(flags.LangISO),
	}

	setters.SetString(params, "data", flags.Data)
	setters.SetString(params, "format", flags.Format)

	setters.SetStringSlice(params, "tags", flags.Tags)
	setters.SetStringSlice(params, "custom_translation_status_ids", flags.CustomTranslationStatusIDs)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"convert-placeholders",
		"convert_placeholders",
		flags.ConvertPlaceholders,
		defaults.ConvertPlaceholders,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"detect-icu-plurals",
		"detect_icu_plurals",
		flags.DetectICUPlurals,
		defaults.DetectICUPlurals,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"tag-inserted-keys",
		"tag_inserted_keys",
		flags.TagInsertedKeys,
		defaults.TagInsertedKeys,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"tag-updated-keys",
		"tag_updated_keys",
		flags.TagUpdatedKeys,
		defaults.TagUpdatedKeys,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"tag-skipped-keys",
		"tag_skipped_keys",
		flags.TagSkippedKeys,
		defaults.TagSkippedKeys,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"replace-modified",
		"replace_modified",
		flags.ReplaceModified,
		defaults.ReplaceModified,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"slashn-to-linebreak",
		"slashn_to_linebreak",
		flags.SlashNToLinebreak,
		defaults.SlashNToLinebreak,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"keys-to-values",
		"keys_to_values",
		flags.KeysToValues,
		defaults.KeysToValues,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"distinguish-by-file",
		"distinguish_by_file",
		flags.DistinguishByFile,
		defaults.DistinguishByFile,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"apply-tm",
		"apply_tm",
		flags.ApplyTM,
		defaults.ApplyTM,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"use-automations",
		"use_automations",
		flags.UseAutomations,
		defaults.UseAutomations,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"hidden-from-contributors",
		"hidden_from_contributors",
		flags.HiddenFromContributors,
		defaults.HiddenFromContributors,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"cleanup-mode",
		"cleanup_mode",
		flags.CleanupMode,
		defaults.CleanupMode,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"custom-translation-status-inserted-keys",
		"custom_translation_status_inserted_keys",
		flags.CustomTranslationStatusInsertedKeys,
		defaults.CustomTranslationStatusInsertedKeys,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"custom-translation-status-updated-keys",
		"custom_translation_status_updated_keys",
		flags.CustomTranslationStatusUpdatedKeys,
		defaults.CustomTranslationStatusUpdatedKeys,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"custom-translation-status-skipped-keys",
		"custom_translation_status_skipped_keys",
		flags.CustomTranslationStatusSkippedKeys,
		defaults.CustomTranslationStatusSkippedKeys,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"skip-detect-lang-iso",
		"skip_detect_lang_iso",
		flags.SkipDetectLangISO,
		defaults.SkipDetectLangISO,
	)

	setters.SetInt64WithDefault(
		cmd,
		params,
		"filter-task-id",
		"filter_task_id",
		flags.FilterTaskID,
		defaults.FilterTaskID,
	)

	return params
}
