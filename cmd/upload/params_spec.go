package upload

import "github.com/bodrovis/lokex-cli/internal/params"

const configPrefix = "upload"

var uploadParamSpecs = []params.Spec{
	{
		Name:  "manifest",
		Kind:  params.KindString,
		Usage: "Path to JSON manifest for batch upload",
	},
	{
		Name:  "filename",
		Kind:  params.KindString,
		Usage: "Filename sent to Lokalise (required)",
	},
	{
		Name:  "src-path",
		Kind:  params.KindString,
		Usage: "Local path to read file contents from (optional)",
	},
	{
		Name:  "data",
		Kind:  params.KindString,
		Usage: "Base64-encoded file contents (optional; if set, file is not read from disk)",
	},
	{
		Name:  "lang-iso",
		Kind:  params.KindString,
		Usage: "Language code of the translations in the file (required)",
	},
	{
		Name:  "poll",
		Kind:  params.KindBool,
		Usage: "Wait until Lokalise finishes processing the upload",
	},
	{
		Name:  "format",
		Kind:  params.KindString,
		Usage: "File format (e.g. json, strings, xml)",
	},
	{
		Name:  "tags",
		Kind:  params.KindStringSlice,
		Usage: "Tags to apply to keys",
	},
	{
		Name:  "custom-translation-status-ids",
		Kind:  params.KindStringSlice,
		Usage: "Custom translation status IDs to add",
	},
	{
		Name:  "convert-placeholders",
		Kind:  params.KindBool,
		Usage: "Convert placeholders to Lokalise universal placeholders",
	},
	{
		Name:  "detect-icu-plurals",
		Kind:  params.KindBool,
		Usage: "Automatically detect and parse ICU plurals",
	},
	{
		Name:  "tag-inserted-keys",
		Kind:  params.KindBool,
		Usage: "Add tags to inserted keys",
	},
	{
		Name:  "tag-updated-keys",
		Kind:  params.KindBool,
		Usage: "Add tags to updated keys",
	},
	{
		Name:  "tag-skipped-keys",
		Kind:  params.KindBool,
		Usage: "Add tags to skipped keys",
	},
	{
		Name:  "replace-modified",
		Kind:  params.KindBool,
		Usage: "Replace modified translations from the uploaded file",
	},
	{
		Name:  "slashn-to-linebreak",
		Kind:  params.KindBool,
		Usage: "Replace \\n with a real line break",
	},
	{
		Name:  "keys-to-values",
		Kind:  params.KindBool,
		Usage: "Replace values with key names",
	},
	{
		Name:  "distinguish-by-file",
		Kind:  params.KindBool,
		Usage: "Allow same key names to coexist across different filenames",
	},
	{
		Name:  "apply-tm",
		Kind:  params.KindBool,
		Usage: "Apply 100% translation memory matches",
	},
	{
		Name:  "use-automations",
		Kind:  params.KindBool,
		Usage: "Run automations for this upload",
	},
	{
		Name:  "hidden-from-contributors",
		Kind:  params.KindBool,
		Usage: "Mark newly created keys as hidden from contributors",
	},
	{
		Name:  "cleanup-mode",
		Kind:  params.KindBool,
		Usage: "Delete keys/translations not present in the uploaded file",
	},
	{
		Name:  "custom-translation-status-inserted-keys",
		Kind:  params.KindBool,
		Usage: "Add custom statuses to inserted keys",
	},
	{
		Name:  "custom-translation-status-updated-keys",
		Kind:  params.KindBool,
		Usage: "Add custom statuses to updated keys",
	},
	{
		Name:  "custom-translation-status-skipped-keys",
		Kind:  params.KindBool,
		Usage: "Add custom statuses to skipped keys",
	},
	{
		Name:  "skip-detect-lang-iso",
		Kind:  params.KindBool,
		Usage: "Skip automatic language detection by filename",
	},
	{
		Name:  "filter-task-id",
		Kind:  params.KindInt64,
		Usage: "Apply import results as a part of a task (offline_xliff only)",
	},
}
