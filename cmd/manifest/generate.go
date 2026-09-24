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

The name pattern describes the structure of local file paths and supports
the following placeholders:

  {lang}  Language code detected from the path
  {name}  File name without the extension
  {ext}   File extension

For example:

  --path ./locales
  --name-pattern "{lang}/{name}.{ext}"

matches files such as:

  locales/en/common.json
  locales/de/common.json

The filename pattern controls the filename stored in the upload manifest.
Its default value is "{name}.{ext}".

Generated source paths are relative to the manifest file location,
so the manifest can be used from any working directory.

If the name pattern does not contain {lang}, use --base-lang to assign
a language to matching files.

Use --out=- to print the generated manifest to stdout instead of writing
it to a file.`,
		Example: `  # Generate a manifest from a nested language directory
  lokex manifest generate \
    --path ./locales \
    --name-pattern "{lang}/{name}.{ext}"

  # Use one language for files without a language in their names
  lokex manifest generate \
    --path ./locales \
    --name-pattern "{name}.{ext}" \
    --base-lang en

  # Exclude generated or test files
  lokex manifest generate \
    --path ./locales \
    --name-pattern "{lang}/{name}.{ext}" \
    --exclude-pattern "**/*.test.json"

  # Generate a manifest for flat filenames
  lokex manifest generate \
    --path ./locales \
    --name-pattern "{name}.{lang}.{ext}"

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

	if !strings.Contains(namePattern, "{lang}") &&
		ptrutil.TrimmedString(cfg.BaseLang) == "" {
		return errors.New(
			"base-lang is required when name-pattern does not contain {lang}",
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

	if ptrutil.TrimmedString(cfg.Out) == "" {
		return errors.New("out is required")
	}

	return nil
}
