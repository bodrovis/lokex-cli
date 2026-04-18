package download

import (
	"github.com/spf13/cobra"

	lokexdownload "github.com/bodrovis/lokex/v2/client/download"

	params "github.com/bodrovis/lokex-cli/internal/params"
)

func buildParams(cmd *cobra.Command, flags *Flags, defaults *DownloadConfig) (lokexdownload.DownloadParams, error) {
	req := lokexdownload.DownloadParams{}
	if err := params.ApplyToRequest(cmd, flags, defaults, req, downloadParamSpecs); err != nil {
		return nil, err
	}
	return req, nil
}
