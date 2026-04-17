package download_config

import (
	"fmt"
	"time"

	vh "github.com/bodrovis/lokex-cli/internal/viper_helpers"
	"github.com/spf13/viper"
)

func LoadDownloadConfig(cfg *DownloadConfig, configFile, envPrefix string) error {
	v := vh.NewConfigViper(configFile, envPrefix)

	if err := bindDownloadEnv(v); err != nil {
		return fmt.Errorf("bind download env: %w", err)
	}

	if err := vh.ReadOptionalConfig(v, configFile); err != nil {
		return fmt.Errorf("read download config: %w", err)
	}

	vh.ApplyConfigValue(v, "download.out", func(val string) { cfg.Out = &val })
	vh.ApplyConfigValue(v, "download.format", func(val string) { cfg.Format = &val })
	vh.ApplyConfigValue(v, "download.context-timeout", func(val time.Duration) { cfg.ContextTimeout = &val })
	vh.ApplyConfigValue(v, "download.async", func(val bool) { cfg.Async = &val })

	vh.ApplyConfigValue(v, "download.original-filenames", func(val bool) { cfg.OriginalFilenames = &val })
	vh.ApplyConfigValue(v, "download.bundle-structure", func(val string) { cfg.BundleStructure = &val })
	vh.ApplyConfigValue(v, "download.directory-prefix", func(val string) { cfg.DirectoryPrefix = &val })
	vh.ApplyConfigValue(v, "download.all-platforms", func(val bool) { cfg.AllPlatforms = &val })
	vh.ApplyConfigStringSlice(v, "download.filter-langs", func(val []string) { cfg.FilterLangs = &val })
	vh.ApplyConfigStringSlice(v, "download.filter-data", func(val []string) { cfg.FilterData = &val })
	vh.ApplyConfigStringSlice(v, "download.filter-filenames", func(val []string) { cfg.FilterFilenames = &val })
	vh.ApplyConfigValue(v, "download.add-newline-eof", func(val bool) { cfg.AddNewlineEOF = &val })
	vh.ApplyConfigStringSlice(v, "download.custom-translation-status-ids", func(val []string) {
		cfg.CustomTranslationStatusIDs = &val
	})
	vh.ApplyConfigStringSlice(v, "download.include-tags", func(val []string) { cfg.IncludeTags = &val })
	vh.ApplyConfigStringSlice(v, "download.exclude-tags", func(val []string) { cfg.ExcludeTags = &val })
	vh.ApplyConfigValue(v, "download.export-sort", func(val string) { cfg.ExportSort = &val })
	vh.ApplyConfigValue(v, "download.export-empty-as", func(val string) { cfg.ExportEmptyAs = &val })
	vh.ApplyConfigValue(v, "download.export-null-as", func(val string) { cfg.ExportNullAs = &val })
	vh.ApplyConfigValue(v, "download.include-comments", func(val bool) { cfg.IncludeComments = &val })
	vh.ApplyConfigValue(v, "download.include-description", func(val bool) { cfg.IncludeDescription = &val })
	vh.ApplyConfigStringSlice(v, "download.include-pids", func(val []string) { cfg.IncludePIDs = &val })
	vh.ApplyConfigStringSlice(v, "download.triggers", func(val []string) { cfg.Triggers = &val })
	vh.ApplyConfigStringSlice(v, "download.filter-repositories", func(val []string) { cfg.FilterRepositories = &val })
	vh.ApplyConfigValue(v, "download.replace-breaks", func(val bool) { cfg.ReplaceBreaks = &val })
	vh.ApplyConfigValue(v, "download.disable-references", func(val bool) { cfg.DisableReferences = &val })
	vh.ApplyConfigValue(v, "download.plural-format", func(val string) { cfg.PluralFormat = &val })
	vh.ApplyConfigValue(v, "download.placeholder-format", func(val string) { cfg.PlaceholderFormat = &val })
	vh.ApplyConfigValue(v, "download.webhook-url", func(val string) { cfg.WebhookURL = &val })
	vh.ApplyConfigValue(v, "download.language-mapping", func(val string) { cfg.LanguageMappingJSON = &val })
	vh.ApplyConfigValue(v, "download.icu-numeric", func(val bool) { cfg.ICUNumeric = &val })
	vh.ApplyConfigValue(v, "download.escape-percent", func(val bool) { cfg.EscapePercent = &val })
	vh.ApplyConfigValue(v, "download.indentation", func(val string) { cfg.Indentation = &val })
	vh.ApplyConfigValue(v, "download.yaml-include-root", func(val bool) { cfg.YAMLIncludeRoot = &val })
	vh.ApplyConfigValue(v, "download.json-unescaped-slashes", func(val bool) { cfg.JSONUnescapedSlashes = &val })
	vh.ApplyConfigValue(v, "download.java-properties-encoding", func(val string) { cfg.JavaPropertiesEncoding = &val })
	vh.ApplyConfigValue(v, "download.java-properties-separator", func(val string) { cfg.JavaPropertiesSeparator = &val })
	vh.ApplyConfigValue(v, "download.bundle-description", func(val string) { cfg.BundleDescription = &val })
	vh.ApplyConfigValue(v, "download.filter-task-id", func(val int64) { cfg.FilterTaskID = &val })
	vh.ApplyConfigValue(v, "download.compact", func(val bool) { cfg.Compact = &val })

	return nil
}

func bindDownloadEnv(v *viper.Viper) error {
	keys := []string{
		"download.out",
		"download.format",
		"download.context-timeout",
		"download.async",

		"download.original-filenames",
		"download.bundle-structure",
		"download.directory-prefix",
		"download.all-platforms",
		"download.filter-langs",
		"download.filter-data",
		"download.filter-filenames",
		"download.add-newline-eof",
		"download.custom-translation-status-ids",
		"download.include-tags",
		"download.exclude-tags",
		"download.export-sort",
		"download.export-empty-as",
		"download.export-null-as",
		"download.include-comments",
		"download.include-description",
		"download.include-pids",
		"download.triggers",
		"download.filter-repositories",
		"download.replace-breaks",
		"download.disable-references",
		"download.plural-format",
		"download.placeholder-format",
		"download.webhook-url",
		"download.language-mapping",
		"download.icu-numeric",
		"download.escape-percent",
		"download.indentation",
		"download.yaml-include-root",
		"download.json-unescaped-slashes",
		"download.java-properties-encoding",
		"download.java-properties-separator",
		"download.bundle-description",
		"download.filter-task-id",
		"download.compact",
	}

	return vh.BindEnvKeys(v, keys)
}
