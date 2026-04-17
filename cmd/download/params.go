package download

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	lokexdownload "github.com/bodrovis/lokex/v2/client/download"

	downloadCfg "github.com/bodrovis/lokex-cli/internal/download_config"
	setters "github.com/bodrovis/lokex-cli/internal/params"
)

func buildParams(cmd *cobra.Command, flags *Flags, defaults *downloadCfg.DownloadConfig) (lokexdownload.DownloadParams, error) {
	if defaults == nil {
		defaults = &downloadCfg.DownloadConfig{}
	}

	params := lokexdownload.DownloadParams{
		"format": flags.Format,
	}

	setters.SetString(params, "bundle_structure", flags.BundleStructure)
	setters.SetString(params, "directory_prefix", flags.DirectoryPrefix)
	setters.SetString(params, "export_sort", flags.ExportSort)
	setters.SetString(params, "export_empty_as", flags.ExportEmptyAs)
	setters.SetString(params, "export_null_as", flags.ExportNullAs)
	setters.SetString(params, "plural_format", flags.PluralFormat)
	setters.SetString(params, "placeholder_format", flags.PlaceholderFormat)
	setters.SetString(params, "webhook_url", flags.WebhookURL)
	setters.SetString(params, "indentation", flags.Indentation)
	setters.SetString(params, "java_properties_encoding", flags.JavaPropertiesEncoding)
	setters.SetString(params, "java_properties_separator", flags.JavaPropertiesSeparator)
	setters.SetString(params, "bundle_description", flags.BundleDescription)

	setters.SetStringSlice(params, "filter_langs", flags.FilterLangs)
	setters.SetStringSlice(params, "filter_data", flags.FilterData)
	setters.SetStringSlice(params, "filter_filenames", flags.FilterFilenames)
	setters.SetStringSlice(params, "custom_translation_status_ids", flags.CustomTranslationStatusIDs)
	setters.SetStringSlice(params, "include_tags", flags.IncludeTags)
	setters.SetStringSlice(params, "exclude_tags", flags.ExcludeTags)
	setters.SetStringSlice(params, "include_pids", flags.IncludePIDs)
	setters.SetStringSlice(params, "triggers", flags.Triggers)
	setters.SetStringSlice(params, "filter_repositories", flags.FilterRepositories)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"original-filenames",
		"original_filenames",
		flags.OriginalFilenames,
		defaults.OriginalFilenames,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"all-platforms",
		"all_platforms",
		flags.AllPlatforms,
		defaults.AllPlatforms,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"add-newline-eof",
		"add_newline_eof",
		flags.AddNewlineEOF,
		defaults.AddNewlineEOF,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"include-comments",
		"include_comments",
		flags.IncludeComments,
		defaults.IncludeComments,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"include-description",
		"include_description",
		flags.IncludeDescription,
		defaults.IncludeDescription,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"replace-breaks",
		"replace_breaks",
		flags.ReplaceBreaks,
		defaults.ReplaceBreaks,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"disable-references",
		"disable_references",
		flags.DisableReferences,
		defaults.DisableReferences,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"icu-numeric",
		"icu_numeric",
		flags.ICUNumeric,
		defaults.ICUNumeric,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"escape-percent",
		"escape_percent",
		flags.EscapePercent,
		defaults.EscapePercent,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"yaml-include-root",
		"yaml_include_root",
		flags.YAMLIncludeRoot,
		defaults.YAMLIncludeRoot,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"json-unescaped-slashes",
		"json_unescaped_slashes",
		flags.JSONUnescapedSlashes,
		defaults.JSONUnescapedSlashes,
	)

	setters.SetBoolWithDefault(
		cmd,
		params,
		"compact",
		"compact",
		flags.Compact,
		defaults.Compact,
	)

	setters.SetInt64WithDefault(
		cmd,
		params,
		"filter-task-id",
		"filter_task_id",
		flags.FilterTaskID,
		defaults.FilterTaskID,
	)

	if strings.TrimSpace(flags.LanguageMappingJSON) != "" {
		languageMapping, err := parseLanguageMapping(flags.LanguageMappingJSON)
		if err != nil {
			return nil, fmt.Errorf("parse --language-mapping: %w", err)
		}

		params["language_mapping"] = languageMapping
	}

	return params, nil
}

func parseLanguageMapping(raw string) ([]map[string]any, error) {
	var out []map[string]any

	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil, err
	}

	return out, nil
}
