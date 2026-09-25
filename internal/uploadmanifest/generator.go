package uploadmanifest

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bodrovis/lokex-cli/internal/languagemapping"
	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
)

type GenerateOptions struct {
	Paths           []string
	ExcludePatterns []string

	NamePattern     string
	FilenamePattern string

	BaseLang        string
	LanguageMapping []languagemapping.Mapping

	ManifestDir string
}

func Generate(
	opts GenerateOptions,
) (File, error) {
	if err := validateGenerateOptions(opts); err != nil {
		return File{}, err
	}

	files, err := discoverGenerateFiles(opts)
	if err != nil {
		return File{}, err
	}

	languageMap := languagemapping.ToMap(
		opts.LanguageMapping,
	)

	items := make(
		[]Item,
		0,
		len(files),
	)

	for _, file := range files {
		item, matched, err := buildGeneratedItem(
			file,
			opts,
			languageMap,
		)
		if err != nil {
			return File{}, err
		}

		if !matched {
			continue
		}

		items = append(
			items,
			item,
		)
	}

	if len(items) == 0 {
		return File{}, fmt.Errorf(
			"no files matched name pattern %q",
			opts.NamePattern,
		)
	}

	return File{
		Items: items,
	}, nil
}

func validateGenerateOptions(
	opts GenerateOptions,
) error {
	if len(nonEmptyStrings(opts.Paths)) == 0 {
		return errors.New("no paths provided")
	}

	if strings.TrimSpace(opts.NamePattern) == "" {
		return errors.New("name pattern is required")
	}

	if strings.TrimSpace(opts.FilenamePattern) == "" {
		return errors.New("filename pattern is required")
	}

	if err := ValidatePathPattern(
		opts.NamePattern,
	); err != nil {
		return fmt.Errorf(
			"invalid name pattern: %w",
			err,
		)
	}

	if err := ValidateRenderPattern(
		opts.FilenamePattern,
	); err != nil {
		return fmt.Errorf(
			"invalid filename pattern: %w",
			err,
		)
	}

	if !strings.Contains(
		opts.NamePattern,
		"{lang}",
	) && strings.TrimSpace(opts.BaseLang) == "" {
		return errors.New(
			"base language is required when name pattern does not contain {lang}",
		)
	}

	if err := ValidatePatternSources(
		opts.NamePattern,
		opts.FilenamePattern,
		opts.BaseLang,
	); err != nil {
		return err
	}

	return nil
}

func discoverGenerateFiles(
	opts GenerateOptions,
) ([]DiscoveredFile, error) {
	paths := nonEmptyStrings(
		opts.Paths,
	)

	files, err := DiscoverFiles(
		paths,
	)
	if err != nil {
		return nil, err
	}

	files, err = FilterFiles(
		files,
		opts.ExcludePatterns,
	)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func nonEmptyStrings(
	values []string,
) []string {
	result := make(
		[]string,
		0,
		len(values),
	)

	for _, value := range values {
		value = strings.TrimSpace(value)

		if value == "" {
			continue
		}

		result = append(
			result,
			value,
		)
	}

	return result
}

func buildGeneratedItem(
	file DiscoveredFile,
	opts GenerateOptions,
	languageMap map[string]string,
) (Item, bool, error) {
	parts, matched, err := MatchPathPattern(
		opts.NamePattern,
		file.RelativePath,
	)
	if err != nil {
		return Item{}, false, err
	}

	if !matched {
		return Item{}, false, nil
	}

	lang, err := resolveLanguage(
		parts,
		opts.BaseLang,
		languageMap,
		file.RelativePath,
	)
	if err != nil {
		return Item{}, false, err
	}

	filename := strings.TrimSpace(
		RenderPattern(
			opts.FilenamePattern,
			parts,
		),
	)

	if filename == "" {
		return Item{}, false, fmt.Errorf(
			"file %q resolves to an empty filename",
			file.RelativePath,
		)
	}

	srcPath, err := relativeSourcePath(
		opts.ManifestDir,
		file.Path,
	)
	if err != nil {
		return Item{}, false, err
	}

	return Item{
		SrcPath: srcPath,
		Params: lokexupload.UploadParams{
			"filename": filename,
			"lang_iso": lang,
		},
	}, true, nil
}

func resolveLanguage(
	parts PathParts,
	baseLang string,
	languageMap map[string]string,
	relativePath string,
) (string, error) {
	lang := strings.TrimSpace(
		parts["lang"],
	)

	if lang == "" {
		lang = strings.TrimSpace(
			baseLang,
		)

		parts["lang"] = lang
	}

	if mapped, ok := languageMap[lang]; ok {
		lang = mapped
	}

	if lang == "" {
		return "", fmt.Errorf(
			"file %q resolves to an empty language",
			relativePath,
		)
	}

	return lang, nil
}

func relativeSourcePath(
	manifestDir string,
	sourcePath string,
) (string, error) {
	if strings.TrimSpace(manifestDir) == "" {
		manifestDir = "."
	}

	manifestDirAbs, err := filepath.Abs(
		manifestDir,
	)
	if err != nil {
		return "", fmt.Errorf(
			"resolve manifest directory %q: %w",
			manifestDir,
			err,
		)
	}

	sourceAbs, err := filepath.Abs(
		sourcePath,
	)
	if err != nil {
		return "", fmt.Errorf(
			"resolve source path %q: %w",
			sourcePath,
			err,
		)
	}

	if !sameVolume(
		manifestDirAbs,
		sourceAbs,
	) {
		return filepath.ToSlash(
			sourceAbs,
		), nil
	}

	relative, err := filepath.Rel(
		manifestDirAbs,
		sourceAbs,
	)
	if err != nil {
		return "", fmt.Errorf(
			"resolve source path %q relative to manifest directory %q: %w",
			sourcePath,
			manifestDir,
			err,
		)
	}

	return filepath.ToSlash(relative), nil
}

func sameVolume(
	left string,
	right string,
) bool {
	return strings.EqualFold(
		filepath.VolumeName(left),
		filepath.VolumeName(right),
	)
}
