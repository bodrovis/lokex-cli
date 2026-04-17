package download

import (
	"github.com/spf13/cobra"

	downloadCfg "github.com/bodrovis/lokex-cli/internal/download_config"
)

func applyDefaults(cmd *cobra.Command, flags *Flags, cfg *downloadCfg.DownloadConfig) {
	if cfg == nil {
		return
	}

	if !cmd.Flags().Changed("out") && cfg.Out != nil {
		flags.Out = *cfg.Out
	}

	if !cmd.Flags().Changed("format") && cfg.Format != nil {
		flags.Format = *cfg.Format
	}

	if !cmd.Flags().Changed("context-timeout") && cfg.ContextTimeout != nil {
		flags.ContextTimeout = *cfg.ContextTimeout
	}

	if !cmd.Flags().Changed("async") && cfg.Async != nil {
		flags.Async = *cfg.Async
	}

	if !cmd.Flags().Changed("original-filenames") && cfg.OriginalFilenames != nil {
		flags.OriginalFilenames = *cfg.OriginalFilenames
	}

	if !cmd.Flags().Changed("bundle-structure") && cfg.BundleStructure != nil {
		flags.BundleStructure = *cfg.BundleStructure
	}

	if !cmd.Flags().Changed("directory-prefix") && cfg.DirectoryPrefix != nil {
		flags.DirectoryPrefix = *cfg.DirectoryPrefix
	}

	if !cmd.Flags().Changed("all-platforms") && cfg.AllPlatforms != nil {
		flags.AllPlatforms = *cfg.AllPlatforms
	}

	if !cmd.Flags().Changed("filter-langs") && cfg.FilterLangs != nil {
		flags.FilterLangs = *cfg.FilterLangs
	}

	if !cmd.Flags().Changed("filter-data") && cfg.FilterData != nil {
		flags.FilterData = *cfg.FilterData
	}

	if !cmd.Flags().Changed("filter-filenames") && cfg.FilterFilenames != nil {
		flags.FilterFilenames = *cfg.FilterFilenames
	}

	if !cmd.Flags().Changed("custom-translation-status-ids") && cfg.CustomTranslationStatusIDs != nil {
		flags.CustomTranslationStatusIDs = *cfg.CustomTranslationStatusIDs
	}

	if !cmd.Flags().Changed("include-tags") && cfg.IncludeTags != nil {
		flags.IncludeTags = *cfg.IncludeTags
	}

	if !cmd.Flags().Changed("exclude-tags") && cfg.ExcludeTags != nil {
		flags.ExcludeTags = *cfg.ExcludeTags
	}

	if !cmd.Flags().Changed("include-pids") && cfg.IncludePIDs != nil {
		flags.IncludePIDs = *cfg.IncludePIDs
	}

	if !cmd.Flags().Changed("triggers") && cfg.Triggers != nil {
		flags.Triggers = *cfg.Triggers
	}

	if !cmd.Flags().Changed("filter-repositories") && cfg.FilterRepositories != nil {
		flags.FilterRepositories = *cfg.FilterRepositories
	}

	if !cmd.Flags().Changed("filter-task-id") && cfg.FilterTaskID != nil {
		flags.FilterTaskID = *cfg.FilterTaskID
	}

	if !cmd.Flags().Changed("add-newline-eof") && cfg.AddNewlineEOF != nil {
		flags.AddNewlineEOF = *cfg.AddNewlineEOF
	}

	if !cmd.Flags().Changed("include-comments") && cfg.IncludeComments != nil {
		flags.IncludeComments = *cfg.IncludeComments
	}

	if !cmd.Flags().Changed("include-description") && cfg.IncludeDescription != nil {
		flags.IncludeDescription = *cfg.IncludeDescription
	}

	if !cmd.Flags().Changed("replace-breaks") && cfg.ReplaceBreaks != nil {
		flags.ReplaceBreaks = *cfg.ReplaceBreaks
	}

	if !cmd.Flags().Changed("disable-references") && cfg.DisableReferences != nil {
		flags.DisableReferences = *cfg.DisableReferences
	}

	if !cmd.Flags().Changed("icu-numeric") && cfg.ICUNumeric != nil {
		flags.ICUNumeric = *cfg.ICUNumeric
	}

	if !cmd.Flags().Changed("escape-percent") && cfg.EscapePercent != nil {
		flags.EscapePercent = *cfg.EscapePercent
	}

	if !cmd.Flags().Changed("yaml-include-root") && cfg.YAMLIncludeRoot != nil {
		flags.YAMLIncludeRoot = *cfg.YAMLIncludeRoot
	}

	if !cmd.Flags().Changed("json-unescaped-slashes") && cfg.JSONUnescapedSlashes != nil {
		flags.JSONUnescapedSlashes = *cfg.JSONUnescapedSlashes
	}

	if !cmd.Flags().Changed("compact") && cfg.Compact != nil {
		flags.Compact = *cfg.Compact
	}

	if !cmd.Flags().Changed("export-sort") && cfg.ExportSort != nil {
		flags.ExportSort = *cfg.ExportSort
	}

	if !cmd.Flags().Changed("export-empty-as") && cfg.ExportEmptyAs != nil {
		flags.ExportEmptyAs = *cfg.ExportEmptyAs
	}

	if !cmd.Flags().Changed("export-null-as") && cfg.ExportNullAs != nil {
		flags.ExportNullAs = *cfg.ExportNullAs
	}

	if !cmd.Flags().Changed("plural-format") && cfg.PluralFormat != nil {
		flags.PluralFormat = *cfg.PluralFormat
	}

	if !cmd.Flags().Changed("placeholder-format") && cfg.PlaceholderFormat != nil {
		flags.PlaceholderFormat = *cfg.PlaceholderFormat
	}

	if !cmd.Flags().Changed("webhook-url") && cfg.WebhookURL != nil {
		flags.WebhookURL = *cfg.WebhookURL
	}

	if !cmd.Flags().Changed("indentation") && cfg.Indentation != nil {
		flags.Indentation = *cfg.Indentation
	}

	if !cmd.Flags().Changed("java-properties-encoding") && cfg.JavaPropertiesEncoding != nil {
		flags.JavaPropertiesEncoding = *cfg.JavaPropertiesEncoding
	}

	if !cmd.Flags().Changed("java-properties-separator") && cfg.JavaPropertiesSeparator != nil {
		flags.JavaPropertiesSeparator = *cfg.JavaPropertiesSeparator
	}

	if !cmd.Flags().Changed("bundle-description") && cfg.BundleDescription != nil {
		flags.BundleDescription = *cfg.BundleDescription
	}

	if !cmd.Flags().Changed("language-mapping") && cfg.LanguageMappingJSON != nil {
		flags.LanguageMappingJSON = *cfg.LanguageMappingJSON
	}
}
