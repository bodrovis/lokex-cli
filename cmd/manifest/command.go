package manifest

import (
	"fmt"
	"path/filepath"

	"github.com/bodrovis/lokex-cli/internal/appstate"
	"github.com/bodrovis/lokex-cli/internal/languagemapping"
	"github.com/bodrovis/lokex-cli/internal/ptrutil"
	"github.com/bodrovis/lokex-cli/internal/uploadmanifest"
	"github.com/spf13/cobra"
)

func NewCommand(
	state *appstate.State,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manifest",
		Short: "Work with upload manifests",
		Long: `Create and work with manifests for batch uploads.

Upload manifests describe local source files and the upload parameters
that should be used for each file.`,
		Args: cobra.NoArgs,
	}

	cmd.AddCommand(
		newGenerateCommand(state),
	)

	return cmd
}

func runGenerateCommand(
	cmd *cobra.Command,
	cfg *GenerateConfig,
) error {
	out := ptrutil.TrimmedString(cfg.Out)

	mapping, err := languagemapping.Parse(
		ptrutil.TrimmedString(cfg.LanguageMapping),
	)
	if err != nil {
		return fmt.Errorf(
			"parse language-mapping: %w",
			err,
		)
	}

	manifestDir := "."

	if out != "-" {
		manifestDir = filepath.Dir(out)
	}

	mf, err := uploadmanifest.Generate(
		uploadmanifest.GenerateOptions{
			Paths: ptrutil.Value(
				cfg.Paths,
			),
			ExcludePatterns: ptrutil.Value(
				cfg.ExcludePatterns,
			),
			NamePattern: ptrutil.TrimmedString(
				cfg.NamePattern,
			),
			FilenamePattern: ptrutil.TrimmedString(
				cfg.FilenamePattern,
			),
			BaseLang: ptrutil.TrimmedString(
				cfg.BaseLang,
			),
			LanguageMapping: mapping,
			ManifestDir:     manifestDir,
		},
	)
	if err != nil {
		return err
	}

	return writeGeneratedManifest(
		cmd,
		out,
		mf,
	)
}
