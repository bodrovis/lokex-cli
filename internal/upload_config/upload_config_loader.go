package upload_config

import (
	"fmt"
	"time"

	vh "github.com/bodrovis/lokex-cli/internal/viper_helpers"
	"github.com/spf13/viper"
)

func LoadUploadConfig(cfg *UploadConfig, configFile, envPrefix string) error {
	v := vh.NewConfigViper(configFile, envPrefix)

	if err := bindUploadEnv(v); err != nil {
		return fmt.Errorf("bind upload env: %w", err)
	}

	if err := vh.ReadOptionalConfig(v, configFile); err != nil {
		return fmt.Errorf("read upload config: %w", err)
	}

	vh.ApplyConfigValue(v, "upload.filename", func(val string) { cfg.Filename = &val })
	vh.ApplyConfigValue(v, "upload.src-path", func(val string) { cfg.SrcPath = &val })
	vh.ApplyConfigValue(v, "upload.data", func(val string) { cfg.Data = &val })
	vh.ApplyConfigValue(v, "upload.lang-iso", func(val string) { cfg.LangISO = &val })
	vh.ApplyConfigValue(v, "upload.poll", func(val bool) { cfg.Poll = &val })
	vh.ApplyConfigValue(v, "upload.context-timeout", func(val time.Duration) { cfg.ContextTimeout = &val })

	vh.ApplyConfigValue(v, "upload.convert-placeholders", func(val bool) { cfg.ConvertPlaceholders = &val })
	vh.ApplyConfigValue(v, "upload.detect-icu-plurals", func(val bool) { cfg.DetectICUPlurals = &val })
	vh.ApplyConfigValue(v, "upload.tag-inserted-keys", func(val bool) { cfg.TagInsertedKeys = &val })
	vh.ApplyConfigValue(v, "upload.tag-updated-keys", func(val bool) { cfg.TagUpdatedKeys = &val })
	vh.ApplyConfigValue(v, "upload.tag-skipped-keys", func(val bool) { cfg.TagSkippedKeys = &val })
	vh.ApplyConfigValue(v, "upload.replace-modified", func(val bool) { cfg.ReplaceModified = &val })
	vh.ApplyConfigValue(v, "upload.slashn-to-linebreak", func(val bool) { cfg.SlashNToLinebreak = &val })
	vh.ApplyConfigValue(v, "upload.keys-to-values", func(val bool) { cfg.KeysToValues = &val })
	vh.ApplyConfigValue(v, "upload.distinguish-by-file", func(val bool) { cfg.DistinguishByFile = &val })
	vh.ApplyConfigValue(v, "upload.apply-tm", func(val bool) { cfg.ApplyTM = &val })
	vh.ApplyConfigValue(v, "upload.use-automations", func(val bool) { cfg.UseAutomations = &val })
	vh.ApplyConfigValue(v, "upload.hidden-from-contributors", func(val bool) { cfg.HiddenFromContributors = &val })
	vh.ApplyConfigValue(v, "upload.cleanup-mode", func(val bool) { cfg.CleanupMode = &val })
	vh.ApplyConfigValue(v, "upload.custom-translation-status-inserted-keys", func(val bool) {
		cfg.CustomTranslationStatusInsertedKeys = &val
	})
	vh.ApplyConfigValue(v, "upload.custom-translation-status-updated-keys", func(val bool) {
		cfg.CustomTranslationStatusUpdatedKeys = &val
	})
	vh.ApplyConfigValue(v, "upload.custom-translation-status-skipped-keys", func(val bool) {
		cfg.CustomTranslationStatusSkippedKeys = &val
	})
	vh.ApplyConfigValue(v, "upload.skip-detect-lang-iso", func(val bool) { cfg.SkipDetectLangISO = &val })
	vh.ApplyConfigValue(v, "upload.format", func(val string) { cfg.Format = &val })
	vh.ApplyConfigValue(v, "upload.filter-task-id", func(val int64) { cfg.FilterTaskID = &val })

	vh.ApplyConfigStringSlice(v, "upload.tags", func(val []string) {
		cfg.Tags = &val
	})
	vh.ApplyConfigStringSlice(v, "upload.custom-translation-status-ids", func(val []string) {
		cfg.CustomTranslationStatusIDs = &val
	})

	return nil
}

func bindUploadEnv(v *viper.Viper) error {
	keys := []string{
		"upload.filename",
		"upload.src-path",
		"upload.data",
		"upload.lang-iso",
		"upload.poll",
		"upload.context-timeout",

		"upload.convert-placeholders",
		"upload.detect-icu-plurals",
		"upload.tags",
		"upload.tag-inserted-keys",
		"upload.tag-updated-keys",
		"upload.tag-skipped-keys",
		"upload.replace-modified",
		"upload.slashn-to-linebreak",
		"upload.keys-to-values",
		"upload.distinguish-by-file",
		"upload.apply-tm",
		"upload.use-automations",
		"upload.hidden-from-contributors",
		"upload.cleanup-mode",
		"upload.custom-translation-status-ids",
		"upload.custom-translation-status-inserted-keys",
		"upload.custom-translation-status-updated-keys",
		"upload.custom-translation-status-skipped-keys",
		"upload.skip-detect-lang-iso",
		"upload.format",
		"upload.filter-task-id",
	}

	return vh.BindEnvKeys(v, keys)
}
