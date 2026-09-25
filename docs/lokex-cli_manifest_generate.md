## lokex-cli manifest generate

Generate an upload manifest from local files

### Synopsis

Discover local translation files and generate an upload manifest.

The command scans the paths provided with --path, filters excluded files,
matches each remaining file against --name-pattern, and generates an upload
manifest that can later be passed to "lokex upload --manifest".

The name pattern describes the structure of local file paths.

Patterns support:
  *         Any characters within a single path segment
  **        Any number of directories
  {lang}    Language code
  {name}    File name
  {ext}     File extension

You can also use custom placeholders such as {app}, {module}, or {namespace}.
Values captured by placeholders can be reused in --filename-pattern.

For example:

  --path ./apps
  --name-pattern "{app}/locales/{lang}/{name}.{ext}"
  --filename-pattern "{app}/{name}.{ext}"

matches:

  apps/mobile/locales/en/common.json

and generates:

  filename = mobile/common.json
  lang_iso = en

Use ** when the directory depth is not fixed. For example:

  --name-pattern "**/locales/{lang}.{ext}"

matches files such as:

  apps/mobile/locales/en.json
  packages/shared/locales/de.json
  locales/fr.json

Neither {name} nor {ext} is required. For files named only by language,
for example en.json and de.json, use:

  --name-pattern "{lang}.{ext}"

If the name pattern does not contain {lang}, use --base-lang to assign a
language to matching files.

The filename pattern controls the filename stored in Lokalise. Its default
value is "{name}.{ext}". Any placeholder used by the filename pattern must
be provided by the name pattern, except {lang}, which can also come from
--base-lang.

Generated source paths are relative to the manifest file location when
possible. If the source and manifest are on different filesystem volumes,
an absolute source path is stored instead.

Use --out=- to print the generated manifest to stdout instead of writing
it to a file.

```
lokex-cli manifest generate [flags]
```

### Examples

```
  # Files named by language: en.json, de.json, fr.json
  lokex manifest generate \
    --path ./locales \
    --name-pattern "{lang}.{ext}" \
    --filename-pattern "{lang}.{ext}"

  # Language directories: en/common.json, de/common.json
  lokex manifest generate \
    --path ./locales \
    --name-pattern "{lang}/{name}.{ext}" \
    --filename-pattern "{name}.{ext}"

  # Match locales directories at any depth
  lokex manifest generate \
    --path . \
    --name-pattern "**/locales/{lang}.{ext}" \
    --filename-pattern "{lang}.{ext}"

  # Preserve application names in a monorepo
  lokex manifest generate \
    --path ./apps \
    --name-pattern "{app}/locales/{lang}.{ext}" \
    --filename-pattern "{app}/{lang}.{ext}"

  # Use one language when it is not present in the path
  lokex manifest generate \
    --path ./locales \
    --name-pattern "{name}.{ext}" \
    --base-lang en

  # Exclude test files
  lokex manifest generate \
    --path ./locales \
    --name-pattern "{lang}/{name}.{ext}" \
    --exclude-pattern "**/*.test.json"

  # Print the manifest to stdout
  lokex manifest generate \
    --path ./locales \
    --name-pattern "{lang}/{name}.{ext}" \
    --out -
```

### Options

```
      --path strings              File or directory to scan for translation files; can be specified multiple times
      --exclude-pattern strings   Glob pattern for files to exclude; can be specified multiple times
      --name-pattern string       Pattern used to extract {lang}, {name}, and {ext} from local paths
      --filename-pattern string   Pattern used to generate filenames stored in the upload manifest (default "{name}.{ext}")
      --base-lang string          Language code to use when --name-pattern does not contain {lang}
      --language-mapping string   JSON language mapping applied to generated lang_iso values
      --out string                Output manifest path; use - to write to stdout (default "./lokex-manifest.json")
  -h, --help                      help for generate
```

### Options inherited from parent commands

```
      --backoff-initial duration     Initial retry backoff (e.g. 400ms, 1s). 0 means library default
      --backoff-max duration         Maximum retry backoff (e.g. 5s, 10s). 0 means library default
      --base-url string              Override Lokalise API base URL
      --config string                Path to YAML config file
      --context-timeout duration     Overall command timeout (e.g. 30s, 2m). 0 disables the timeout (default 2m30s)
      --http-timeout duration        HTTP client timeout (e.g. 30s, 1m). 0 means library default
      --poll-initial-wait duration   Initial wait between polling rounds (e.g. 1s, 2s). 0 means library default
      --poll-max-wait duration       Maximum total wait for polling (e.g. 120s, 5m). 0 means library default
      --project-id string            Lokalise project ID
      --retries int                  Number of retries after the first attempt. -1 means library default (default -1)
      --token string                 Lokalise API token
      --user-agent string            User-Agent header (default "lokex-cli/dev")
```

### SEE ALSO

* [lokex-cli manifest](lokex-cli_manifest.md)	 - Work with upload manifests

###### Auto generated by spf13/cobra on 25-Sep-2026
