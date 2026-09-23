package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/bodrovis/lokex-cli/internal/appstate"
	"github.com/bodrovis/lokex-cli/internal/global_config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestRootCmd_HasExpectedCommands(t *testing.T) {
	root := RootCmd()

	if root == nil {
		t.Fatal("RootCmd() returned nil")
	}

	if root.Use != "lokex-cli" {
		t.Fatalf("unexpected root Use: got %q, want %q", root.Use, "lokex-cli")
	}

	expected := []string{"version", "gendocs", "download", "upload"}
	for _, name := range expected {
		if findSubcommand(root, name) == nil {
			t.Fatalf("expected subcommand %q to be registered", name)
		}
	}
}

func TestNewPersistentPreRunE_ReturnsGlobalConfigError(t *testing.T) {
	cfg := &global_config.GlobalConfig{
		UserAgent: "test-agent",
	}

	state := &appstate.State{}
	configFile := ""

	expectedErr := errors.New("boom")

	preRun := newPersistentPreRunE(
		cfg,
		state,
		&configFile,
		func(
			*viper.Viper,
			*cobra.Command,
			*global_config.GlobalConfig,
		) error {
			return expectedErr
		},
	)

	cmd := &cobra.Command{
		Use: "upload",
	}

	err := preRun(cmd, nil)

	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, state.Viper)
}

func TestNewPersistentPreRunE_LoadsGlobalConfigAndStoresViper(t *testing.T) {
	dir := t.TempDir()
	configFile := filepath.Join(dir, "lokex.yaml")

	require.NoError(t, os.WriteFile(
		configFile,
		[]byte(`
token: file-token
project-id: file-project
`),
		0o644,
	))

	cfg := &global_config.GlobalConfig{
		UserAgent: "test-agent",
	}

	state := &appstate.State{}

	globalCalled := false

	preRun := newPersistentPreRunE(
		cfg,
		state,
		&configFile,
		func(
			v *viper.Viper,
			cmd *cobra.Command,
			gotCfg *global_config.GlobalConfig,
		) error {
			globalCalled = true

			require.Same(t, cfg, gotCfg)
			require.Equal(t, "file-token", v.GetString("token"))
			require.Equal(t, "file-project", v.GetString("project-id"))
			require.Equal(t, "upload", cmd.Name())

			return nil
		},
	)

	cmd := &cobra.Command{
		Use: "upload",
	}

	err := preRun(cmd, nil)
	require.NoError(t, err)

	require.True(t, globalCalled)
	require.NotNil(t, state.Viper)
}

func TestNewPersistentPreRunE_SkipsConfigWhenAnnotated(t *testing.T) {
	cfg := &global_config.GlobalConfig{
		UserAgent: "test-agent",
	}

	state := &appstate.State{}
	configFile := "does-not-exist.yaml"

	globalCalled := false

	preRun := newPersistentPreRunE(
		cfg,
		state,
		&configFile,
		func(
			*viper.Viper,
			*cobra.Command,
			*global_config.GlobalConfig,
		) error {
			globalCalled = true
			return nil
		},
	)

	cmd := &cobra.Command{
		Use: "version",
	}

	markSkipConfig(cmd)

	err := preRun(cmd, nil)

	require.NoError(t, err)
	require.False(t, globalCalled)
	require.Nil(t, state.Viper)
}

func TestNewPersistentPreRunE_ReturnsGlobalLoadError(t *testing.T) {
	cfg := &global_config.GlobalConfig{
		UserAgent: "test-agent",
	}

	state := &appstate.State{}
	configFile := ""

	wantErr := errors.New("global load failed")

	preRun := newPersistentPreRunE(
		cfg,
		state,
		&configFile,
		func(
			*viper.Viper,
			*cobra.Command,
			*global_config.GlobalConfig,
		) error {
			return wantErr
		},
	)

	cmd := &cobra.Command{
		Use: "upload",
	}

	err := preRun(cmd, nil)

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, state.Viper)
}

func TestNewPersistentPreRunE_ReturnsConfigReadError(t *testing.T) {
	cfg := &global_config.GlobalConfig{
		UserAgent: "test-agent",
	}

	state := &appstate.State{}
	configFile := filepath.Join(t.TempDir(), "missing.yaml")

	globalCalled := false

	preRun := newPersistentPreRunE(
		cfg,
		state,
		&configFile,
		func(
			*viper.Viper,
			*cobra.Command,
			*global_config.GlobalConfig,
		) error {
			globalCalled = true
			return nil
		},
	)

	cmd := &cobra.Command{
		Use: "upload",
	}

	err := preRun(cmd, nil)

	require.Error(t, err)
	require.False(t, globalCalled)
	require.Nil(t, state.Viper)
}

func findSubcommand(root *cobra.Command, name string) *cobra.Command {
	for _, cmd := range root.Commands() {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}
