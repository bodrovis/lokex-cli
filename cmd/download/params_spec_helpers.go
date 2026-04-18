package download

import (
	"encoding/json"
	"fmt"
	"strings"

	setters "github.com/bodrovis/lokex-cli/internal/params"
	lokexdownload "github.com/bodrovis/lokex/v2/client/download"
	"github.com/spf13/cobra"
)

func reqString(
	apiKey string,
	get func(*Flags) string,
) func(*cobra.Command, *Flags, *DownloadConfig, lokexdownload.DownloadParams) error {
	return func(_ *cobra.Command, flags *Flags, _ *DownloadConfig, req lokexdownload.DownloadParams) error {
		setters.SetString(req, apiKey, get(flags))

		return nil
	}
}

func reqDirectString(
	apiKey string,
	get func(*Flags) string,
) func(*cobra.Command, *Flags, *DownloadConfig, lokexdownload.DownloadParams) error {
	return func(_ *cobra.Command, flags *Flags, _ *DownloadConfig, req lokexdownload.DownloadParams) error {
		req[apiKey] = get(flags)

		return nil
	}
}

func reqStringSlice(
	apiKey string,
	get func(*Flags) []string,
) func(*cobra.Command, *Flags, *DownloadConfig, lokexdownload.DownloadParams) error {
	return func(_ *cobra.Command, flags *Flags, _ *DownloadConfig, req lokexdownload.DownloadParams) error {
		setters.SetStringSlice(req, apiKey, get(flags))

		return nil
	}
}

func reqBoolWithDefault(
	flagName string,
	apiKey string,
	getFlag func(*Flags) bool,
	getDefault func(*DownloadConfig) *bool,
) func(*cobra.Command, *Flags, *DownloadConfig, lokexdownload.DownloadParams) error {
	return func(cmd *cobra.Command, flags *Flags, defaults *DownloadConfig, req lokexdownload.DownloadParams) error {
		var def *bool
		if defaults != nil {
			def = getDefault(defaults)
		}

		setters.SetBoolWithDefault(
			cmd,
			req,
			flagName,
			apiKey,
			getFlag(flags),
			def,
		)

		return nil
	}
}

func reqInt64WithDefault(
	flagName string,
	apiKey string,
	getFlag func(*Flags) int64,
	getDefault func(*DownloadConfig) *int64,
) func(*cobra.Command, *Flags, *DownloadConfig, lokexdownload.DownloadParams) error {
	return func(cmd *cobra.Command, flags *Flags, defaults *DownloadConfig, req lokexdownload.DownloadParams) error {
		var def *int64
		if defaults != nil {
			def = getDefault(defaults)
		}

		setters.SetInt64WithDefault(
			cmd,
			req,
			flagName,
			apiKey,
			getFlag(flags),
			def,
		)

		return nil
	}
}

func reqLanguageMapping(
	flagName string,
	get func(*Flags) string,
) func(*cobra.Command, *Flags, *DownloadConfig, lokexdownload.DownloadParams) error {
	return func(_ *cobra.Command, flags *Flags, _ *DownloadConfig, req lokexdownload.DownloadParams) error {
		raw := strings.TrimSpace(get(flags))
		if raw == "" {
			return nil
		}

		languageMapping, err := parseLanguageMapping(raw)
		if err != nil {
			return fmt.Errorf("parse --%s: %w", flagName, err)
		}

		req["language_mapping"] = languageMapping
		return nil
	}
}

func parseLanguageMapping(raw string) ([]map[string]any, error) {
	var out []map[string]any

	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil, err
	}

	return out, nil
}
