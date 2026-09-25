package uploadmanifest

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

type PathParts map[string]string

type compiledPathPattern struct {
	regexp   *regexp.Regexp
	captures []string
}

var placeholderNameRE = regexp.MustCompile(
	`^[A-Za-z_][A-Za-z0-9_-]*$`,
)

func RenderPattern(
	pattern string,
	parts PathParts,
) string {
	result := pattern

	for name, value := range parts {
		result = strings.ReplaceAll(
			result,
			"{"+name+"}",
			value,
		)
	}

	return result
}

func MatchPathPattern(
	pattern string,
	path string,
) (PathParts, bool, error) {
	compiled, err := compilePathPattern(pattern)
	if err != nil {
		return nil, false, err
	}

	path = filepath.ToSlash(path)

	match := compiled.regexp.FindStringSubmatch(path)
	if match == nil {
		return nil, false, nil
	}

	parts := make(
		PathParts,
		len(compiled.captures),
	)

	for i, name := range compiled.captures {
		parts[name] = match[i+1]
	}

	return parts, true, nil
}

func ValidatePathPattern(
	pattern string,
) error {
	pattern = strings.TrimSpace(pattern)

	if pattern == "" {
		return errors.New("name pattern is empty")
	}

	placeholders, err := patternPlaceholders(pattern)
	if err != nil {
		return err
	}

	seen := make(map[string]struct{})

	for _, name := range placeholders {
		if _, ok := seen[name]; ok {
			return fmt.Errorf(
				"duplicate name pattern placeholder {%s}",
				name,
			)
		}

		seen[name] = struct{}{}
	}

	return nil
}

func ValidateRenderPattern(
	pattern string,
) error {
	pattern = strings.TrimSpace(pattern)

	if pattern == "" {
		return errors.New("render pattern is empty")
	}

	if strings.Contains(
		pattern,
		"*",
	) {
		return errors.New(
			"render pattern must not contain wildcards",
		)
	}

	_, err := patternPlaceholders(pattern)

	return err
}

func ValidatePatternSources(
	namePattern string,
	filenamePattern string,
	baseLang string,
) error {
	source, err := patternPlaceholders(namePattern)
	if err != nil {
		return err
	}

	rendered, err := patternPlaceholders(filenamePattern)
	if err != nil {
		return err
	}

	available := make(
		map[string]struct{},
		len(source)+1,
	)

	for _, name := range source {
		available[name] = struct{}{}
	}

	if strings.TrimSpace(baseLang) != "" {
		available["lang"] = struct{}{}
	}

	for _, name := range rendered {
		if _, ok := available[name]; ok {
			continue
		}

		return fmt.Errorf(
			"filename pattern uses {%s}, but name pattern does not provide {%s}",
			name,
			name,
		)
	}

	return nil
}

func compilePathPattern(
	pattern string,
) (*compiledPathPattern, error) {
	pattern = filepath.ToSlash(
		strings.TrimSpace(pattern),
	)

	if err := ValidatePathPattern(pattern); err != nil {
		return nil, err
	}

	source, captures, err := buildPathPatternRegex(
		pattern,
	)
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(source)
	if err != nil {
		return nil, fmt.Errorf(
			"compile name pattern: %w",
			err,
		)
	}

	return &compiledPathPattern{
		regexp:   re,
		captures: captures,
	}, nil
}

func buildPathPatternRegex(
	pattern string,
) (string, []string, error) {
	var result strings.Builder
	var captures []string

	result.WriteString("^")

	for i := 0; i < len(pattern); {
		fragment, capture, next, err := compilePathPatternToken(
			pattern,
			i,
		)
		if err != nil {
			return "", nil, err
		}

		result.WriteString(fragment)

		if capture != "" {
			captures = append(
				captures,
				capture,
			)
		}

		i = next
	}

	result.WriteString("$")

	return result.String(), captures, nil
}

func compilePathPatternToken(
	pattern string,
	index int,
) (string, string, int, error) {
	switch {
	case strings.HasPrefix(
		pattern[index:],
		"**/",
	):
		return `(?:.*/)?`, "", index + 3, nil

	case strings.HasPrefix(
		pattern[index:],
		"**",
	):
		return `.*`, "", index + 2, nil

	case pattern[index] == '*':
		return `[^/]*`, "", index + 1, nil

	case pattern[index] == '{':
		return compilePlaceholderToken(
			pattern,
			index,
		)

	default:
		return compileLiteralToken(
			pattern,
			index,
		)
	}
}

func compilePlaceholderToken(
	pattern string,
	index int,
) (string, string, int, error) {
	end := strings.IndexByte(
		pattern[index:],
		'}',
	)

	if end == -1 {
		return "", "", 0, errors.New(
			"unclosed pattern placeholder",
		)
	}

	name := pattern[index+1 : index+end]

	return capturePattern(name),
		name,
		index + end + 1,
		nil
}

func capturePattern(
	name string,
) string {
	switch name {
	case "lang", "ext":
		return `([^/.]+?)`

	default:
		return `([^/]+?)`
	}
}

func compileLiteralToken(
	pattern string,
	index int,
) (string, string, int, error) {
	end := index

	for end < len(pattern) &&
		pattern[end] != '*' &&
		pattern[end] != '{' {
		end++
	}

	return regexp.QuoteMeta(
		pattern[index:end],
	), "", end, nil
}

func patternPlaceholders(
	pattern string,
) ([]string, error) {
	var result []string

	for i := 0; i < len(pattern); {
		switch pattern[i] {
		case '}':
			return nil, errors.New(
				"unexpected closing brace in pattern",
			)

		case '{':
			end := strings.IndexByte(
				pattern[i+1:],
				'}',
			)

			if end == -1 {
				return nil, errors.New(
					"unclosed pattern placeholder",
				)
			}

			end += i + 1

			if strings.Contains(
				pattern[i+1:end],
				"{",
			) {
				return nil, errors.New(
					"unclosed pattern placeholder",
				)
			}

			name := pattern[i+1 : end]

			if !placeholderNameRE.MatchString(name) {
				return nil, fmt.Errorf(
					"invalid pattern placeholder {%s}",
					name,
				)
			}

			result = append(
				result,
				name,
			)

			i = end + 1

		default:
			i++
		}
	}

	return result, nil
}
