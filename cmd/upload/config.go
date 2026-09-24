package upload

type UploadConfig struct {
	// Command behavior.
	Manifest *string `mapstructure:"manifest"`
	SrcPath  *string `mapstructure:"src-path"`
	Poll     *bool   `mapstructure:"poll"`

	// Upload request.
	Filename *string `mapstructure:"filename"`
	Data     *string `mapstructure:"data"`
	LangISO  *string `mapstructure:"lang-iso"`
	Format   *string `mapstructure:"format"`

	Tags                       *[]string `mapstructure:"tags"`
	CustomTranslationStatusIDs *[]string `mapstructure:"custom-translation-status-ids"`

	ConvertPlaceholders    *bool `mapstructure:"convert-placeholders"`
	DetectICUPlurals       *bool `mapstructure:"detect-icu-plurals"`
	TagInsertedKeys        *bool `mapstructure:"tag-inserted-keys"`
	TagUpdatedKeys         *bool `mapstructure:"tag-updated-keys"`
	TagSkippedKeys         *bool `mapstructure:"tag-skipped-keys"`
	ReplaceModified        *bool `mapstructure:"replace-modified"`
	SlashNToLinebreak      *bool `mapstructure:"slashn-to-linebreak"`
	KeysToValues           *bool `mapstructure:"keys-to-values"`
	DistinguishByFile      *bool `mapstructure:"distinguish-by-file"`
	ApplyTM                *bool `mapstructure:"apply-tm"`
	UseAutomations         *bool `mapstructure:"use-automations"`
	HiddenFromContributors *bool `mapstructure:"hidden-from-contributors"`
	CleanupMode            *bool `mapstructure:"cleanup-mode"`
	SkipDetectLangISO      *bool `mapstructure:"skip-detect-lang-iso"`

	CustomTranslationStatusInsertedKeys *bool `mapstructure:"custom-translation-status-inserted-keys"`
	CustomTranslationStatusUpdatedKeys  *bool `mapstructure:"custom-translation-status-updated-keys"`
	CustomTranslationStatusSkippedKeys  *bool `mapstructure:"custom-translation-status-skipped-keys"`

	FilterTaskID *int64 `mapstructure:"filter-task-id"`
}
