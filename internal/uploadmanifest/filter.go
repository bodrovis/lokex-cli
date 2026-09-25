package uploadmanifest

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

func FilterFiles(
	files []DiscoveredFile,
	excludePatterns []string,
) ([]DiscoveredFile, error) {
	patterns, err := normalizeExcludePatterns(
		excludePatterns,
	)
	if err != nil {
		return nil, err
	}

	if len(patterns) == 0 {
		return files, nil
	}

	result := make(
		[]DiscoveredFile,
		0,
		len(files),
	)

	for _, file := range files {
		path := filepath.ToSlash(
			file.RelativePath,
		)

		excluded, err := matchesAnyPattern(
			path,
			patterns,
		)
		if err != nil {
			return nil, err
		}

		if !excluded {
			result = append(
				result,
				file,
			)
		}
	}

	return result, nil
}

func normalizeExcludePatterns(
	patterns []string,
) ([]string, error) {
	result := make(
		[]string,
		0,
		len(patterns),
	)

	for _, pattern := range patterns {
		pattern = filepath.ToSlash(
			strings.TrimSpace(pattern),
		)

		if pattern == "" {
			continue
		}

		if _, err := doublestar.Match(
			pattern,
			"",
		); err != nil {
			return nil, fmt.Errorf(
				"invalid exclude pattern %q: %w",
				pattern,
				err,
			)
		}

		result = append(
			result,
			pattern,
		)
	}

	return result, nil
}

func matchesAnyPattern(
	path string,
	patterns []string,
) (bool, error) {
	for _, pattern := range patterns {
		matched, err := doublestar.Match(
			pattern,
			path,
		)
		if err != nil {
			return false, fmt.Errorf(
				"match exclude pattern %q: %w",
				pattern,
				err,
			)
		}

		if matched {
			return true, nil
		}
	}

	return false, nil
}
