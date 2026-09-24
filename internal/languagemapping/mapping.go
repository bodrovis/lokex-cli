package languagemapping

import (
	"encoding/json/v2"
	"errors"
	"strings"
)

type Mapping struct {
	OriginalLanguageISO string `json:"original_language_iso"`
	CustomLanguageISO   string `json:"custom_language_iso"`
}

func Parse(
	raw string,
) ([]Mapping, error) {
	raw = strings.TrimSpace(raw)

	if raw == "" {
		return nil, nil
	}

	var mappings []Mapping

	if err := json.Unmarshal(
		[]byte(raw),
		&mappings,
		json.RejectUnknownMembers(true),
	); err != nil {
		return nil, err
	}

	for i, mapping := range mappings {
		if strings.TrimSpace(
			mapping.OriginalLanguageISO,
		) == "" {
			return nil, errors.New(
				"original_language_iso is required",
			)
		}

		if strings.TrimSpace(
			mapping.CustomLanguageISO,
		) == "" {
			return nil, errors.New(
				"custom_language_iso is required",
			)
		}

		mappings[i].OriginalLanguageISO =
			strings.TrimSpace(mapping.OriginalLanguageISO)

		mappings[i].CustomLanguageISO =
			strings.TrimSpace(mapping.CustomLanguageISO)
	}

	return mappings, nil
}

func ToMap(
	mappings []Mapping,
) map[string]string {
	if len(mappings) == 0 {
		return nil
	}

	result := make(
		map[string]string,
		len(mappings),
	)

	for _, mapping := range mappings {
		result[mapping.OriginalLanguageISO] =
			mapping.CustomLanguageISO
	}

	return result
}
