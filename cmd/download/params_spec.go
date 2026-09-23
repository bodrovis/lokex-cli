package download

import "github.com/bodrovis/lokex-cli/internal/params"

const configPrefix = "download"

var downloadParamSpecs = []params.Spec{
	{
		Name:    "out",
		Kind:    params.KindString,
		Usage:   "Directory to unzip downloaded bundle into",
		Default: "./locales",
	},
	{
		Name:  "format",
		Kind:  params.KindString,
		Usage: "File format (e.g. json, strings, xml)",
	},
	{
		Name:  "async",
		Kind:  params.KindBool,
		Usage: "Use Lokalise async download flow",
	},
	{
		Name:  "original-filenames",
		Kind:  params.KindBool,
		Usage: "Use original filenames/formats",
	},
	{
		Name:  "bundle-structure",
		Kind:  params.KindString,
		Usage: "Bundle structure when original-filenames=false",
	},
	{
		Name:  "directory-prefix",
		Kind:  params.KindString,
		Usage: "Directory prefix in bundle when original-filenames=true",
	},
	{
		Name:  "all-platforms",
		Kind:  params.KindBool,
		Usage: "Include all platform keys",
	},
	{
		Name:  "filter-langs",
		Kind:  params.KindStringSlice,
		Usage: "Languages to export",
	},
	{
		Name:  "filter-data",
		Kind:  params.KindStringSlice,
		Usage: "Narrow export data range",
	},
	{
		Name:  "filter-filenames",
		Kind:  params.KindStringSlice,
		Usage: "Only include keys attributed to selected files",
	},
	{
		Name:  "custom-translation-status-ids",
		Kind:  params.KindStringSlice,
		Usage: "Only include translations with selected custom status IDs",
	},
	{
		Name:  "include-tags",
		Kind:  params.KindStringSlice,
		Usage: "Only include keys with these tags",
	},
	{
		Name:  "exclude-tags",
		Kind:  params.KindStringSlice,
		Usage: "Exclude keys with these tags",
	},
	{
		Name:  "include-pids",
		Kind:  params.KindStringSlice,
		Usage: "Include keys from other project IDs",
	},
	{
		Name:  "triggers",
		Kind:  params.KindStringSlice,
		Usage: "Trigger integration exports",
	},
	{
		Name:  "filter-repositories",
		Kind:  params.KindStringSlice,
		Usage: "Only process selected repositories in organization/repository format",
	},
	{
		Name:  "filter-task-id",
		Kind:  params.KindInt64,
		Usage: "Only include keys attributed to this task (offline_xliff only)",
	},
	{
		Name:  "add-newline-eof",
		Kind:  params.KindBool,
		Usage: "Add newline at end of file when supported",
	},
	{
		Name:  "include-comments",
		Kind:  params.KindBool,
		Usage: "Include key comments and description when supported",
	},
	{
		Name:  "include-description",
		Kind:  params.KindBool,
		Usage: "Include key description when supported",
	},
	{
		Name:  "replace-breaks",
		Kind:  params.KindBool,
		Usage: "Replace line breaks in exported translations with \\n",
	},
	{
		Name:  "disable-references",
		Kind:  params.KindBool,
		Usage: "Disable automatic replacement of key reference placeholders",
	},
	{
		Name:  "icu-numeric",
		Kind:  params.KindBool,
		Usage: "Replace ICU plural forms zero/one/two with =0/=1/=2",
	},
	{
		Name:  "escape-percent",
		Kind:  params.KindBool,
		Usage: "Escape universal percent placeholders for printf format",
	},
	{
		Name:  "yaml-include-root",
		Kind:  params.KindBool,
		Usage: "Include language ISO code as root key for YAML export",
	},
	{
		Name:  "json-unescaped-slashes",
		Kind:  params.KindBool,
		Usage: "Leave forward slashes unescaped in JSON export",
	},
	{
		Name:  "compact",
		Kind:  params.KindBool,
		Usage: "Export compact ARB structure",
	},
	{
		Name:  "export-sort",
		Kind:  params.KindString,
		Usage: "Export key sort mode",
	},
	{
		Name:  "export-empty-as",
		Kind:  params.KindString,
		Usage: "How to export empty translations",
	},
	{
		Name:  "export-null-as",
		Kind:  params.KindString,
		Usage: "How to export null translations (Ruby on Rails YAML only)",
	},
	{
		Name:  "plural-format",
		Kind:  params.KindString,
		Usage: "Override default plural format",
	},
	{
		Name:  "placeholder-format",
		Kind:  params.KindString,
		Usage: "Override default placeholder format",
	},
	{
		Name:  "webhook-url",
		Kind:  params.KindString,
		Usage: "Send POST with generated bundle URL to this URL when export completes",
	},
	{
		Name:  "indentation",
		Kind:  params.KindString,
		Usage: "Override indentation in supported files",
	},
	{
		Name:  "java-properties-encoding",
		Kind:  params.KindString,
		Usage: "Encoding for Java .properties export",
	},
	{
		Name:  "java-properties-separator",
		Kind:  params.KindString,
		Usage: "Separator for Java .properties export",
	},
	{
		Name:  "bundle-description",
		Kind:  params.KindString,
		Usage: "Description for ios_sdk/android_sdk OTA SDK bundles",
	},
	{
		Name:  "language-mapping",
		Kind:  params.KindString,
		Usage: "Language mapping as JSON array of objects",
	},
}
