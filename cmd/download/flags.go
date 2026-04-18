package download

import (
	"time"

	params "github.com/bodrovis/lokex-cli/internal/params"
	"github.com/spf13/cobra"
)

type Flags struct {
	Out            string
	Format         string
	ContextTimeout time.Duration
	Async          bool

	OriginalFilenames          bool
	BundleStructure            string
	DirectoryPrefix            string
	AllPlatforms               bool
	FilterLangs                []string
	FilterData                 []string
	FilterFilenames            []string
	AddNewlineEOF              bool
	CustomTranslationStatusIDs []string
	IncludeTags                []string
	ExcludeTags                []string
	ExportSort                 string
	ExportEmptyAs              string
	ExportNullAs               string
	IncludeComments            bool
	IncludeDescription         bool
	IncludePIDs                []string
	Triggers                   []string
	FilterRepositories         []string
	ReplaceBreaks              bool
	DisableReferences          bool
	PluralFormat               string
	PlaceholderFormat          string
	WebhookURL                 string
	LanguageMappingJSON        string
	ICUNumeric                 bool
	EscapePercent              bool
	Indentation                string
	YAMLIncludeRoot            bool
	JSONUnescapedSlashes       bool
	JavaPropertiesEncoding     string
	JavaPropertiesSeparator    string
	BundleDescription          string
	FilterTaskID               int64
	Compact                    bool
}

func newFlags() *Flags {
	return &Flags{
		Out:            "./locales",
		ContextTimeout: 150 * time.Second,
	}
}

func bindFlags(cmd *cobra.Command, flags *Flags) {
	cmd.Flags().SortFlags = false
	params.BindFlags(cmd, flags, downloadParamSpecs)
}
