package upload

import (
	"errors"
	"strings"

	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
)

func buildParams(cfg *UploadConfig) (lokexupload.UploadParams, error) {
	if cfg == nil {
		return nil, errors.New("upload config is nil")
	}

	req := lokexupload.UploadParams{}

	if cfg.Filename != nil {
		req["filename"] = normalizeFilename(*cfg.Filename)
	}

	if cfg.Data != nil && *cfg.Data != "" {
		req["data"] = *cfg.Data
	}

	if cfg.LangISO != nil {
		langISO := strings.TrimSpace(*cfg.LangISO)

		if langISO != "" {
			req["lang_iso"] = langISO
		}
	}

	if cfg.Format != nil && *cfg.Format != "" {
		req["format"] = *cfg.Format
	}

	if cfg.Tags != nil && len(*cfg.Tags) > 0 {
		req["tags"] = *cfg.Tags
	}

	if cfg.CustomTranslationStatusIDs != nil &&
		len(*cfg.CustomTranslationStatusIDs) > 0 {
		req["custom_translation_status_ids"] = *cfg.CustomTranslationStatusIDs
	}

	putBool(req, "convert_placeholders", cfg.ConvertPlaceholders)
	putBool(req, "detect_icu_plurals", cfg.DetectICUPlurals)
	putBool(req, "tag_inserted_keys", cfg.TagInsertedKeys)
	putBool(req, "tag_updated_keys", cfg.TagUpdatedKeys)
	putBool(req, "tag_skipped_keys", cfg.TagSkippedKeys)
	putBool(req, "replace_modified", cfg.ReplaceModified)
	putBool(req, "slashn_to_linebreak", cfg.SlashNToLinebreak)
	putBool(req, "keys_to_values", cfg.KeysToValues)
	putBool(req, "distinguish_by_file", cfg.DistinguishByFile)
	putBool(req, "apply_tm", cfg.ApplyTM)
	putBool(req, "use_automations", cfg.UseAutomations)
	putBool(req, "hidden_from_contributors", cfg.HiddenFromContributors)
	putBool(req, "cleanup_mode", cfg.CleanupMode)
	putBool(
		req,
		"custom_translation_status_inserted_keys",
		cfg.CustomTranslationStatusInsertedKeys,
	)
	putBool(
		req,
		"custom_translation_status_updated_keys",
		cfg.CustomTranslationStatusUpdatedKeys,
	)
	putBool(
		req,
		"custom_translation_status_skipped_keys",
		cfg.CustomTranslationStatusSkippedKeys,
	)
	putBool(req, "skip_detect_lang_iso", cfg.SkipDetectLangISO)

	if cfg.FilterTaskID != nil {
		req["filter_task_id"] = *cfg.FilterTaskID
	}

	return req, nil
}

func putBool(
	req lokexupload.UploadParams,
	key string,
	value *bool,
) {
	if value != nil {
		req[key] = *value
	}
}

func normalizeFilename(filename string) string {
	return strings.ReplaceAll(strings.TrimSpace(filename), `\`, `/`)
}
