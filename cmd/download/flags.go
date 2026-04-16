package download

import (
	"time"

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
	cmd.Flags().StringVar(
		&flags.Out,
		"out",
		flags.Out,
		"Directory to unzip downloaded bundle into",
	)

	cmd.Flags().StringVar(
		&flags.Format,
		"format",
		"",
		"File format (e.g. json, strings, xml)",
	)

	cmd.Flags().DurationVar(
		&flags.ContextTimeout,
		"context-timeout",
		flags.ContextTimeout,
		"Overall command timeout (e.g. 30s, 2m). 0 disables the timeout",
	)

	cmd.Flags().BoolVar(
		&flags.Async,
		"async",
		false,
		"Use Lokalise async download flow",
	)

	cmd.Flags().BoolVar(
		&flags.OriginalFilenames,
		"original-filenames",
		false,
		"Use original filenames/formats",
	)

	cmd.Flags().StringVar(
		&flags.BundleStructure,
		"bundle-structure",
		"",
		"Bundle structure when original-filenames=false",
	)

	cmd.Flags().StringVar(
		&flags.DirectoryPrefix,
		"directory-prefix",
		"",
		"Directory prefix in bundle when original-filenames=true",
	)

	cmd.Flags().BoolVar(
		&flags.AllPlatforms,
		"all-platforms",
		false,
		"Include all platform keys",
	)

	cmd.Flags().StringSliceVar(
		&flags.FilterLangs,
		"filter-langs",
		nil,
		"Languages to export",
	)

	cmd.Flags().StringSliceVar(
		&flags.FilterData,
		"filter-data",
		nil,
		"Narrow export data range",
	)

	cmd.Flags().StringSliceVar(
		&flags.FilterFilenames,
		"filter-filenames",
		nil,
		"Only include keys attributed to selected files",
	)

	cmd.Flags().BoolVar(
		&flags.AddNewlineEOF,
		"add-newline-eof",
		false,
		"Add newline at end of file when supported",
	)

	cmd.Flags().StringSliceVar(
		&flags.CustomTranslationStatusIDs,
		"custom-translation-status-ids",
		nil,
		"Only include translations with selected custom status IDs",
	)

	cmd.Flags().StringSliceVar(
		&flags.IncludeTags,
		"include-tags",
		nil,
		"Only include keys with these tags",
	)

	cmd.Flags().StringSliceVar(
		&flags.ExcludeTags,
		"exclude-tags",
		nil,
		"Exclude keys with these tags",
	)

	cmd.Flags().StringVar(
		&flags.ExportSort,
		"export-sort",
		"",
		"Export key sort mode",
	)

	cmd.Flags().StringVar(
		&flags.ExportEmptyAs,
		"export-empty-as",
		"",
		"How to export empty translations",
	)

	cmd.Flags().StringVar(
		&flags.ExportNullAs,
		"export-null-as",
		"",
		"How to export null translations (Ruby on Rails YAML only)",
	)

	cmd.Flags().BoolVar(
		&flags.IncludeComments,
		"include-comments",
		false,
		"Include key comments and description when supported",
	)

	cmd.Flags().BoolVar(
		&flags.IncludeDescription,
		"include-description",
		false,
		"Include key description when supported",
	)

	cmd.Flags().StringSliceVar(
		&flags.IncludePIDs,
		"include-pids",
		nil,
		"Include keys from other project IDs",
	)

	cmd.Flags().StringSliceVar(
		&flags.Triggers,
		"triggers",
		nil,
		"Trigger integration exports",
	)

	cmd.Flags().StringSliceVar(
		&flags.FilterRepositories,
		"filter-repositories",
		nil,
		"Only process selected repositories in organization/repository format",
	)

	cmd.Flags().BoolVar(
		&flags.ReplaceBreaks,
		"replace-breaks",
		false,
		"Replace line breaks in exported translations with \\n",
	)

	cmd.Flags().BoolVar(
		&flags.DisableReferences,
		"disable-references",
		false,
		"Disable automatic replacement of key reference placeholders",
	)

	cmd.Flags().StringVar(
		&flags.PluralFormat,
		"plural-format",
		"",
		"Override default plural format",
	)

	cmd.Flags().StringVar(
		&flags.PlaceholderFormat,
		"placeholder-format",
		"",
		"Override default placeholder format",
	)

	cmd.Flags().StringVar(
		&flags.WebhookURL,
		"webhook-url",
		"",
		"Send POST with generated bundle URL to this URL when export completes",
	)

	cmd.Flags().StringVar(
		&flags.LanguageMappingJSON,
		"language-mapping",
		"",
		"Language mapping as JSON array of objects",
	)

	cmd.Flags().BoolVar(
		&flags.ICUNumeric,
		"icu-numeric",
		false,
		"Replace ICU plural forms zero/one/two with =0/=1/=2",
	)

	cmd.Flags().BoolVar(
		&flags.EscapePercent,
		"escape-percent",
		false,
		"Escape universal percent placeholders for printf format",
	)

	cmd.Flags().StringVar(
		&flags.Indentation,
		"indentation",
		"",
		"Override indentation in supported files",
	)

	cmd.Flags().BoolVar(
		&flags.YAMLIncludeRoot,
		"yaml-include-root",
		false,
		"Include language ISO code as root key for YAML export",
	)

	cmd.Flags().BoolVar(
		&flags.JSONUnescapedSlashes,
		"json-unescaped-slashes",
		false,
		"Leave forward slashes unescaped in JSON export",
	)

	cmd.Flags().StringVar(
		&flags.JavaPropertiesEncoding,
		"java-properties-encoding",
		"",
		"Encoding for Java .properties export",
	)

	cmd.Flags().StringVar(
		&flags.JavaPropertiesSeparator,
		"java-properties-separator",
		"",
		"Separator for Java .properties export",
	)

	cmd.Flags().StringVar(
		&flags.BundleDescription,
		"bundle-description",
		"",
		"Description for ios_sdk/android_sdk OTA SDK bundles",
	)

	cmd.Flags().Int64Var(
		&flags.FilterTaskID,
		"filter-task-id",
		0,
		"Only include keys attributed to this task (offline_xliff only)",
	)

	cmd.Flags().BoolVar(
		&flags.Compact,
		"compact",
		false,
		"Export compact ARB structure",
	)
}
