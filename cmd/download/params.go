package download

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	lokexdownload "github.com/bodrovis/lokex/v2/client/download"

	cliutils "github.com/bodrovis/lokex-cli/internal/cli"
)

func buildParams(cmd *cobra.Command, flags *Flags) (lokexdownload.DownloadParams, error) {
	params := lokexdownload.DownloadParams{
		"format": flags.Format,
	}

	cliutils.SetString(params, "bundle_structure", flags.BundleStructure)
	cliutils.SetString(params, "directory_prefix", flags.DirectoryPrefix)
	cliutils.SetString(params, "export_sort", flags.ExportSort)
	cliutils.SetString(params, "export_empty_as", flags.ExportEmptyAs)
	cliutils.SetString(params, "export_null_as", flags.ExportNullAs)
	cliutils.SetString(params, "plural_format", flags.PluralFormat)
	cliutils.SetString(params, "placeholder_format", flags.PlaceholderFormat)
	cliutils.SetString(params, "webhook_url", flags.WebhookURL)
	cliutils.SetString(params, "indentation", flags.Indentation)
	cliutils.SetString(params, "java_properties_encoding", flags.JavaPropertiesEncoding)
	cliutils.SetString(params, "java_properties_separator", flags.JavaPropertiesSeparator)
	cliutils.SetString(params, "bundle_description", flags.BundleDescription)

	cliutils.SetStringSlice(params, "filter_langs", flags.FilterLangs)
	cliutils.SetStringSlice(params, "filter_data", flags.FilterData)
	cliutils.SetStringSlice(params, "filter_filenames", flags.FilterFilenames)
	cliutils.SetStringSlice(params, "custom_translation_status_ids", flags.CustomTranslationStatusIDs)
	cliutils.SetStringSlice(params, "include_tags", flags.IncludeTags)
	cliutils.SetStringSlice(params, "exclude_tags", flags.ExcludeTags)
	cliutils.SetStringSlice(params, "include_pids", flags.IncludePIDs)
	cliutils.SetStringSlice(params, "triggers", flags.Triggers)
	cliutils.SetStringSlice(params, "filter_repositories", flags.FilterRepositories)

	cliutils.SetChangedBool(cmd, params, "original-filenames", "original_filenames", flags.OriginalFilenames)
	cliutils.SetChangedBool(cmd, params, "all-platforms", "all_platforms", flags.AllPlatforms)
	cliutils.SetChangedBool(cmd, params, "add-newline-eof", "add_newline_eof", flags.AddNewlineEOF)
	cliutils.SetChangedBool(cmd, params, "include-comments", "include_comments", flags.IncludeComments)
	cliutils.SetChangedBool(cmd, params, "include-description", "include_description", flags.IncludeDescription)
	cliutils.SetChangedBool(cmd, params, "replace-breaks", "replace_breaks", flags.ReplaceBreaks)
	cliutils.SetChangedBool(cmd, params, "disable-references", "disable_references", flags.DisableReferences)
	cliutils.SetChangedBool(cmd, params, "icu-numeric", "icu_numeric", flags.ICUNumeric)
	cliutils.SetChangedBool(cmd, params, "escape-percent", "escape_percent", flags.EscapePercent)
	cliutils.SetChangedBool(cmd, params, "yaml-include-root", "yaml_include_root", flags.YAMLIncludeRoot)
	cliutils.SetChangedBool(cmd, params, "json-unescaped-slashes", "json_unescaped_slashes", flags.JSONUnescapedSlashes)
	cliutils.SetChangedBool(cmd, params, "compact", "compact", flags.Compact)

	if cmd.Flags().Changed("filter-task-id") {
		params["filter_task_id"] = flags.FilterTaskID
	}

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
