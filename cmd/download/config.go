package download

import (
	"time"
)

type DownloadConfig struct {
	Out            *string        `mapstructure:"out"`
	Format         *string        `mapstructure:"format"`
	ContextTimeout *time.Duration `mapstructure:"context-timeout"`
	Async          *bool          `mapstructure:"async"`

	OriginalFilenames          *bool     `mapstructure:"original-filenames"`
	BundleStructure            *string   `mapstructure:"bundle-structure"`
	DirectoryPrefix            *string   `mapstructure:"directory-prefix"`
	AllPlatforms               *bool     `mapstructure:"all-platforms"`
	FilterLangs                *[]string `mapstructure:"filter-langs"`
	FilterData                 *[]string `mapstructure:"filter-data"`
	FilterFilenames            *[]string `mapstructure:"filter-filenames"`
	AddNewlineEOF              *bool     `mapstructure:"add-newline-eof"`
	CustomTranslationStatusIDs *[]string `mapstructure:"custom-translation-status-ids"`
	IncludeTags                *[]string `mapstructure:"include-tags"`
	ExcludeTags                *[]string `mapstructure:"exclude-tags"`
	ExportSort                 *string   `mapstructure:"export-sort"`
	ExportEmptyAs              *string   `mapstructure:"export-empty-as"`
	ExportNullAs               *string   `mapstructure:"export-null-as"`
	IncludeComments            *bool     `mapstructure:"include-comments"`
	IncludeDescription         *bool     `mapstructure:"include-description"`
	IncludePIDs                *[]string `mapstructure:"include-pids"`
	Triggers                   *[]string `mapstructure:"triggers"`
	FilterRepositories         *[]string `mapstructure:"filter-repositories"`
	ReplaceBreaks              *bool     `mapstructure:"replace-breaks"`
	DisableReferences          *bool     `mapstructure:"disable-references"`
	PluralFormat               *string   `mapstructure:"plural-format"`
	PlaceholderFormat          *string   `mapstructure:"placeholder-format"`
	WebhookURL                 *string   `mapstructure:"webhook-url"`
	LanguageMappingJSON        *string   `mapstructure:"language-mapping"`
	ICUNumeric                 *bool     `mapstructure:"icu-numeric"`
	EscapePercent              *bool     `mapstructure:"escape-percent"`
	Indentation                *string   `mapstructure:"indentation"`
	YAMLIncludeRoot            *bool     `mapstructure:"yaml-include-root"`
	JSONUnescapedSlashes       *bool     `mapstructure:"json-unescaped-slashes"`
	JavaPropertiesEncoding     *string   `mapstructure:"java-properties-encoding"`
	JavaPropertiesSeparator    *string   `mapstructure:"java-properties-separator"`
	BundleDescription          *string   `mapstructure:"bundle-description"`
	FilterTaskID               *int64    `mapstructure:"filter-task-id"`
	Compact                    *bool     `mapstructure:"compact"`
}
