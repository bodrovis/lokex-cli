package upload

import (
	params "github.com/bodrovis/lokex-cli/internal/params"
	vh "github.com/bodrovis/lokex-cli/internal/viper_helpers"

	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type UploadParamSpec = params.ParamSpec[Flags, UploadConfig, lokexupload.UploadParams]

var uploadParamSpecs = []UploadParamSpec{
	{
		FlagName:  "filename",
		ConfigKey: "upload.filename",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.Filename,
				"filename",
				flags.Filename,
				"Filename sent to Lokalise (required)",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("filename") && cfg.Filename != nil {
				flags.Filename = *cfg.Filename
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.filename", func(val string) {
				cfg.Filename = &val
			})
		},
		ApplyToRequest: reqNormalizedFilename("filename", func(f *Flags) string {
			return f.Filename
		}),
	},
	{
		FlagName:  "src-path",
		ConfigKey: "upload.src-path",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.SrcPath,
				"src-path",
				flags.SrcPath,
				"Local path to read file contents from (optional)",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("src-path") && cfg.SrcPath != nil {
				flags.SrcPath = *cfg.SrcPath
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.src-path", func(val string) {
				cfg.SrcPath = &val
			})
		},
	},
	{
		FlagName:  "data",
		ConfigKey: "upload.data",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.Data,
				"data",
				flags.Data,
				"Base64-encoded file contents (optional; if set, file is not read from disk)",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("data") && cfg.Data != nil {
				flags.Data = *cfg.Data
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.data", func(val string) {
				cfg.Data = &val
			})
		},
		ApplyToRequest: reqString("data", func(f *Flags) string {
			return f.Data
		}),
	},
	{
		FlagName:  "lang-iso",
		ConfigKey: "upload.lang-iso",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.LangISO,
				"lang-iso",
				flags.LangISO,
				"Language code of the translations in the file (required)",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("lang-iso") && cfg.LangISO != nil {
				flags.LangISO = *cfg.LangISO
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.lang-iso", func(val string) {
				cfg.LangISO = &val
			})
		},
		ApplyToRequest: reqTrimmedString("lang_iso", func(f *Flags) string {
			return f.LangISO
		}),
	},
	{
		FlagName:  "poll",
		ConfigKey: "upload.poll",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.Poll,
				"poll",
				flags.Poll,
				"Wait until Lokalise finishes processing the upload",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("poll") && cfg.Poll != nil {
				flags.Poll = *cfg.Poll
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.poll", func(val bool) {
				cfg.Poll = &val
			})
		},
	},
	{
		FlagName:  "format",
		ConfigKey: "upload.format",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringVar(
				&flags.Format,
				"format",
				flags.Format,
				"File format (e.g. json, strings, xml)",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("format") && cfg.Format != nil {
				flags.Format = *cfg.Format
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.format", func(val string) {
				cfg.Format = &val
			})
		},
		ApplyToRequest: reqString("format", func(f *Flags) string {
			return f.Format
		}),
	},
	{
		FlagName:  "tags",
		ConfigKey: "upload.tags",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringSliceVar(
				&flags.Tags,
				"tags",
				flags.Tags,
				"Tags to apply to keys",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("tags") && cfg.Tags != nil {
				flags.Tags = append([]string(nil), (*cfg.Tags)...)
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigStringSlice(v, "upload.tags", func(val []string) {
				cfg.Tags = &val
			})
		},
		ApplyToRequest: reqStringSlice("tags", func(f *Flags) []string {
			return f.Tags
		}),
	},
	{
		FlagName:  "custom-translation-status-ids",
		ConfigKey: "upload.custom-translation-status-ids",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().StringSliceVar(
				&flags.CustomTranslationStatusIDs,
				"custom-translation-status-ids",
				flags.CustomTranslationStatusIDs,
				"Custom translation status IDs to add",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("custom-translation-status-ids") && cfg.CustomTranslationStatusIDs != nil {
				flags.CustomTranslationStatusIDs = append([]string(nil), (*cfg.CustomTranslationStatusIDs)...)
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigStringSlice(v, "upload.custom-translation-status-ids", func(val []string) {
				cfg.CustomTranslationStatusIDs = &val
			})
		},
		ApplyToRequest: reqStringSlice("custom_translation_status_ids", func(f *Flags) []string {
			return f.CustomTranslationStatusIDs
		}),
	},
	{
		FlagName:  "convert-placeholders",
		ConfigKey: "upload.convert-placeholders",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.ConvertPlaceholders,
				"convert-placeholders",
				flags.ConvertPlaceholders,
				"Convert placeholders to Lokalise universal placeholders",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("convert-placeholders") && cfg.ConvertPlaceholders != nil {
				flags.ConvertPlaceholders = *cfg.ConvertPlaceholders
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.convert-placeholders", func(val bool) {
				cfg.ConvertPlaceholders = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"convert-placeholders",
			"convert_placeholders",
			func(f *Flags) bool { return f.ConvertPlaceholders },
			func(c *UploadConfig) *bool { return c.ConvertPlaceholders },
		),
	},
	{
		FlagName:  "detect-icu-plurals",
		ConfigKey: "upload.detect-icu-plurals",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.DetectICUPlurals,
				"detect-icu-plurals",
				flags.DetectICUPlurals,
				"Automatically detect and parse ICU plurals",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("detect-icu-plurals") && cfg.DetectICUPlurals != nil {
				flags.DetectICUPlurals = *cfg.DetectICUPlurals
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.detect-icu-plurals", func(val bool) {
				cfg.DetectICUPlurals = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"detect-icu-plurals",
			"detect_icu_plurals",
			func(f *Flags) bool { return f.DetectICUPlurals },
			func(c *UploadConfig) *bool { return c.DetectICUPlurals },
		),
	},
	{
		FlagName:  "tag-inserted-keys",
		ConfigKey: "upload.tag-inserted-keys",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.TagInsertedKeys,
				"tag-inserted-keys",
				flags.TagInsertedKeys,
				"Add tags to inserted keys",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("tag-inserted-keys") && cfg.TagInsertedKeys != nil {
				flags.TagInsertedKeys = *cfg.TagInsertedKeys
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.tag-inserted-keys", func(val bool) {
				cfg.TagInsertedKeys = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"tag-inserted-keys",
			"tag_inserted_keys",
			func(f *Flags) bool { return f.TagInsertedKeys },
			func(c *UploadConfig) *bool { return c.TagInsertedKeys },
		),
	},
	{
		FlagName:  "tag-updated-keys",
		ConfigKey: "upload.tag-updated-keys",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.TagUpdatedKeys,
				"tag-updated-keys",
				flags.TagUpdatedKeys,
				"Add tags to updated keys",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("tag-updated-keys") && cfg.TagUpdatedKeys != nil {
				flags.TagUpdatedKeys = *cfg.TagUpdatedKeys
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.tag-updated-keys", func(val bool) {
				cfg.TagUpdatedKeys = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"tag-updated-keys",
			"tag_updated_keys",
			func(f *Flags) bool { return f.TagUpdatedKeys },
			func(c *UploadConfig) *bool { return c.TagUpdatedKeys },
		),
	},
	{
		FlagName:  "tag-skipped-keys",
		ConfigKey: "upload.tag-skipped-keys",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.TagSkippedKeys,
				"tag-skipped-keys",
				flags.TagSkippedKeys,
				"Add tags to skipped keys",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("tag-skipped-keys") && cfg.TagSkippedKeys != nil {
				flags.TagSkippedKeys = *cfg.TagSkippedKeys
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.tag-skipped-keys", func(val bool) {
				cfg.TagSkippedKeys = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"tag-skipped-keys",
			"tag_skipped_keys",
			func(f *Flags) bool { return f.TagSkippedKeys },
			func(c *UploadConfig) *bool { return c.TagSkippedKeys },
		),
	},
	{
		FlagName:  "replace-modified",
		ConfigKey: "upload.replace-modified",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.ReplaceModified,
				"replace-modified",
				flags.ReplaceModified,
				"Replace modified translations from the uploaded file",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("replace-modified") && cfg.ReplaceModified != nil {
				flags.ReplaceModified = *cfg.ReplaceModified
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.replace-modified", func(val bool) {
				cfg.ReplaceModified = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"replace-modified",
			"replace_modified",
			func(f *Flags) bool { return f.ReplaceModified },
			func(c *UploadConfig) *bool { return c.ReplaceModified },
		),
	},
	{
		FlagName:  "slashn-to-linebreak",
		ConfigKey: "upload.slashn-to-linebreak",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.SlashNToLinebreak,
				"slashn-to-linebreak",
				flags.SlashNToLinebreak,
				"Replace \\n with a real line break",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("slashn-to-linebreak") && cfg.SlashNToLinebreak != nil {
				flags.SlashNToLinebreak = *cfg.SlashNToLinebreak
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.slashn-to-linebreak", func(val bool) {
				cfg.SlashNToLinebreak = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"slashn-to-linebreak",
			"slashn_to_linebreak",
			func(f *Flags) bool { return f.SlashNToLinebreak },
			func(c *UploadConfig) *bool { return c.SlashNToLinebreak },
		),
	},
	{
		FlagName:  "keys-to-values",
		ConfigKey: "upload.keys-to-values",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.KeysToValues,
				"keys-to-values",
				flags.KeysToValues,
				"Replace values with key names",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("keys-to-values") && cfg.KeysToValues != nil {
				flags.KeysToValues = *cfg.KeysToValues
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.keys-to-values", func(val bool) {
				cfg.KeysToValues = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"keys-to-values",
			"keys_to_values",
			func(f *Flags) bool { return f.KeysToValues },
			func(c *UploadConfig) *bool { return c.KeysToValues },
		),
	},
	{
		FlagName:  "distinguish-by-file",
		ConfigKey: "upload.distinguish-by-file",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.DistinguishByFile,
				"distinguish-by-file",
				flags.DistinguishByFile,
				"Allow same key names to coexist across different filenames",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("distinguish-by-file") && cfg.DistinguishByFile != nil {
				flags.DistinguishByFile = *cfg.DistinguishByFile
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.distinguish-by-file", func(val bool) {
				cfg.DistinguishByFile = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"distinguish-by-file",
			"distinguish_by_file",
			func(f *Flags) bool { return f.DistinguishByFile },
			func(c *UploadConfig) *bool { return c.DistinguishByFile },
		),
	},
	{
		FlagName:  "apply-tm",
		ConfigKey: "upload.apply-tm",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.ApplyTM,
				"apply-tm",
				flags.ApplyTM,
				"Apply 100% translation memory matches",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("apply-tm") && cfg.ApplyTM != nil {
				flags.ApplyTM = *cfg.ApplyTM
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.apply-tm", func(val bool) {
				cfg.ApplyTM = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"apply-tm",
			"apply_tm",
			func(f *Flags) bool { return f.ApplyTM },
			func(c *UploadConfig) *bool { return c.ApplyTM },
		),
	},
	{
		FlagName:  "use-automations",
		ConfigKey: "upload.use-automations",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.UseAutomations,
				"use-automations",
				flags.UseAutomations,
				"Run automations for this upload",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("use-automations") && cfg.UseAutomations != nil {
				flags.UseAutomations = *cfg.UseAutomations
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.use-automations", func(val bool) {
				cfg.UseAutomations = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"use-automations",
			"use_automations",
			func(f *Flags) bool { return f.UseAutomations },
			func(c *UploadConfig) *bool { return c.UseAutomations },
		),
	},
	{
		FlagName:  "hidden-from-contributors",
		ConfigKey: "upload.hidden-from-contributors",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.HiddenFromContributors,
				"hidden-from-contributors",
				flags.HiddenFromContributors,
				"Mark newly created keys as hidden from contributors",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("hidden-from-contributors") && cfg.HiddenFromContributors != nil {
				flags.HiddenFromContributors = *cfg.HiddenFromContributors
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.hidden-from-contributors", func(val bool) {
				cfg.HiddenFromContributors = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"hidden-from-contributors",
			"hidden_from_contributors",
			func(f *Flags) bool { return f.HiddenFromContributors },
			func(c *UploadConfig) *bool { return c.HiddenFromContributors },
		),
	},
	{
		FlagName:  "cleanup-mode",
		ConfigKey: "upload.cleanup-mode",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.CleanupMode,
				"cleanup-mode",
				flags.CleanupMode,
				"Delete keys/translations not present in the uploaded file",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("cleanup-mode") && cfg.CleanupMode != nil {
				flags.CleanupMode = *cfg.CleanupMode
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.cleanup-mode", func(val bool) {
				cfg.CleanupMode = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"cleanup-mode",
			"cleanup_mode",
			func(f *Flags) bool { return f.CleanupMode },
			func(c *UploadConfig) *bool { return c.CleanupMode },
		),
	},
	{
		FlagName:  "custom-translation-status-inserted-keys",
		ConfigKey: "upload.custom-translation-status-inserted-keys",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.CustomTranslationStatusInsertedKeys,
				"custom-translation-status-inserted-keys",
				flags.CustomTranslationStatusInsertedKeys,
				"Add custom statuses to inserted keys",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("custom-translation-status-inserted-keys") && cfg.CustomTranslationStatusInsertedKeys != nil {
				flags.CustomTranslationStatusInsertedKeys = *cfg.CustomTranslationStatusInsertedKeys
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.custom-translation-status-inserted-keys", func(val bool) {
				cfg.CustomTranslationStatusInsertedKeys = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"custom-translation-status-inserted-keys",
			"custom_translation_status_inserted_keys",
			func(f *Flags) bool { return f.CustomTranslationStatusInsertedKeys },
			func(c *UploadConfig) *bool { return c.CustomTranslationStatusInsertedKeys },
		),
	},
	{
		FlagName:  "custom-translation-status-updated-keys",
		ConfigKey: "upload.custom-translation-status-updated-keys",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.CustomTranslationStatusUpdatedKeys,
				"custom-translation-status-updated-keys",
				flags.CustomTranslationStatusUpdatedKeys,
				"Add custom statuses to updated keys",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("custom-translation-status-updated-keys") && cfg.CustomTranslationStatusUpdatedKeys != nil {
				flags.CustomTranslationStatusUpdatedKeys = *cfg.CustomTranslationStatusUpdatedKeys
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.custom-translation-status-updated-keys", func(val bool) {
				cfg.CustomTranslationStatusUpdatedKeys = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"custom-translation-status-updated-keys",
			"custom_translation_status_updated_keys",
			func(f *Flags) bool { return f.CustomTranslationStatusUpdatedKeys },
			func(c *UploadConfig) *bool { return c.CustomTranslationStatusUpdatedKeys },
		),
	},
	{
		FlagName:  "custom-translation-status-skipped-keys",
		ConfigKey: "upload.custom-translation-status-skipped-keys",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.CustomTranslationStatusSkippedKeys,
				"custom-translation-status-skipped-keys",
				flags.CustomTranslationStatusSkippedKeys,
				"Add custom statuses to skipped keys",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("custom-translation-status-skipped-keys") && cfg.CustomTranslationStatusSkippedKeys != nil {
				flags.CustomTranslationStatusSkippedKeys = *cfg.CustomTranslationStatusSkippedKeys
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.custom-translation-status-skipped-keys", func(val bool) {
				cfg.CustomTranslationStatusSkippedKeys = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"custom-translation-status-skipped-keys",
			"custom_translation_status_skipped_keys",
			func(f *Flags) bool { return f.CustomTranslationStatusSkippedKeys },
			func(c *UploadConfig) *bool { return c.CustomTranslationStatusSkippedKeys },
		),
	},
	{
		FlagName:  "skip-detect-lang-iso",
		ConfigKey: "upload.skip-detect-lang-iso",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().BoolVar(
				&flags.SkipDetectLangISO,
				"skip-detect-lang-iso",
				flags.SkipDetectLangISO,
				"Skip automatic language detection by filename",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("skip-detect-lang-iso") && cfg.SkipDetectLangISO != nil {
				flags.SkipDetectLangISO = *cfg.SkipDetectLangISO
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.skip-detect-lang-iso", func(val bool) {
				cfg.SkipDetectLangISO = &val
			})
		},
		ApplyToRequest: reqBoolWithDefault(
			"skip-detect-lang-iso",
			"skip_detect_lang_iso",
			func(f *Flags) bool { return f.SkipDetectLangISO },
			func(c *UploadConfig) *bool { return c.SkipDetectLangISO },
		),
	},
	{
		FlagName:  "filter-task-id",
		ConfigKey: "upload.filter-task-id",
		BindFlag: func(cmd *cobra.Command, flags *Flags) {
			cmd.Flags().Int64Var(
				&flags.FilterTaskID,
				"filter-task-id",
				flags.FilterTaskID,
				"Apply import results as a part of a task (offline_xliff only)",
			)
		},
		ApplyDefault: func(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
			if !cmd.Flags().Changed("filter-task-id") && cfg.FilterTaskID != nil {
				flags.FilterTaskID = *cfg.FilterTaskID
			}
		},
		LoadFromViper: func(v *viper.Viper, cfg *UploadConfig) {
			vh.ApplyConfigValue(v, "upload.filter-task-id", func(val int64) {
				cfg.FilterTaskID = &val
			})
		},
		ApplyToRequest: reqInt64WithDefault(
			"filter-task-id",
			"filter_task_id",
			func(f *Flags) int64 { return f.FilterTaskID },
			func(c *UploadConfig) *int64 { return c.FilterTaskID },
		),
	},
}
