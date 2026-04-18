package upload

import (
	params "github.com/bodrovis/lokex-cli/internal/params"
	"github.com/spf13/cobra"
)

func applyDefaults(cmd *cobra.Command, flags *Flags, cfg *UploadConfig) {
	params.ApplyDefaults(cmd, flags, cfg, uploadParamSpecs)
}
