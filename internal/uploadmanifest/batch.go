package uploadmanifest

import (
	"fmt"
	"path/filepath"
	"strings"

	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
)

func BuildBatchUploadItems(
	manifestPath string,
	mf File,
) ([]lokexupload.BatchUploadItem, error) {
	manifestDir := filepath.Dir(manifestPath)

	items := make(
		[]lokexupload.BatchUploadItem,
		0,
		len(mf.Items),
	)

	for i, item := range mf.Items {
		if err := validateItem(item, i+1); err != nil {
			return nil, err
		}

		srcPath := item.SrcPath

		if srcPath != "" && !filepath.IsAbs(srcPath) {
			srcPath = filepath.Join(
				manifestDir,
				srcPath,
			)
		}

		items = append(
			items,
			lokexupload.BatchUploadItem{
				Params:  item.Params,
				SrcPath: srcPath,
			},
		)
	}

	return items, nil
}

func validateItem(
	item Item,
	index int,
) error {
	if item.Params == nil {
		return fmt.Errorf(
			"manifest item %d: params is required",
			index,
		)
	}

	if _, ok := requiredString(
		item.Params,
		"filename",
	); !ok {
		return fmt.Errorf(
			"manifest item %d: params.filename is required",
			index,
		)
	}

	if _, ok := requiredString(
		item.Params,
		"lang_iso",
	); !ok {
		return fmt.Errorf(
			"manifest item %d: params.lang_iso is required",
			index,
		)
	}

	return nil
}

func requiredString(
	params lokexupload.UploadParams,
	key string,
) (string, bool) {
	val, ok := params[key]
	if !ok {
		return "", false
	}

	s, ok := val.(string)
	if !ok {
		return "", false
	}

	s = strings.TrimSpace(s)

	if s == "" {
		return "", false
	}

	return s, true
}
