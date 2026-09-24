package upload

import (
	"encoding/json/v2"
	"fmt"
	"os"

	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
)

type manifestFile struct {
	Items []manifestItem `json:"items"`
}

type manifestItem struct {
	Params  lokexupload.UploadParams `json:"params"`
	SrcPath string                   `json:"src_path"`
}

func loadManifestFile(path string) (manifestFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return manifestFile{}, fmt.Errorf(
			"read manifest file %q: %w",
			path,
			err,
		)
	}
	defer func() { _ = f.Close() }()

	var mf manifestFile

	if err := json.UnmarshalRead(
		f,
		&mf,
		preserveJSONNumbers,
	); err != nil {
		return manifestFile{}, fmt.Errorf(
			"parse manifest file %q: %w",
			path,
			err,
		)
	}

	if len(mf.Items) == 0 {
		return manifestFile{}, fmt.Errorf(
			"manifest file %q contains no items",
			path,
		)
	}

	return mf, nil
}
