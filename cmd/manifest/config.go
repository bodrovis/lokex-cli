package manifest

type GenerateConfig struct {
	Paths           *[]string `mapstructure:"path"`
	ExcludePatterns *[]string `mapstructure:"exclude-pattern"`

	NamePattern     *string `mapstructure:"name-pattern"`
	FilenamePattern *string `mapstructure:"filename-pattern"`

	BaseLang        *string `mapstructure:"base-lang"`
	LanguageMapping *string `mapstructure:"language-mapping"`

	Out *string `mapstructure:"out"`
}
