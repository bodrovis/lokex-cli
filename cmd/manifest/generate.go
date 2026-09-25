package manifest

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bodrovis/lokex-cli/internal/appstate"
	"github.com/bodrovis/lokex-cli/internal/params"
	"github.com/bodrovis/lokex-cli/internal/ptrutil"
	"github.com/bodrovis/lokex-cli/internal/uploadmanifest"
	"github.com/spf13/cobra"
)

func newGenerateCommand(
	state *appstate.State,
) *cobra.Command {
	cfg := &GenerateConfig{}

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate an upload manifest from local files",
		Long: `Discover local translation files and generate an upload manifest.

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
it to a file.`,
		Example: `  # Files named by language: en.json, de.json, fr.json
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
    --out -`,
		Args: cobra.NoArgs,
		PreRunE: func(
			cmd *cobra.Command,
			_ []string,
		) error {
			if state == nil {
				return errors.New("app state is nil")
			}

			if state.Viper == nil {
				return errors.New("viper is nil")
			}

			if err := LoadGenerateConfig(
				state.Viper,
				cmd,
				cfg,
			); err != nil {
				return err
			}

			return validateGenerateConfig(cfg)
		},
		RunE: func(
			cmd *cobra.Command,
			_ []string,
		) error {
			return runGenerateCommand(
				cmd,
				cfg,
			)
		},
	}

	params.BindFlags(
		cmd,
		generateParamSpecs,
	)

	return cmd
}

func validateGenerateConfig(
	cfg *GenerateConfig,
) error {
	if cfg == nil {
		return errors.New("manifest generate config is nil")
	}

	if cfg.Paths == nil || len(*cfg.Paths) == 0 {
		return errors.New("path is required")
	}

	hasPath := false

	for _, path := range *cfg.Paths {
		if strings.TrimSpace(path) != "" {
			hasPath = true
			break
		}
	}

	if !hasPath {
		return errors.New("path is required")
	}

	namePattern := ptrutil.TrimmedString(
		cfg.NamePattern,
	)

	if namePattern == "" {
		return errors.New("name-pattern is required")
	}

	if err := uploadmanifest.ValidatePathPattern(
		namePattern,
	); err != nil {
		return fmt.Errorf(
			"invalid name-pattern: %w",
			err,
		)
	}

	filenamePattern := ptrutil.TrimmedString(
		cfg.FilenamePattern,
	)

	if filenamePattern == "" {
		return errors.New("filename-pattern is required")
	}

	if err := uploadmanifest.ValidateRenderPattern(
		filenamePattern,
	); err != nil {
		return fmt.Errorf(
			"invalid filename-pattern: %w",
			err,
		)
	}

	baseLang := ptrutil.TrimmedString(
		cfg.BaseLang,
	)

	if !strings.Contains(
		namePattern,
		"{lang}",
	) && baseLang == "" {
		return errors.New(
			"base-lang is required when name-pattern does not contain {lang}",
		)
	}

	if err := uploadmanifest.ValidatePatternSources(
		namePattern,
		filenamePattern,
		baseLang,
	); err != nil {
		return fmt.Errorf(
			"invalid filename-pattern: %w",
			err,
		)
	}

	if ptrutil.TrimmedString(cfg.Out) == "" {
		return errors.New("out is required")
	}

	return nil
}
