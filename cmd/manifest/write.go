package manifest

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bodrovis/lokex-cli/internal/uploadmanifest"
	"github.com/spf13/cobra"
)

func writeGeneratedManifest(
	cmd *cobra.Command,
	out string,
	mf uploadmanifest.File,
) error {
	data, err := json.Marshal(
		mf,
		jsontext.WithIndent("  "),
	)
	if err != nil {
		return fmt.Errorf(
			"encode upload manifest: %w",
			err,
		)
	}

	data = append(data, '\n')

	if out == "-" {
		if _, err := cmd.OutOrStdout().Write(data); err != nil {
			return fmt.Errorf(
				"write upload manifest to stdout: %w",
				err,
			)
		}

		return nil
	}

	dir := filepath.Dir(out)

	if err := os.MkdirAll(
		dir,
		0o755,
	); err != nil {
		return fmt.Errorf(
			"create manifest directory %q: %w",
			dir,
			err,
		)
	}

	if err := os.WriteFile(
		out,
		data,
		0o644,
	); err != nil {
		return fmt.Errorf(
			"write upload manifest %q: %w",
			out,
			err,
		)
	}

	if _, err := fmt.Fprintf(
		cmd.OutOrStdout(),
		"Generated upload manifest: %s (%d items)\n",
		out,
		len(mf.Items),
	); err != nil {
		return fmt.Errorf(
			"write manifest result: %w",
			err,
		)
	}

	return nil
}
