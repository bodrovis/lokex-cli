package global_config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestLoadGlobalConfigInput_ConfigOnly(t *testing.T) {
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "lokex.yaml")

	err := os.WriteFile(cfgFile, []byte(`
token: file-token
project-id: file-project
http-timeout: 45s
retries: 2
context-timeout: 100s
base-url: https://example.com
`), 0o644)
	require.NoError(t, err)

	input, err := LoadGlobalConfigInput("lokex-cli/test", LoadOptions{
		ConfigFile: cfgFile,
		EnvPrefix:  "LOKEX",
	})
	require.NoError(t, err)
	require.NotNil(t, input)

	require.NotNil(t, input.Token)
	require.Equal(t, "file-token", *input.Token)

	require.NotNil(t, input.ProjectID)
	require.Equal(t, "file-project", *input.ProjectID)

	require.NotNil(t, input.HTTPTimeout)
	require.Equal(t, 45*time.Second, *input.HTTPTimeout)

	require.NotNil(t, input.ContextTimeout)
	require.Equal(t, 100*time.Second, *input.ContextTimeout)

	require.NotNil(t, input.MaxRetries)
	require.Equal(t, 2, *input.MaxRetries)

	require.NotNil(t, input.BaseURL)
	require.Equal(t, "https://example.com", *input.BaseURL)
}

func TestLoadGlobalConfigInput_NoConfigFile(t *testing.T) {
	input, err := LoadGlobalConfigInput("lokex-cli/test", LoadOptions{
		ConfigFile: "",
		EnvPrefix:  "LOKEX",
	})
	require.NoError(t, err)
	require.NotNil(t, input)
}

func TestLoadGlobalConfigInput_ConfigCanSetZeroValues(t *testing.T) {
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "lokex.yaml")

	err := os.WriteFile(cfgFile, []byte(`
token: file-token
project-id: file-project
http-timeout: 0s
retries: -1
backoff-initial: 0s
backoff-max: 0s
poll-initial-wait: 0s
poll-max-wait: 0s
context-timeout: 0s
`), 0o644)
	require.NoError(t, err)

	input, err := LoadGlobalConfigInput("lokex-cli/test", LoadOptions{
		ConfigFile: cfgFile,
		EnvPrefix:  "LOKEX",
	})
	require.NoError(t, err)

	require.Equal(t, time.Duration(0), *input.HTTPTimeout)
	require.Equal(t, time.Duration(0), *input.ContextTimeout)
	require.Equal(t, -1, *input.MaxRetries)
	require.Equal(t, time.Duration(0), *input.InitialBackoff)
	require.Equal(t, time.Duration(0), *input.MaxBackoff)
	require.Equal(t, time.Duration(0), *input.PollInitialWait)
	require.Equal(t, time.Duration(0), *input.PollMaxWait)
}

func TestApplyGlobalInput_FlagsOverrideInput(t *testing.T) {
	cfg := &GlobalConfig{
		UserAgent: "lokex-cli/test",
	}

	cmd := &cobra.Command{Use: "test"}
	BindPersistentFlags(cmd.PersistentFlags(), cfg)

	err := cmd.ParseFlags([]string{
		"--token=cli-token",
	})
	require.NoError(t, err)

	projectID := "file-project"
	timeout := 30 * time.Second
	retries := 2
	token := "file-token"
	contextTimeout := 60 * time.Second
	baseUrl := "https://example.com"

	input := &GlobalConfigInput{
		Token:          &token,
		ProjectID:      &projectID,
		HTTPTimeout:    &timeout,
		MaxRetries:     &retries,
		ContextTimeout: &contextTimeout,
		BaseURL:        &baseUrl,
	}

	ApplyGlobalInput(cmd, cfg, input)

	require.Equal(t, "cli-token", cfg.Token)
	require.Equal(t, "file-project", cfg.ProjectID)
	require.Equal(t, 30*time.Second, cfg.HTTPTimeout)
	require.Equal(t, 60*time.Second, cfg.ContextTimeout)
	require.Equal(t, 2, cfg.MaxRetries)
	require.Equal(t, "https://example.com", cfg.BaseURL)
}

func TestApplyGlobalInput_InputOnly(t *testing.T) {
	cfg := &GlobalConfig{
		UserAgent: "lokex-cli/test",
	}

	cmd := &cobra.Command{Use: "test"}
	BindPersistentFlags(cmd.PersistentFlags(), cfg)

	token := "file-token"
	projectID := "file-project"
	timeout := 45 * time.Second
	contextTimeout := 60 * time.Second

	input := &GlobalConfigInput{
		Token:          &token,
		ProjectID:      &projectID,
		HTTPTimeout:    &timeout,
		ContextTimeout: &contextTimeout,
	}

	ApplyGlobalInput(cmd, cfg, input)

	require.Equal(t, "file-token", cfg.Token)
	require.Equal(t, "file-project", cfg.ProjectID)
	require.Equal(t, 45*time.Second, cfg.HTTPTimeout)
	require.Equal(t, 60*time.Second, cfg.ContextTimeout)
}

func TestApplyGlobalInput_InputCanSetZeroValues(t *testing.T) {
	cfg := &GlobalConfig{
		UserAgent:       "lokex-cli/test",
		HTTPTimeout:     10 * time.Second,
		MaxRetries:      5,
		InitialBackoff:  1 * time.Second,
		MaxBackoff:      2 * time.Second,
		PollInitialWait: 3 * time.Second,
		PollMaxWait:     4 * time.Second,
		ContextTimeout:  5 * time.Second,
	}

	cmd := &cobra.Command{Use: "test"}
	BindPersistentFlags(cmd.PersistentFlags(), cfg)

	zeroDuration := time.Duration(0)
	minusOne := -1

	input := &GlobalConfigInput{
		HTTPTimeout:     &zeroDuration,
		MaxRetries:      &minusOne,
		InitialBackoff:  &zeroDuration,
		MaxBackoff:      &zeroDuration,
		PollInitialWait: &zeroDuration,
		PollMaxWait:     &zeroDuration,
		ContextTimeout:  &zeroDuration,
	}

	ApplyGlobalInput(cmd, cfg, input)

	require.Equal(t, time.Duration(0), cfg.HTTPTimeout)
	require.Equal(t, -1, cfg.MaxRetries)
	require.Equal(t, time.Duration(0), cfg.InitialBackoff)
	require.Equal(t, time.Duration(0), cfg.MaxBackoff)
	require.Equal(t, time.Duration(0), cfg.PollInitialWait)
	require.Equal(t, time.Duration(0), cfg.PollMaxWait)
	require.Equal(t, time.Duration(0), cfg.ContextTimeout)
}

func TestApplyGlobalInput_FlagOverridesExplicitZeroFromInput(t *testing.T) {
	cfg := &GlobalConfig{
		UserAgent: "lokex-cli/test",
	}

	cmd := &cobra.Command{Use: "test"}
	BindPersistentFlags(cmd.PersistentFlags(), cfg)

	err := cmd.ParseFlags([]string{
		"--http-timeout=20s",
	})
	require.NoError(t, err)

	zeroDuration := time.Duration(0)
	input := &GlobalConfigInput{
		HTTPTimeout: &zeroDuration,
	}

	ApplyGlobalInput(cmd, cfg, input)

	require.Equal(t, 20*time.Second, cfg.HTTPTimeout)
}
