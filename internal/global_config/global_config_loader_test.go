package global_config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bodrovis/lokex-cli/internal/viper_helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestLoadGlobalConfig_ConfigOnly(t *testing.T) {
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "lokex.yaml")

	require.NoError(t, os.WriteFile(cfgFile, []byte(`
token: file-token
project-id: file-project
http-timeout: 45s
retries: 2
context-timeout: 100s
base-url: https://example.com
`), 0o644))

	const userAgent = "lokex-cli/test"

	v := newTestConfigViper(t, cfgFile)
	cmd := newTestCommand(t, userAgent)

	cfg := &GlobalConfig{
		UserAgent: userAgent,
	}

	require.NoError(
		t,
		LoadGlobalConfig(v, cmd, cfg),
	)

	require.Equal(t, "file-token", cfg.Token)
	require.Equal(t, "file-project", cfg.ProjectID)
	require.Equal(t, 45*time.Second, cfg.HTTPTimeout)
	require.Equal(t, 100*time.Second, cfg.ContextTimeout)
	require.Equal(t, 2, cfg.MaxRetries)
	require.Equal(t, "https://example.com", cfg.BaseURL)
	require.Equal(t, userAgent, cfg.UserAgent)
}

func TestLoadGlobalConfig_NoConfigFile(t *testing.T) {
	const userAgent = "lokex-cli/test"

	v := newTestConfigViper(t, "")
	cmd := newTestCommand(t, userAgent)

	cfg := &GlobalConfig{
		UserAgent: userAgent,
	}

	require.NoError(
		t,
		LoadGlobalConfig(v, cmd, cfg),
	)

	require.Equal(t, userAgent, cfg.UserAgent)
	require.Equal(t, -1, cfg.MaxRetries)
	require.Equal(t, 150*time.Second, cfg.ContextTimeout)
}

func TestLoadGlobalConfig_ConfigCanSetZeroValues(t *testing.T) {
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "lokex.yaml")

	require.NoError(t, os.WriteFile(cfgFile, []byte(`
token: file-token
project-id: file-project
http-timeout: 0s
retries: -1
backoff-initial: 0s
backoff-max: 0s
poll-initial-wait: 0s
poll-max-wait: 0s
context-timeout: 0s
`), 0o644))

	const userAgent = "lokex-cli/test"

	v := newTestConfigViper(t, cfgFile)
	cmd := newTestCommand(t, userAgent)

	cfg := &GlobalConfig{
		UserAgent: userAgent,
	}

	require.NoError(
		t,
		LoadGlobalConfig(v, cmd, cfg),
	)

	require.Equal(t, time.Duration(0), cfg.HTTPTimeout)
	require.Equal(t, -1, cfg.MaxRetries)
	require.Equal(t, time.Duration(0), cfg.InitialBackoff)
	require.Equal(t, time.Duration(0), cfg.MaxBackoff)
	require.Equal(t, time.Duration(0), cfg.PollInitialWait)
	require.Equal(t, time.Duration(0), cfg.PollMaxWait)
	require.Equal(t, time.Duration(0), cfg.ContextTimeout)
}

func TestLoadGlobalConfig_FlagsOverrideConfig(t *testing.T) {
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "lokex.yaml")

	require.NoError(t, os.WriteFile(cfgFile, []byte(`
token: file-token
project-id: file-project
http-timeout: 30s
retries: 2
context-timeout: 60s
base-url: https://example.com
`), 0o644))

	const userAgent = "lokex-cli/test"

	v := newTestConfigViper(t, cfgFile)

	cmd := newTestCommand(
		t,
		userAgent,
		"--token=cli-token",
	)

	cfg := &GlobalConfig{
		UserAgent: userAgent,
	}

	require.NoError(
		t,
		LoadGlobalConfig(v, cmd, cfg),
	)

	require.Equal(t, "cli-token", cfg.Token)
	require.Equal(t, "file-project", cfg.ProjectID)
	require.Equal(t, 30*time.Second, cfg.HTTPTimeout)
	require.Equal(t, 2, cfg.MaxRetries)
	require.Equal(t, 60*time.Second, cfg.ContextTimeout)
	require.Equal(t, "https://example.com", cfg.BaseURL)
}

func TestLoadGlobalConfig_FlagOverridesExplicitZeroFromConfig(t *testing.T) {
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "lokex.yaml")

	require.NoError(t, os.WriteFile(cfgFile, []byte(`
http-timeout: 0s
`), 0o644))

	const userAgent = "lokex-cli/test"

	v := newTestConfigViper(t, cfgFile)

	cmd := newTestCommand(
		t,
		userAgent,
		"--http-timeout=20s",
	)

	cfg := &GlobalConfig{
		UserAgent: userAgent,
	}

	require.NoError(
		t,
		LoadGlobalConfig(v, cmd, cfg),
	)

	require.Equal(t, 20*time.Second, cfg.HTTPTimeout)
}

func TestLoadGlobalConfig_EnvOverridesConfig(t *testing.T) {
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "lokex.yaml")

	require.NoError(t, os.WriteFile(cfgFile, []byte(`
token: file-token
project-id: file-project
http-timeout: 10s
`), 0o644))

	t.Setenv("LOKEX_TOKEN", "env-token")
	t.Setenv("LOKEX_HTTP_TIMEOUT", "20s")

	const userAgent = "lokex-cli/test"

	v := newTestConfigViper(t, cfgFile)
	cmd := newTestCommand(t, userAgent)

	cfg := &GlobalConfig{
		UserAgent: userAgent,
	}

	require.NoError(
		t,
		LoadGlobalConfig(v, cmd, cfg),
	)

	require.Equal(t, "env-token", cfg.Token)
	require.Equal(t, "file-project", cfg.ProjectID)
	require.Equal(t, 20*time.Second, cfg.HTTPTimeout)
}

func TestLoadGlobalConfig_Precedence(t *testing.T) {
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "lokex.yaml")

	require.NoError(t, os.WriteFile(cfgFile, []byte(`
token: file-token
project-id: file-project
http-timeout: 10s
context-timeout: 90s
`), 0o644))

	t.Setenv("LOKEX_TOKEN", "env-token")
	t.Setenv("LOKEX_HTTP_TIMEOUT", "20s")

	const userAgent = "lokex-cli/test"

	v := newTestConfigViper(t, cfgFile)

	cmd := newTestCommand(
		t,
		userAgent,
		"--token=cli-token",
	)

	cfg := &GlobalConfig{
		UserAgent: userAgent,
	}

	require.NoError(
		t,
		LoadGlobalConfig(v, cmd, cfg),
	)

	// CLI > env > config > default.
	require.Equal(t, "cli-token", cfg.Token)
	require.Equal(t, 20*time.Second, cfg.HTTPTimeout)
	require.Equal(t, "file-project", cfg.ProjectID)
	require.Equal(t, 90*time.Second, cfg.ContextTimeout)

	// Not provided anywhere: declared default.
	require.Equal(t, -1, cfg.MaxRetries)
}

func TestLoadGlobalConfig_Errors(t *testing.T) {
	t.Run("nil viper", func(t *testing.T) {
		cmd := &cobra.Command{Use: "test"}
		cfg := &GlobalConfig{}

		err := LoadGlobalConfig(nil, cmd, cfg)

		require.EqualError(t, err, "viper is nil")
	})

	t.Run("nil command", func(t *testing.T) {
		v := viper.New()
		cfg := &GlobalConfig{}

		err := LoadGlobalConfig(v, nil, cfg)

		require.EqualError(t, err, "command is nil")
	})

	t.Run("nil config", func(t *testing.T) {
		v := viper.New()
		cmd := &cobra.Command{Use: "test"}

		err := LoadGlobalConfig(v, cmd, nil)

		require.EqualError(t, err, "global config is nil")
	})
}

func TestLoadGlobalConfig_ReturnsApplyChangedFlagsError(t *testing.T) {
	v := viper.New()

	cmd := &cobra.Command{
		Use: "test",
	}

	cmd.Flags().Int(
		"token",
		0,
		"wrong type on purpose",
	)

	require.NoError(
		t,
		cmd.Flags().Set("token", "123"),
	)

	cfg := &GlobalConfig{
		UserAgent: "lokex-cli/test",
	}

	err := LoadGlobalConfig(
		v,
		cmd,
		cfg,
	)

	require.Error(t, err)
	require.Contains(
		t,
		err.Error(),
		"apply global flags:",
	)
}

func TestLoadGlobalConfig_ReturnsDecodeError(t *testing.T) {
	v := viper.New()

	v.Set(
		"retries",
		"definitely-not-an-int",
	)

	cmd := newTestCommand(
		t,
		"lokex-cli/test",
	)

	cfg := &GlobalConfig{
		UserAgent: "lokex-cli/test",
	}

	err := LoadGlobalConfig(
		v,
		cmd,
		cfg,
	)

	require.Error(t, err)
	require.Contains(
		t,
		err.Error(),
		"decode global config:",
	)
}

func TestLoadGlobalConfig_DoesNotValidateResolvedConfig(t *testing.T) {
	v := viper.New()

	v.Set("token", "token")
	v.Set("project-id", "project-id")
	v.Set("http-timeout", "-1s")

	cmd := newTestCommand(
		t,
		"lokex-cli/test",
	)

	cfg := &GlobalConfig{
		UserAgent: "lokex-cli/test",
	}

	require.NoError(
		t,
		LoadGlobalConfig(v, cmd, cfg),
	)

	require.Equal(
		t,
		-time.Second,
		cfg.HTTPTimeout,
	)

	require.EqualError(
		t,
		cfg.Validate(),
		"http-timeout must be >= 0",
	)
}

func newTestConfigViper(
	t *testing.T,
	configFile string,
) *viper.Viper {
	t.Helper()

	v := viper_helpers.NewConfigViper(
		configFile,
		"LOKEX",
	)

	require.NoError(
		t,
		viper_helpers.ReadOptionalConfig(v, configFile),
	)

	return v
}

func newTestCommand(
	t *testing.T,
	userAgent string,
	args ...string,
) *cobra.Command {
	t.Helper()

	cmd := &cobra.Command{
		Use: "test",
	}

	BindPersistentFlags(
		cmd.PersistentFlags(),
		userAgent,
	)

	require.NoError(
		t,
		cmd.ParseFlags(args),
	)

	return cmd
}
