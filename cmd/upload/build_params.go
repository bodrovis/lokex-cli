package upload

import (
	"github.com/spf13/cobra"

	params "github.com/bodrovis/lokex-cli/internal/params"
	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
)

func buildParams(cmd *cobra.Command, flags *Flags, defaults *UploadConfig) (lokexupload.UploadParams, error) {
	req := lokexupload.UploadParams{}
	if err := params.ApplyToRequest(cmd, flags, defaults, req, uploadParamSpecs); err != nil {
		return nil, err
	}
	return req, nil
}
