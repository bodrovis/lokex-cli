package manifest

import "github.com/bodrovis/lokex-cli/internal/params"

const generateConfigPrefix = "manifest.generate"

var generateParamSpecs = []params.Spec{
	{
		Name:  "path",
		Kind:  params.KindStringSlice,
		Usage: "File or directory to scan for translation files; can be specified multiple times",
	},
	{
		Name:  "exclude-pattern",
		Kind:  params.KindStringSlice,
		Usage: "Glob pattern for files to exclude; can be specified multiple times",
	},
	{
		Name:  "name-pattern",
		Kind:  params.KindString,
		Usage: "Pattern used to extract {lang}, {name}, and {ext} from local paths",
	},
	{
		Name:    "filename-pattern",
		Kind:    params.KindString,
		Usage:   "Pattern used to generate filenames stored in the upload manifest",
		Default: "{name}.{ext}",
	},
	{
		Name:  "base-lang",
		Kind:  params.KindString,
		Usage: "Language code to use when --name-pattern does not contain {lang}",
	},
	{
		Name:  "language-mapping",
		Kind:  params.KindString,
		Usage: "JSON language mapping applied to generated lang_iso values",
	},
	{
		Name:    "out",
		Kind:    params.KindString,
		Usage:   "Output manifest path; use - to write to stdout",
		Default: "./lokex-manifest.json",
	},
}
