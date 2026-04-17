package upload

import (
	"github.com/spf13/cobra"

	uploadCfg "github.com/bodrovis/lokex-cli/internal/upload_config"
)

func applyDefaults(cmd *cobra.Command, flags *Flags, cfg *uploadCfg.UploadConfig) {
	if cfg == nil {
		return
	}

	if !cmd.Flags().Changed("filename") && cfg.Filename != nil {
		flags.Filename = *cfg.Filename
	}

	if !cmd.Flags().Changed("src-path") && cfg.SrcPath != nil {
		flags.SrcPath = *cfg.SrcPath
	}

	if !cmd.Flags().Changed("data") && cfg.Data != nil {
		flags.Data = *cfg.Data
	}

	if !cmd.Flags().Changed("lang-iso") && cfg.LangISO != nil {
		flags.LangISO = *cfg.LangISO
	}

	if !cmd.Flags().Changed("poll") && cfg.Poll != nil {
		flags.Poll = *cfg.Poll
	}

	if !cmd.Flags().Changed("context-timeout") && cfg.ContextTimeout != nil {
		flags.ContextTimeout = *cfg.ContextTimeout
	}

	if !cmd.Flags().Changed("format") && cfg.Format != nil {
		flags.Format = *cfg.Format
	}

	if !cmd.Flags().Changed("tags") && cfg.Tags != nil {
		flags.Tags = append([]string(nil), (*cfg.Tags)...)
	}

	if !cmd.Flags().Changed("custom-translation-status-ids") && cfg.CustomTranslationStatusIDs != nil {
		flags.CustomTranslationStatusIDs = append([]string(nil), (*cfg.CustomTranslationStatusIDs)...)
	}

	if !cmd.Flags().Changed("convert-placeholders") && cfg.ConvertPlaceholders != nil {
		flags.ConvertPlaceholders = *cfg.ConvertPlaceholders
	}

	if !cmd.Flags().Changed("detect-icu-plurals") && cfg.DetectICUPlurals != nil {
		flags.DetectICUPlurals = *cfg.DetectICUPlurals
	}

	if !cmd.Flags().Changed("tag-inserted-keys") && cfg.TagInsertedKeys != nil {
		flags.TagInsertedKeys = *cfg.TagInsertedKeys
	}

	if !cmd.Flags().Changed("tag-updated-keys") && cfg.TagUpdatedKeys != nil {
		flags.TagUpdatedKeys = *cfg.TagUpdatedKeys
	}

	if !cmd.Flags().Changed("tag-skipped-keys") && cfg.TagSkippedKeys != nil {
		flags.TagSkippedKeys = *cfg.TagSkippedKeys
	}

	if !cmd.Flags().Changed("replace-modified") && cfg.ReplaceModified != nil {
		flags.ReplaceModified = *cfg.ReplaceModified
	}

	if !cmd.Flags().Changed("slashn-to-linebreak") && cfg.SlashNToLinebreak != nil {
		flags.SlashNToLinebreak = *cfg.SlashNToLinebreak
	}

	if !cmd.Flags().Changed("keys-to-values") && cfg.KeysToValues != nil {
		flags.KeysToValues = *cfg.KeysToValues
	}

	if !cmd.Flags().Changed("distinguish-by-file") && cfg.DistinguishByFile != nil {
		flags.DistinguishByFile = *cfg.DistinguishByFile
	}

	if !cmd.Flags().Changed("apply-tm") && cfg.ApplyTM != nil {
		flags.ApplyTM = *cfg.ApplyTM
	}

	if !cmd.Flags().Changed("use-automations") && cfg.UseAutomations != nil {
		flags.UseAutomations = *cfg.UseAutomations
	}

	if !cmd.Flags().Changed("hidden-from-contributors") && cfg.HiddenFromContributors != nil {
		flags.HiddenFromContributors = *cfg.HiddenFromContributors
	}

	if !cmd.Flags().Changed("cleanup-mode") && cfg.CleanupMode != nil {
		flags.CleanupMode = *cfg.CleanupMode
	}

	if !cmd.Flags().Changed("custom-translation-status-inserted-keys") && cfg.CustomTranslationStatusInsertedKeys != nil {
		flags.CustomTranslationStatusInsertedKeys = *cfg.CustomTranslationStatusInsertedKeys
	}

	if !cmd.Flags().Changed("custom-translation-status-updated-keys") && cfg.CustomTranslationStatusUpdatedKeys != nil {
		flags.CustomTranslationStatusUpdatedKeys = *cfg.CustomTranslationStatusUpdatedKeys
	}

	if !cmd.Flags().Changed("custom-translation-status-skipped-keys") && cfg.CustomTranslationStatusSkippedKeys != nil {
		flags.CustomTranslationStatusSkippedKeys = *cfg.CustomTranslationStatusSkippedKeys
	}

	if !cmd.Flags().Changed("skip-detect-lang-iso") && cfg.SkipDetectLangISO != nil {
		flags.SkipDetectLangISO = *cfg.SkipDetectLangISO
	}

	if !cmd.Flags().Changed("filter-task-id") && cfg.FilterTaskID != nil {
		flags.FilterTaskID = *cfg.FilterTaskID
	}
}
