package manifest

import (
	"testing"

	"github.com/bodrovis/lokex-cli/internal/appstate"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestNewCommand(t *testing.T) {
	state := &appstate.State{
		Viper: viper.New(),
	}

	cmd := NewCommand(state)

	require.Equal(t, "manifest", cmd.Use)
	require.NotEmpty(t, cmd.Short)
	require.NotEmpty(t, cmd.Long)

	generateCmd, _, err := cmd.Find([]string{"generate"})
	require.NoError(t, err)
	require.NotNil(t, generateCmd)

	require.Equal(t, "generate", generateCmd.Use)
	require.NotEmpty(t, generateCmd.Short)
	require.NotEmpty(t, generateCmd.Long)
	require.NotEmpty(t, generateCmd.Example)

	require.NotNil(t, generateCmd.Flags().Lookup("path"))
	require.NotNil(t, generateCmd.Flags().Lookup("exclude-pattern"))
	require.NotNil(t, generateCmd.Flags().Lookup("name-pattern"))
	require.NotNil(t, generateCmd.Flags().Lookup("filename-pattern"))
	require.NotNil(t, generateCmd.Flags().Lookup("base-lang"))
	require.NotNil(t, generateCmd.Flags().Lookup("language-mapping"))
	require.NotNil(t, generateCmd.Flags().Lookup("out"))
}

func TestGenerateCommand_FlagDefaults(t *testing.T) {
	state := &appstate.State{
		Viper: viper.New(),
	}

	cmd := newGenerateCommand(state)

	filenamePattern, err := cmd.Flags().GetString(
		"filename-pattern",
	)
	require.NoError(t, err)
	require.Equal(
		t,
		"{name}.{ext}",
		filenamePattern,
	)

	out, err := cmd.Flags().GetString(
		"out",
	)
	require.NoError(t, err)
	require.Equal(
		t,
		"./lokex-manifest.json",
		out,
	)
}

func TestGenerateCommand_RejectsArguments(t *testing.T) {
	state := &appstate.State{
		Viper: viper.New(),
	}

	cmd := newGenerateCommand(state)

	err := cmd.Args(
		cmd,
		[]string{"unexpected"},
	)

	require.Error(t, err)
}

func TestManifestCommand_RejectsArguments(t *testing.T) {
	cmd := NewCommand(
		&appstate.State{
			Viper: viper.New(),
		},
	)

	err := cmd.Args(
		cmd,
		[]string{"unexpected"},
	)

	require.Error(t, err)
}
func TestGenerateCommand_PreRunE_NilState(t *testing.T) {
	cmd := newGenerateCommand(nil)

	err := cmd.PreRunE(
		cmd,
		nil,
	)

	require.EqualError(
		t,
		err,
		"app state is nil",
	)
}

func TestGenerateCommand_PreRunE_NilViper(t *testing.T) {
	cmd := newGenerateCommand(
		&appstate.State{},
	)

	err := cmd.PreRunE(
		cmd,
		nil,
	)

	require.EqualError(
		t,
		err,
		"viper is nil",
	)
}
