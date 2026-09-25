package uploadmanifest

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

type DiscoveredFile struct {
	Path         string
	RelativePath string
}

func DiscoverFiles(
	paths []string,
) ([]DiscoveredFile, error) {
	var result []DiscoveredFile

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf(
				"inspect path %q: %w",
				path,
				err,
			)
		}

		if !info.IsDir() {
			result = append(
				result,
				DiscoveredFile{
					Path:         path,
					RelativePath: filepath.Base(path),
				},
			)

			continue
		}

		err = filepath.WalkDir(
			path,
			func(
				current string,
				entry fs.DirEntry,
				err error,
			) error {
				if err != nil {
					return err
				}

				if entry.IsDir() {
					return nil
				}

				relative, err := filepath.Rel(
					path,
					current,
				)
				if err != nil {
					return fmt.Errorf(
						"resolve path %q relative to %q: %w",
						current,
						path,
						err,
					)
				}

				result = append(
					result,
					DiscoveredFile{
						Path:         current,
						RelativePath: filepath.ToSlash(relative),
					},
				)

				return nil
			},
		)
		if err != nil {
			return nil, fmt.Errorf(
				"scan path %q: %w",
				path,
				err,
			)
		}
	}

	sort.Slice(
		result,
		func(i, j int) bool {
			return result[i].Path < result[j].Path
		},
	)

	return result, nil
}
