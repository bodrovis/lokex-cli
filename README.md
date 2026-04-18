# Lokex CLI: Content exchange tool for Lokalise

![GitHub Release](https://img.shields.io/github/v/release/bodrovis/lokex-cli)
![CI](https://github.com/bodrovis/lokex-cli/actions/workflows/ci.yml/badge.svg)
[![Code Coverage](https://qlty.sh/gh/bodrovis/projects/lokex-cli/coverage.svg)](https://qlty.sh/gh/bodrovis/projects/lokex-cli)
[![Maintainability](https://qlty.sh/gh/bodrovis/projects/lokex-cli/maintainability.svg)](https://qlty.sh/gh/bodrovis/projects/lokex-cli)

`lokex-cli` is a focused CLI built specifically for **file exchange with Lokalise** on top of [`lokex`](https://github.com/bodrovis/lokex).

It is intentionally narrow in scope, so you can only upload and download files. This tool is meant to be a fast, optimized workflow for import/export operations.

If you need a broader Lokalise command set, use the [official CLI instead](https://github.com/lokalise/lokalise-cli-2-go). 

## Documentation

Detailed command docs are generated in:

- [`docs/lokex-cli_upload.md`](https://github.com/bodrovis/lokex-cli/blob/master/docs/lokex-cli_upload.md)
- [`docs/lokex-cli_download.md`](https://github.com/bodrovis/lokex-cli/blob/master/docs/lokex-cli_download.md)

Those files list all supported flags and API-related parameters for each command.

## Global flags

These flags are shared by both `upload` and `download`.

> [Find full list of global flags in the docs](https://github.com/bodrovis/lokex-cli/blob/master/docs/lokex-cli.md).

### Required in all cases

These two flags are always required:

- `--token`
- `--project-id`

Example:

```bash
lokex-cli download \
  --token YOUR_TOKEN \
  --project-id YOUR_PROJECT_ID \
  --format json
```

## Download examples

### Basic download

```
lokex-cli download \
  --token YOUR_TOKEN \
  --project-id YOUR_PROJECT_ID \
  --format json
```

This downloads files from Lokalise and extracts them into `./locales` by default.

### Download into a custom directory

```
lokex-cli download \
  --token YOUR_TOKEN \
  --project-id YOUR_PROJECT_ID \
  --format json \
  --out ./tmp/locales
```

### Async download flow

```
lokex-cli download \
  --token YOUR_TOKEN \
  --project-id YOUR_PROJECT_ID \
  --format json \
  --async
```

### Download with extra API parameters

```
lokex-cli download \
  --token YOUR_TOKEN \
  --project-id YOUR_PROJECT_ID \
  --format json \
  --original-filenames \
  --filter-langs en,fr,de \
  --include-tags mobile,release
```

## Upload examples

### Basic upload

```
lokex-cli upload \
  --token YOUR_TOKEN \
  --project-id YOUR_PROJECT_ID \
  --filename locales/en.json \
  --lang-iso en
```

### Upload and wait until processing finishes

```
lokex-cli upload \
  --token YOUR_TOKEN \
  --project-id YOUR_PROJECT_ID \
  --filename en.json \
  --lang-iso en \
  --poll
```

### Upload with extra import options

```
lokex-cli upload \
  --token YOUR_TOKEN \
  --project-id YOUR_PROJECT_ID \
  --filename en.json \
  --lang-iso en \
  --replace-modified \
  --convert-placeholders \
  --tags backend,release
```

## `filename` and `src-path`

For most uploads, you only need `--filename`.

Example:

```bash
lokex-cli upload \
  --token YOUR_TOKEN \
  --project-id YOUR_PROJECT_ID \
  --filename ./locales/en.json \
  --lang-iso en
```

That is totally fine.

If needed, you can also set `--src-path` explicitly to tell the CLI where to read the local file from, while keeping a different filename for Lokalise:

```
lokex-cli upload \
  --token YOUR_TOKEN \
  --project-id YOUR_PROJECT_ID \
  --filename en.json \
  --src-path ./locales/en.json \
  --lang-iso en
```

### About data

The upload API supports a `data` field.

You can pass it with:

```
--data BASE64_ENCODED_CONTENT
```

But in normal usage you usually do not need this.

If `data` is not provided, lokex-cli reads the file from disk and prepares the payload for you automatically.

So the normal upload flow is:

- provide `--filename`
- optionally provide `--src-path`
- provide `--lang-iso`
- let the tool read and encode the file itself

## Passing arrays and JSON via CLI

Some flags accept arrays (lists of values) or structured JSON. Here is how to pass them in a shell.

### String arrays

Flags like `--filter-langs`, `--include-tags`, etc. accept multiple values.

You can pass them in two ways:

#### Comma-separated

```
lokex-cli download \
  --format=json \
  --filter-langs=en,fr,de \
  --include-tags=mobile,web
```

#### Repeated flags

```
lokex-cli download \
  --format=json \
  --filter-langs=en \
  --filter-langs=fr \
  --filter-langs=de
```

### JSON arrays

Some flags like `--language-mapping` expect structured JSON (array of objects).

Example:

```
lokex-cli download \
  --format=json \
  --language-mapping='[
    {"lang_iso":"en","custom_iso":"en_US"},
    {"lang_iso":"pt","custom_iso":"pt_BR"}
  ]'
```

**Important notes**:

- Always wrap JSON in quotes ('...') to prevent the shell from interpreting it.
- On Windows (PowerShell), you may need to use double quotes:

```
--language-mapping "[{\"lang_iso\":\"en\",\"custom_iso\":\"en_US\"}]"
```

## Configuration via environment variables and YAML

In addition to CLI flags, `lokex-cli` can also read command defaults from:

- environment variables
- an optional YAML config file

### YAML file

Config file can be provided explicitly:

```
lokex-cli --config path/to/config.yaml <command>
```

If `--config` is not set, `lokex-cli` will try to automatically find a config file named `lokex.yaml` in the following locations:

- current working directory (`./lokex.yaml`)
- user config directory: `~/.config/lokex-cli/lokex.yaml`

If no config file is found, execution continues without error.

Command-specific options must be placed under the matching namespace:

- `download.*` for `lokex-cli download`
- `upload.*` for `lokex-cli upload`

Example:

```yaml
token: token
project-id: project-id

download:
  format: json
  out: ./tmp/locales
  async: true
  filter-langs:
    - en
    - fr
  include-tags:
    - mobile
    - release

upload:
  filename: en.json
  src-path: ./locales/en.json
  lang-iso: en
  poll: true
  replace-modified: true
  tags:
    - backend
    - release
```

### Environment variables

Environment variables follow the same logical keys as the YAML config, using the `LOKEX` env prefix.

```
LOKEX_TOKEN=token
LOKEX_PROJECT_ID=project-id
LOKEX_DOWNLOAD_FORMAT=json
LOKEX_DOWNLOAD_OUT=./locales
LOKEX_UPLOAD_FILENAME=en.json
LOKEX_UPLOAD_LANG_ISO=en
LOKEX_UPLOAD_POLL=true
```

### Precedence

If you pass a value explicitly as a CLI flag, it overrides the value from YAML or environment variables.

So the effective order is:

- explicit CLI flag
- environment variable / YAML config
- built-in default

## Testing

Run:

```
go test -count=1 ./... -shuffle=on -race
```

## License

(c) [Elijah S. Krukowski](https://bodrovis.tech). Licensed under BSD 3-Clause
