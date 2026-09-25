package uploadmanifest

import (
	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
)

type File struct {
	Items []Item `json:"items"`
}

type Item struct {
	Params  lokexupload.UploadParams `json:"params"`
	SrcPath string                   `json:"src_path"`
}
