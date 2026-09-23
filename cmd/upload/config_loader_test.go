package upload

import (
	"testing"

	"github.com/bodrovis/lokex-cli/internal/params"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestLoadUploadConfig_ConfigValues(t *testing.T) {
	v := viper.New()

	v.Set("upload.filename", "en.json")
	v.Set("upload.lang-iso", "en")
	v.Set("upload.poll", true)
	v.Set("upload.apply-tm", true)

	cmd := newTestUploadCommand()

	cfg := &UploadConfig{}

	err := LoadUploadConfig(
		v,
		cmd,
		cfg,
	)
	require.NoError(t, err)

	require.NotNil(t, cfg.Filename)
	require.Equal(t, "en.json", *cfg.Filename)

	require.NotNil(t, cfg.LangISO)
	require.Equal(t, "en", *cfg.LangISO)

	require.NotNil(t, cfg.Poll)
	require.True(t, *cfg.Poll)

	require.NotNil(t, cfg.ApplyTM)
	require.True(t, *cfg.ApplyTM)
}

func TestLoadUploadConfig_FlagsOverrideConfig(t *testing.T) {
	v := viper.New()

	v.Set("upload.filename", "default.json")
	v.Set("upload.lang-iso", "fr")
	v.Set("upload.poll", true)
	v.Set("upload.apply-tm", true)

	cmd := newTestUploadCommand()

	require.NoError(
		t,
		cmd.ParseFlags([]string{
			"--filename=explicit.json",
			"--lang-iso=en",
			"--poll=false",
			"--apply-tm=false",
		}),
	)

	cfg := &UploadConfig{}

	err := LoadUploadConfig(
		v,
		cmd,
		cfg,
	)
	require.NoError(t, err)

	require.Equal(t, "explicit.json", *cfg.Filename)
	require.Equal(t, "en", *cfg.LangISO)

	require.NotNil(t, cfg.Poll)
	require.False(t, *cfg.Poll)

	require.NotNil(t, cfg.ApplyTM)
	require.False(t, *cfg.ApplyTM)
}

func TestLoadUploadConfig_NilArguments(t *testing.T) {
	tests := []struct {
		name    string
		v       *viper.Viper
		cmd     *cobra.Command
		cfg     *UploadConfig
		wantErr string
	}{
		{
			name:    "nil viper",
			v:       nil,
			cmd:     newTestUploadCommand(),
			cfg:     &UploadConfig{},
			wantErr: "viper is nil",
		},
		{
			name:    "nil command",
			v:       viper.New(),
			cmd:     nil,
			cfg:     &UploadConfig{},
			wantErr: "upload command is nil",
		},
		{
			name:    "nil config",
			v:       viper.New(),
			cmd:     newTestUploadCommand(),
			cfg:     nil,
			wantErr: "upload config is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := LoadUploadConfig(
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

func TestLoadUploadConfig_ReturnsDecodeError(t *testing.T) {
	v := viper.New()

	v.Set(
		"upload.filter-task-id",
		"not-an-int64",
	)

	cmd := newTestUploadCommand()

	cfg := &UploadConfig{}

	err := LoadUploadConfig(
		v,
		cmd,
		cfg,
	)

	require.Error(t, err)
	require.Contains(
		t,
		err.Error(),
		"decode upload config:",
	)
}

func TestLoadUploadConfig_UnsetValuesRemainNil(t *testing.T) {
	v := viper.New()
	v.Set("upload.filename", "en.json")

	cmd := newTestUploadCommand()

	cfg := &UploadConfig{}

	require.NoError(
		t,
		LoadUploadConfig(v, cmd, cfg),
	)

	require.NotNil(t, cfg.Filename)
	require.Nil(t, cfg.Poll)
	require.Nil(t, cfg.ApplyTM)
	require.Nil(t, cfg.Tags)
	require.Nil(t, cfg.FilterTaskID)
}

func newTestUploadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "upload",
	}

	params.BindFlags(
		cmd,
		uploadParamSpecs,
	)

	return cmd
}
