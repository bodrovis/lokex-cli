package upload

import (
	"time"
)

type UploadConfig struct {
	Filename       *string        `mapstructure:"filename"`
	SrcPath        *string        `mapstructure:"src-path"`
	Data           *string        `mapstructure:"data"`
	LangISO        *string        `mapstructure:"lang-iso"`
	Poll           *bool          `mapstructure:"poll"`
	ContextTimeout *time.Duration `mapstructure:"context-timeout"`

	ConvertPlaceholders                 *bool     `mapstructure:"convert-placeholders"`
	DetectICUPlurals                    *bool     `mapstructure:"detect-icu-plurals"`
	Tags                                *[]string `mapstructure:"tags"`
	TagInsertedKeys                     *bool     `mapstructure:"tag-inserted-keys"`
	TagUpdatedKeys                      *bool     `mapstructure:"tag-updated-keys"`
	TagSkippedKeys                      *bool     `mapstructure:"tag-skipped-keys"`
	ReplaceModified                     *bool     `mapstructure:"replace-modified"`
	SlashNToLinebreak                   *bool     `mapstructure:"slashn-to-linebreak"`
	KeysToValues                        *bool     `mapstructure:"keys-to-values"`
	DistinguishByFile                   *bool     `mapstructure:"distinguish-by-file"`
	ApplyTM                             *bool     `mapstructure:"apply-tm"`
	UseAutomations                      *bool     `mapstructure:"use-automations"`
	HiddenFromContributors              *bool     `mapstructure:"hidden-from-contributors"`
	CleanupMode                         *bool     `mapstructure:"cleanup-mode"`
	CustomTranslationStatusIDs          *[]string `mapstructure:"custom-translation-status-ids"`
	CustomTranslationStatusInsertedKeys *bool     `mapstructure:"custom-translation-status-inserted-keys"`
	CustomTranslationStatusUpdatedKeys  *bool     `mapstructure:"custom-translation-status-updated-keys"`
	CustomTranslationStatusSkippedKeys  *bool     `mapstructure:"custom-translation-status-skipped-keys"`
	SkipDetectLangISO                   *bool     `mapstructure:"skip-detect-lang-iso"`
	Format                              *string   `mapstructure:"format"`
	FilterTaskID                        *int64    `mapstructure:"filter-task-id"`
}
