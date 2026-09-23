package download

import (
	"testing"

	"github.com/bodrovis/lokex-cli/internal/params"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestLoadDownloadConfig_DefaultsAndConfig(t *testing.T) {
	v := viper.New()

	v.Set("download.format", "json")
	v.Set("download.async", true)

	cmd := newTestDownloadCommand()

	cfg := &DownloadConfig{}

	err := LoadDownloadConfig(
		v,
		cmd,
		cfg,
	)
	require.NoError(t, err)

	require.NotNil(t, cfg.Out)
	require.Equal(t, "./locales", *cfg.Out)

	require.NotNil(t, cfg.Format)
	require.Equal(t, "json", *cfg.Format)

	require.NotNil(t, cfg.Async)
	require.True(t, *cfg.Async)
}

func TestLoadDownloadConfig_FlagsOverrideConfig(t *testing.T) {
	v := viper.New()

	v.Set("download.format", "json")
	v.Set("download.out", "./from-config")
	v.Set("download.async", true)

	cmd := newTestDownloadCommand()

	require.NoError(
		t,
		cmd.ParseFlags([]string{
			"--format=xml",
			"--out=./from-cli",
			"--async=false",
		}),
	)

	cfg := &DownloadConfig{}

	err := LoadDownloadConfig(
		v,
		cmd,
		cfg,
	)
	require.NoError(t, err)

	require.Equal(t, "xml", *cfg.Format)
	require.Equal(t, "./from-cli", *cfg.Out)
	require.False(t, *cfg.Async)
}

func TestLoadDownloadConfig_NilArguments(t *testing.T) {
	tests := []struct {
		name    string
		v       *viper.Viper
		cmd     *cobra.Command
		cfg     *DownloadConfig
		wantErr string
	}{
		{
			name:    "nil viper",
			v:       nil,
			cmd:     newTestDownloadCommand(),
			cfg:     &DownloadConfig{},
			wantErr: "viper is nil",
		},
		{
			name:    "nil command",
			v:       viper.New(),
			cmd:     nil,
			cfg:     &DownloadConfig{},
			wantErr: "download command is nil",
		},
		{
			name:    "nil config",
			v:       viper.New(),
			cmd:     newTestDownloadCommand(),
			cfg:     nil,
			wantErr: "download config is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := LoadDownloadConfig(
				tt.v,
				tt.cmd,
				tt.cfg,
			)

			require.EqualError(
				t,
				err,
				tt.wantErr,
			)
		})
	}
}

func TestLoadDownloadConfig_ReturnsDecodeError(t *testing.T) {
	v := viper.New()

	v.Set(
		"download.filter-task-id",
		"not-an-int64",
	)

	cmd := newTestDownloadCommand()

	cfg := &DownloadConfig{}

	err := LoadDownloadConfig(
		v,
		cmd,
		cfg,
	)

	require.Error(t, err)
	require.Contains(
		t,
		err.Error(),
		"decode download config:",
	)
}

func newTestDownloadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "download",
	}

	params.BindFlags(
		cmd,
		downloadParamSpecs,
	)

	return cmd
}
