# Lokex CLI: Content exchange tool for Lokalise

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

## License

(c) [Elijah S. Krukowski](https://bodrovis.tech). Licensed under BSD 3-Clause