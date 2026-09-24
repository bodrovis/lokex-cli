package download

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"strings"

	lokexdownload "github.com/bodrovis/lokex/v2/client/download"
)

func buildParams(cfg *DownloadConfig) (lokexdownload.DownloadParams, error) {
	if cfg == nil {
		return nil, errors.New("download config is nil")
	}

	req := lokexdownload.DownloadParams{}

	putParam(req, "format", cfg.Format)
	putParam(req, "original_filenames", cfg.OriginalFilenames)
	putParam(req, "bundle_structure", cfg.BundleStructure)
	putParam(req, "directory_prefix", cfg.DirectoryPrefix)
	putParam(req, "all_platforms", cfg.AllPlatforms)
	putParam(req, "filter_langs", cfg.FilterLangs)
	putParam(req, "filter_data", cfg.FilterData)
	putParam(req, "filter_filenames", cfg.FilterFilenames)
	putParam(req, "add_newline_eof", cfg.AddNewlineEOF)
	putParam(req, "custom_translation_status_ids", cfg.CustomTranslationStatusIDs)
	putParam(req, "include_tags", cfg.IncludeTags)
	putParam(req, "exclude_tags", cfg.ExcludeTags)
	putParam(req, "export_sort", cfg.ExportSort)
	putParam(req, "export_empty_as", cfg.ExportEmptyAs)
	putParam(req, "export_null_as", cfg.ExportNullAs)
	putParam(req, "include_comments", cfg.IncludeComments)
	putParam(req, "include_description", cfg.IncludeDescription)
	putParam(req, "include_pids", cfg.IncludePIDs)
	putParam(req, "triggers", cfg.Triggers)
	putParam(req, "filter_repositories", cfg.FilterRepositories)
	putParam(req, "replace_breaks", cfg.ReplaceBreaks)
	putParam(req, "disable_references", cfg.DisableReferences)
	putParam(req, "plural_format", cfg.PluralFormat)
	putParam(req, "placeholder_format", cfg.PlaceholderFormat)
	putParam(req, "webhook_url", cfg.WebhookURL)
	putParam(req, "icu_numeric", cfg.ICUNumeric)
	putParam(req, "escape_percent", cfg.EscapePercent)
	putParam(req, "indentation", cfg.Indentation)
	putParam(req, "yaml_include_root", cfg.YAMLIncludeRoot)
	putParam(req, "json_unescaped_slashes", cfg.JSONUnescapedSlashes)
	putParam(req, "java_properties_encoding", cfg.JavaPropertiesEncoding)
	putParam(req, "java_properties_separator", cfg.JavaPropertiesSeparator)
	putParam(req, "bundle_description", cfg.BundleDescription)
	putParam(req, "filter_task_id", cfg.FilterTaskID)
	putParam(req, "compact", cfg.Compact)

	if err := applyLanguageMapping(req, cfg.LanguageMappingJSON); err != nil {
		return nil, err
	}

	return req, nil
}

func putParam[T any](
	req lokexdownload.DownloadParams,
	key string,
	value *T,
) {
	if value != nil {
		req[key] = *value
	}
}

func applyLanguageMapping(
	req lokexdownload.DownloadParams,
	raw *string,
) error {
	if raw == nil {
		return nil
	}

	value := strings.TrimSpace(*raw)
	if value == "" {
		return nil
	}

	languageMapping, err := parseLanguageMapping(value)
	if err != nil {
		return fmt.Errorf("parse language-mapping: %w", err)
	}

	req["language_mapping"] = languageMapping

	return nil
}

type languageMapping struct {
	OriginalLanguageISO string `json:"original_language_iso"`
	CustomLanguageISO   string `json:"custom_language_iso"`
}

func parseLanguageMapping(raw string) ([]languageMapping, error) {
	var out []languageMapping

	if err := json.Unmarshal(
		[]byte(raw),
		&out,
		json.RejectUnknownMembers(true),
	); err != nil {
		return nil, err
	}

	return out, nil
}
