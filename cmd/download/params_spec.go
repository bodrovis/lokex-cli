package download

import (
	params "github.com/bodrovis/lokex-cli/internal/params"
	vh "github.com/bodrovis/lokex-cli/internal/viper_helpers"
	lokexdownload "github.com/bodrovis/lokex/v2/client/download"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DownloadParamSpec = params.ParamSpec[Flags, DownloadConfig, lokexdownload.DownloadParams]

var downloadParamSpecs = []DownloadParamSpec{
	{
		FlagName:  "out",
		ConfigKey: "download.out",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.Out,
				"out",
				flags.Out,
				"Directory to unzip downloaded bundle into",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("out") && cfg.Out != nil {
				flags.Out = *cfg.Out
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.out", func(val string) {
				cfg.Out = &val
			})
		},
	},
	{
		FlagName:  "format",
		ConfigKey: "download.format",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.Format,
				"format",
				flags.Format,
				"File format (e.g. json, strings, xml)",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("format") && cfg.Format != nil {
				flags.Format = *cfg.Format
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.format", func(val string) {
				cfg.Format = &val
			})
		},
		ApplyToRequest: reqDirectString("format", func(f *Flags) string {
			return f.Format
		}),
	},
	{
		FlagName:  "async",
		ConfigKey: "download.async",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.Async,
				"async",
				flags.Async,
				"Use Lokalise async download flow",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("async") && cfg.Async != nil {
				flags.Async = *cfg.Async
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.async", func(val bool) {
				cfg.Async = &val
			})
		},
	},
	{
		FlagName:  "original-filenames",
		ConfigKey: "download.original-filenames",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.OriginalFilenames,
				"original-filenames",
				flags.OriginalFilenames,
				"Use original filenames/formats",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("original-filenames") && cfg.OriginalFilenames != nil {
				flags.OriginalFilenames = *cfg.OriginalFilenames
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.original-filenames", func(val bool) {
				cfg.OriginalFilenames = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"original-filenames",
			"original_filenames",
			func(f *Flags) bool { return f.OriginalFilenames },
			func(c *DownloadConfig) *bool { return c.OriginalFilenames },
		),
	},
	{
		FlagName:  "bundle-structure",
		ConfigKey: "download.bundle-structure",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.BundleStructure,
				"bundle-structure",
				flags.BundleStructure,
				"Bundle structure when original-filenames=false",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("bundle-structure") && cfg.BundleStructure != nil {
				flags.BundleStructure = *cfg.BundleStructure
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.bundle-structure", func(val string) {
				cfg.BundleStructure = &val
			})
		},
		ApplyToRequest: reqString("bundle_structure", func(f *Flags) string {
			return f.BundleStructure
		}),
	},
	{
		FlagName:  "directory-prefix",
		ConfigKey: "download.directory-prefix",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.DirectoryPrefix,
				"directory-prefix",
				flags.DirectoryPrefix,
				"Directory prefix in bundle when original-filenames=true",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("directory-prefix") && cfg.DirectoryPrefix != nil {
				flags.DirectoryPrefix = *cfg.DirectoryPrefix
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.directory-prefix", func(val string) {
				cfg.DirectoryPrefix = &val
			})
		},
		ApplyToRequest: reqString("directory_prefix", func(f *Flags) string {
			return f.DirectoryPrefix
		}),
	},
	{
		FlagName:  "all-platforms",
		ConfigKey: "download.all-platforms",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.AllPlatforms,
				"all-platforms",
				flags.AllPlatforms,
				"Include all platform keys",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("all-platforms") && cfg.AllPlatforms != nil {
				flags.AllPlatforms = *cfg.AllPlatforms
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.all-platforms", func(val bool) {
				cfg.AllPlatforms = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"all-platforms",
			"all_platforms",
			func(f *Flags) bool { return f.AllPlatforms },
			func(c *DownloadConfig) *bool { return c.AllPlatforms },
		),
	},
	{
		FlagName:  "filter-langs",
		ConfigKey: "download.filter-langs",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringSliceVar(
				&flags.FilterLangs,
				"filter-langs",
				flags.FilterLangs,
				"Languages to export",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("filter-langs") && cfg.FilterLangs != nil {
				flags.FilterLangs = append([]string(nil), (*cfg.FilterLangs)...)
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigStringSlice(v, "download.filter-langs", func(val []string) {
				cfg.FilterLangs = &val
			})
		},
		ApplyToRequest: reqStringSlice("filter_langs", func(f *Flags) []string {
			return f.FilterLangs
		}),
	},
	{
		FlagName:  "filter-data",
		ConfigKey: "download.filter-data",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringSliceVar(
				&flags.FilterData,
				"filter-data",
				flags.FilterData,
				"Narrow export data range",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("filter-data") && cfg.FilterData != nil {
				flags.FilterData = append([]string(nil), (*cfg.FilterData)...)
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigStringSlice(v, "download.filter-data", func(val []string) {
				cfg.FilterData = &val
			})
		},
		ApplyToRequest: reqStringSlice("filter_data", func(f *Flags) []string {
			return f.FilterData
		}),
	},
	{
		FlagName:  "filter-filenames",
		ConfigKey: "download.filter-filenames",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringSliceVar(
				&flags.FilterFilenames,
				"filter-filenames",
				flags.FilterFilenames,
				"Only include keys attributed to selected files",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("filter-filenames") && cfg.FilterFilenames != nil {
				flags.FilterFilenames = append([]string(nil), (*cfg.FilterFilenames)...)
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigStringSlice(v, "download.filter-filenames", func(val []string) {
				cfg.FilterFilenames = &val
			})
		},
		ApplyToRequest: reqStringSlice("filter_filenames", func(f *Flags) []string {
			return f.FilterFilenames
		}),
	},
	{
		FlagName:  "custom-translation-status-ids",
		ConfigKey: "download.custom-translation-status-ids",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringSliceVar(
				&flags.CustomTranslationStatusIDs,
				"custom-translation-status-ids",
				flags.CustomTranslationStatusIDs,
				"Only include translations with selected custom status IDs",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("custom-translation-status-ids") && cfg.CustomTranslationStatusIDs != nil {
				flags.CustomTranslationStatusIDs = append([]string(nil), (*cfg.CustomTranslationStatusIDs)...)
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigStringSlice(v, "download.custom-translation-status-ids", func(val []string) {
				cfg.CustomTranslationStatusIDs = &val
			})
		},
		ApplyToRequest: reqStringSlice("custom_translation_status_ids", func(f *Flags) []string {
			return f.CustomTranslationStatusIDs
		}),
	},
	{
		FlagName:  "include-tags",
		ConfigKey: "download.include-tags",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringSliceVar(
				&flags.IncludeTags,
				"include-tags",
				flags.IncludeTags,
				"Only include keys with these tags",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("include-tags") && cfg.IncludeTags != nil {
				flags.IncludeTags = append([]string(nil), (*cfg.IncludeTags)...)
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigStringSlice(v, "download.include-tags", func(val []string) {
				cfg.IncludeTags = &val
			})
		},
		ApplyToRequest: reqStringSlice("include_tags", func(f *Flags) []string {
			return f.IncludeTags
		}),
	},
	{
		FlagName:  "exclude-tags",
		ConfigKey: "download.exclude-tags",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringSliceVar(
				&flags.ExcludeTags,
				"exclude-tags",
				flags.ExcludeTags,
				"Exclude keys with these tags",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("exclude-tags") && cfg.ExcludeTags != nil {
				flags.ExcludeTags = append([]string(nil), (*cfg.ExcludeTags)...)
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigStringSlice(v, "download.exclude-tags", func(val []string) {
				cfg.ExcludeTags = &val
			})
		},
		ApplyToRequest: reqStringSlice("exclude_tags", func(f *Flags) []string {
			return f.ExcludeTags
		}),
	},
	{
		FlagName:  "include-pids",
		ConfigKey: "download.include-pids",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringSliceVar(
				&flags.IncludePIDs,
				"include-pids",
				flags.IncludePIDs,
				"Include keys from other project IDs",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("include-pids") && cfg.IncludePIDs != nil {
				flags.IncludePIDs = append([]string(nil), (*cfg.IncludePIDs)...)
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigStringSlice(v, "download.include-pids", func(val []string) {
				cfg.IncludePIDs = &val
			})
		},
		ApplyToRequest: reqStringSlice("include_pids", func(f *Flags) []string {
			return f.IncludePIDs
		}),
	},
	{
		FlagName:  "triggers",
		ConfigKey: "download.triggers",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringSliceVar(
				&flags.Triggers,
				"triggers",
				flags.Triggers,
				"Trigger integration exports",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("triggers") && cfg.Triggers != nil {
				flags.Triggers = append([]string(nil), (*cfg.Triggers)...)
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigStringSlice(v, "download.triggers", func(val []string) {
				cfg.Triggers = &val
			})
		},
		ApplyToRequest: reqStringSlice("triggers", func(f *Flags) []string {
			return f.Triggers
		}),
	},
	{
		FlagName:  "filter-repositories",
		ConfigKey: "download.filter-repositories",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringSliceVar(
				&flags.FilterRepositories,
				"filter-repositories",
				flags.FilterRepositories,
				"Only process selected repositories in organization/repository format",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("filter-repositories") && cfg.FilterRepositories != nil {
				flags.FilterRepositories = append([]string(nil), (*cfg.FilterRepositories)...)
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigStringSlice(v, "download.filter-repositories", func(val []string) {
				cfg.FilterRepositories = &val
			})
		},
		ApplyToRequest: reqStringSlice("filter_repositories", func(f *Flags) []string {
			return f.FilterRepositories
		}),
	},
	{
		FlagName:  "filter-task-id",
		ConfigKey: "download.filter-task-id",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().Int64Var(
				&flags.FilterTaskID,
				"filter-task-id",
				flags.FilterTaskID,
				"Only include keys attributed to this task (offline_xliff only)",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("filter-task-id") && cfg.FilterTaskID != nil {
				flags.FilterTaskID = *cfg.FilterTaskID
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.filter-task-id", func(val int64) {
				cfg.FilterTaskID = &val
			})
		},
		ApplyToRequest: reqInt64WithDefault(
			"filter-task-id",
			"filter_task_id",
			func(f *Flags) int64 { return f.FilterTaskID },
			func(c *DownloadConfig) *int64 { return c.FilterTaskID },
		),
	},
	{
		FlagName:  "add-newline-eof",
		ConfigKey: "download.add-newline-eof",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.AddNewlineEOF,
				"add-newline-eof",
				flags.AddNewlineEOF,
				"Add newline at end of file when supported",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("add-newline-eof") && cfg.AddNewlineEOF != nil {
				flags.AddNewlineEOF = *cfg.AddNewlineEOF
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.add-newline-eof", func(val bool) {
				cfg.AddNewlineEOF = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"add-newline-eof",
			"add_newline_eof",
			func(f *Flags) bool { return f.AddNewlineEOF },
			func(c *DownloadConfig) *bool { return c.AddNewlineEOF },
		),
	},
	{
		FlagName:  "include-comments",
		ConfigKey: "download.include-comments",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.IncludeComments,
				"include-comments",
				flags.IncludeComments,
				"Include key comments and description when supported",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("include-comments") && cfg.IncludeComments != nil {
				flags.IncludeComments = *cfg.IncludeComments
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.include-comments", func(val bool) {
				cfg.IncludeComments = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"include-comments",
			"include_comments",
			func(f *Flags) bool { return f.IncludeComments },
			func(c *DownloadConfig) *bool { return c.IncludeComments },
		),
	},
	{
		FlagName:  "include-description",
		ConfigKey: "download.include-description",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.IncludeDescription,
				"include-description",
				flags.IncludeDescription,
				"Include key description when supported",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("include-description") && cfg.IncludeDescription != nil {
				flags.IncludeDescription = *cfg.IncludeDescription
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.include-description", func(val bool) {
				cfg.IncludeDescription = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"include-description",
			"include_description",
			func(f *Flags) bool { return f.IncludeDescription },
			func(c *DownloadConfig) *bool { return c.IncludeDescription },
		),
	},
	{
		FlagName:  "replace-breaks",
		ConfigKey: "download.replace-breaks",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.ReplaceBreaks,
				"replace-breaks",
				flags.ReplaceBreaks,
				"Replace line breaks in exported translations with \\n",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("replace-breaks") && cfg.ReplaceBreaks != nil {
				flags.ReplaceBreaks = *cfg.ReplaceBreaks
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.replace-breaks", func(val bool) {
				cfg.ReplaceBreaks = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"replace-breaks",
			"replace_breaks",
			func(f *Flags) bool { return f.ReplaceBreaks },
			func(c *DownloadConfig) *bool { return c.ReplaceBreaks },
		),
	},
	{
		FlagName:  "disable-references",
		ConfigKey: "download.disable-references",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.DisableReferences,
				"disable-references",
				flags.DisableReferences,
				"Disable automatic replacement of key reference placeholders",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("disable-references") && cfg.DisableReferences != nil {
				flags.DisableReferences = *cfg.DisableReferences
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.disable-references", func(val bool) {
				cfg.DisableReferences = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"disable-references",
			"disable_references",
			func(f *Flags) bool { return f.DisableReferences },
			func(c *DownloadConfig) *bool { return c.DisableReferences },
		),
	},
	{
		FlagName:  "icu-numeric",
		ConfigKey: "download.icu-numeric",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.ICUNumeric,
				"icu-numeric",
				flags.ICUNumeric,
				"Replace ICU plural forms zero/one/two with =0/=1/=2",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("icu-numeric") && cfg.ICUNumeric != nil {
				flags.ICUNumeric = *cfg.ICUNumeric
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.icu-numeric", func(val bool) {
				cfg.ICUNumeric = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"icu-numeric",
			"icu_numeric",
			func(f *Flags) bool { return f.ICUNumeric },
			func(c *DownloadConfig) *bool { return c.ICUNumeric },
		),
	},
	{
		FlagName:  "escape-percent",
		ConfigKey: "download.escape-percent",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.EscapePercent,
				"escape-percent",
				flags.EscapePercent,
				"Escape universal percent placeholders for printf format",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("escape-percent") && cfg.EscapePercent != nil {
				flags.EscapePercent = *cfg.EscapePercent
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.escape-percent", func(val bool) {
				cfg.EscapePercent = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"escape-percent",
			"escape_percent",
			func(f *Flags) bool { return f.EscapePercent },
			func(c *DownloadConfig) *bool { return c.EscapePercent },
		),
	},
	{
		FlagName:  "yaml-include-root",
		ConfigKey: "download.yaml-include-root",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.YAMLIncludeRoot,
				"yaml-include-root",
				flags.YAMLIncludeRoot,
				"Include language ISO code as root key for YAML export",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("yaml-include-root") && cfg.YAMLIncludeRoot != nil {
				flags.YAMLIncludeRoot = *cfg.YAMLIncludeRoot
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.yaml-include-root", func(val bool) {
				cfg.YAMLIncludeRoot = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"yaml-include-root",
			"yaml_include_root",
			func(f *Flags) bool { return f.YAMLIncludeRoot },
			func(c *DownloadConfig) *bool { return c.YAMLIncludeRoot },
		),
	},
	{
		FlagName:  "json-unescaped-slashes",
		ConfigKey: "download.json-unescaped-slashes",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.JSONUnescapedSlashes,
				"json-unescaped-slashes",
				flags.JSONUnescapedSlashes,
				"Leave forward slashes unescaped in JSON export",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("json-unescaped-slashes") && cfg.JSONUnescapedSlashes != nil {
				flags.JSONUnescapedSlashes = *cfg.JSONUnescapedSlashes
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.json-unescaped-slashes", func(val bool) {
				cfg.JSONUnescapedSlashes = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"json-unescaped-slashes",
			"json_unescaped_slashes",
			func(f *Flags) bool { return f.JSONUnescapedSlashes },
			func(c *DownloadConfig) *bool { return c.JSONUnescapedSlashes },
		),
	},
	{
		FlagName:  "compact",
		ConfigKey: "download.compact",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.Compact,
				"compact",
				flags.Compact,
				"Export compact ARB structure",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("compact") && cfg.Compact != nil {
				flags.Compact = *cfg.Compact
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.compact", func(val bool) {
				cfg.Compact = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"compact",
			"compact",
			func(f *Flags) bool { return f.Compact },
			func(c *DownloadConfig) *bool { return c.Compact },
		),
	},
	{
		FlagName:  "export-sort",
		ConfigKey: "download.export-sort",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.ExportSort,
				"export-sort",
				flags.ExportSort,
				"Export key sort mode",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("export-sort") && cfg.ExportSort != nil {
				flags.ExportSort = *cfg.ExportSort
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.export-sort", func(val string) {
				cfg.ExportSort = &val
			})
		},
		ApplyToRequest: reqString("export_sort", func(f *Flags) string {
			return f.ExportSort
		}),
	},
	{
		FlagName:  "export-empty-as",
		ConfigKey: "download.export-empty-as",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.ExportEmptyAs,
				"export-empty-as",
				flags.ExportEmptyAs,
				"How to export empty translations",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("export-empty-as") && cfg.ExportEmptyAs != nil {
				flags.ExportEmptyAs = *cfg.ExportEmptyAs
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.export-empty-as", func(val string) {
				cfg.ExportEmptyAs = &val
			})
		},
		ApplyToRequest: reqString("export_empty_as", func(f *Flags) string {
			return f.ExportEmptyAs
		}),
	},
	{
		FlagName:  "export-null-as",
		ConfigKey: "download.export-null-as",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.ExportNullAs,
				"export-null-as",
				flags.ExportNullAs,
				"How to export null translations (Ruby on Rails YAML only)",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("export-null-as") && cfg.ExportNullAs != nil {
				flags.ExportNullAs = *cfg.ExportNullAs
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.export-null-as", func(val string) {
				cfg.ExportNullAs = &val
			})
		},
		ApplyToRequest: reqString("export_null_as", func(f *Flags) string {
			return f.ExportNullAs
		}),
	},
	{
		FlagName:  "plural-format",
		ConfigKey: "download.plural-format",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.PluralFormat,
				"plural-format",
				flags.PluralFormat,
				"Override default plural format",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("plural-format") && cfg.PluralFormat != nil {
				flags.PluralFormat = *cfg.PluralFormat
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.plural-format", func(val string) {
				cfg.PluralFormat = &val
			})
		},
		ApplyToRequest: reqString("plural_format", func(f *Flags) string {
			return f.PluralFormat
		}),
	},
	{
		FlagName:  "placeholder-format",
		ConfigKey: "download.placeholder-format",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.PlaceholderFormat,
				"placeholder-format",
				flags.PlaceholderFormat,
				"Override default placeholder format",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("placeholder-format") && cfg.PlaceholderFormat != nil {
				flags.PlaceholderFormat = *cfg.PlaceholderFormat
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.placeholder-format", func(val string) {
				cfg.PlaceholderFormat = &val
			})
		},
		ApplyToRequest: reqString("placeholder_format", func(f *Flags) string {
			return f.PlaceholderFormat
		}),
	},
	{
		FlagName:  "webhook-url",
		ConfigKey: "download.webhook-url",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.WebhookURL,
				"webhook-url",
				flags.WebhookURL,
				"Send POST with generated bundle URL to this URL when export completes",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("webhook-url") && cfg.WebhookURL != nil {
				flags.WebhookURL = *cfg.WebhookURL
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.webhook-url", func(val string) {
				cfg.WebhookURL = &val
			})
		},
		ApplyToRequest: reqString("webhook_url", func(f *Flags) string {
			return f.WebhookURL
		}),
	},
	{
		FlagName:  "indentation",
		ConfigKey: "download.indentation",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.Indentation,
				"indentation",
				flags.Indentation,
				"Override indentation in supported files",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("indentation") && cfg.Indentation != nil {
				flags.Indentation = *cfg.Indentation
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.indentation", func(val string) {
				cfg.Indentation = &val
			})
		},
		ApplyToRequest: reqString("indentation", func(f *Flags) string {
			return f.Indentation
		}),
	},
	{
		FlagName:  "java-properties-encoding",
		ConfigKey: "download.java-properties-encoding",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.JavaPropertiesEncoding,
				"java-properties-encoding",
				flags.JavaPropertiesEncoding,
				"Encoding for Java .properties export",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("java-properties-encoding") && cfg.JavaPropertiesEncoding != nil {
				flags.JavaPropertiesEncoding = *cfg.JavaPropertiesEncoding
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.java-properties-encoding", func(val string) {
				cfg.JavaPropertiesEncoding = &val
			})
		},
		ApplyToRequest: reqString("java_properties_encoding", func(f *Flags) string {
			return f.JavaPropertiesEncoding
		}),
	},
	{
		FlagName:  "java-properties-separator",
		ConfigKey: "download.java-properties-separator",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.JavaPropertiesSeparator,
				"java-properties-separator",
				flags.JavaPropertiesSeparator,
				"Separator for Java .properties export",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("java-properties-separator") && cfg.JavaPropertiesSeparator != nil {
				flags.JavaPropertiesSeparator = *cfg.JavaPropertiesSeparator
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.java-properties-separator", func(val string) {
				cfg.JavaPropertiesSeparator = &val
			})
		},
		ApplyToRequest: reqString("java_properties_separator", func(f *Flags) string {
			return f.JavaPropertiesSeparator
		}),
	},
	{
		FlagName:  "bundle-description",
		ConfigKey: "download.bundle-description",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.BundleDescription,
				"bundle-description",
				flags.BundleDescription,
				"Description for ios_sdk/android_sdk OTA SDK bundles",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("bundle-description") && cfg.BundleDescription != nil {
				flags.BundleDescription = *cfg.BundleDescription
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.bundle-description", func(val string) {
				cfg.BundleDescription = &val
			})
		},
		ApplyToRequest: reqString("bundle_description", func(f *Flags) string {
			return f.BundleDescription
		}),
	},
	{
		FlagName:  "language-mapping",
		ConfigKey: "download.language-mapping",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.LanguageMappingJSON,
				"language-mapping",
				flags.LanguageMappingJSON,
				"Language mapping as JSON array of objects",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *DownloadConfig) {
			if !cmd.Flags().Changed("language-mapping") && cfg.LanguageMappingJSON != nil {
				flags.LanguageMappingJSON = *cfg.LanguageMappingJSON
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *DownloadConfig) {
			vh.ApplyConfigValue(v, "download.language-mapping", func(val string) {
				cfg.LanguageMappingJSON = &val
			})
		},
		ApplyToRequest: reqLanguageMapping("language-mapping", func(f *Flags) string {
			return f.LanguageMappingJSON
		}),
	},
}
