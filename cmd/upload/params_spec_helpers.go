package upload

import (
	"strings"

	setters "github.com/bodrovis/lokex-cli/internal/params"
	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
	"github.com/spf13/cobra"
)

func reqString(
	apiKey string,
	get func(*Flags) string,
) func(*cobra.Command, *Flags, *UploadConfig, lokexupload.UploadParams) error {
	return func(_ *cobra.Command, flags *Flags, _ *UploadConfig, req lokexupload.UploadParams) error {
		setters.SetString(req, apiKey, get(flags))

		return nil
	}
}

func reqTrimmedString(
	apiKey string,
	get func(*Flags) string,
) func(*cobra.Command, *Flags, *UploadConfig, lokexupload.UploadParams) error {
	return func(_ *cobra.Command, flags *Flags, _ *UploadConfig, req lokexupload.UploadParams) error {
		req[apiKey] = strings.TrimSpace(get(flags))

		return nil
	}
}

func reqNormalizedFilename(
	apiKey string,
	get func(*Flags) string,
) func(*cobra.Command, *Flags, *UploadConfig, lokexupload.UploadParams) error {
	return func(_ *cobra.Command, flags *Flags, _ *UploadConfig, req lokexupload.UploadParams) error {
		val := strings.TrimSpace(get(flags))
		val = strings.ReplaceAll(val, `\`, `/`)
		req[apiKey] = val
		return nil
	}
}

func reqStringSlice(
	apiKey string,
	get func(*Flags) []string,
) func(*cobra.Command, *Flags, *UploadConfig, lokexupload.UploadParams) error {
	return func(_ *cobra.Command, flags *Flags, _ *UploadConfig, req lokexupload.UploadParams) error {
		setters.SetStringSlice(req, apiKey, get(flags))

		return nil
	}
}

func reqBoolWithDefault(
	flagName string,
	apiKey string,
	getFlag func(*Flags) bool,
	getDefault func(*UploadConfig) *bool,
) func(*cobra.Command, *Flags, *UploadConfig, lokexupload.UploadParams) error {
	return func(cmd *cobra.Command, flags *Flags, defaults *UploadConfig, req lokexupload.UploadParams) error {
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
	getDefault func(*UploadConfig) *int64,
) func(*cobra.Command, *Flags, *UploadConfig, lokexupload.UploadParams) error {
	return func(cmd *cobra.Command, flags *Flags, defaults *UploadConfig, req lokexupload.UploadParams) error {
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
