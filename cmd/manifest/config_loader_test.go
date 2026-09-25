package manifest

import (
	"strings"
	"testing"

	"github.com/bodrovis/lokex-cli/internal/params"
	"github.com/bodrovis/lokex-cli/internal/viper_helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestLoadGenerateConfig_UsesDefaults(t *testing.T) {
	v := viper.New()

	v.Set(
		"manifest.generate.path",
		[]string{"./locales"},
	)

	v.Set(
		"manifest.generate.name-pattern",
		"{lang}/{name}.{ext}",
	)

	cmd := newGenerateTestCommand()

	cfg := &GenerateConfig{}

	err := LoadGenerateConfig(
		v,
		cmd,
		cfg,
	)
	require.NoError(t, err)

	require.Equal(
		t,
		[]string{"./locales"},
		*cfg.Paths,
	)

	require.Equal(
		t,
		"{lang}/{name}.{ext}",
		*cfg.NamePattern,
	)

	require.Equal(
		t,
		"{name}.{ext}",
		*cfg.FilenamePattern,
	)

	require.Equal(
		t,
		"./lokex-manifest.json",
		*cfg.Out,
	)
}

func TestLoadGenerateConfig_FlagsOverrideConfig(t *testing.T) {
	v := viper.New()

	v.Set(
		"manifest.generate.path",
		[]string{"./from-config"},
	)

	v.Set(
		"manifest.generate.name-pattern",
		"{lang}/{name}.{ext}",
	)

	v.Set(
		"manifest.generate.out",
		"./config.json",
	)

	cmd := newGenerateTestCommand()

	require.NoError(
		t,
		cmd.ParseFlags([]string{
			"--path=./from-cli",
			"--out=./cli.json",
		}),
	)

	cfg := &GenerateConfig{}

	err := LoadGenerateConfig(
		v,
		cmd,
		cfg,
	)
	require.NoError(t, err)

	require.Equal(
		t,
		[]string{"./from-cli"},
		*cfg.Paths,
	)

	require.Equal(
		t,
		"./cli.json",
		*cfg.Out,
	)

	// Untouched config value remains.
	require.Equal(
		t,
		"{lang}/{name}.{ext}",
		*cfg.NamePattern,
	)

	// Declared default.
	require.Equal(
		t,
		"{name}.{ext}",
		*cfg.FilenamePattern,
	)
}

func TestLoadGenerateConfig_EnvValues(t *testing.T) {
	t.Setenv(
		"LOKEX_MANIFEST_GENERATE_NAME_PATTERN",
		"{name}.{lang}.{ext}",
	)

	t.Setenv(
		"LOKEX_MANIFEST_GENERATE_OUT",
		"./from-env.json",
	)

	v := viper_helpers.NewConfigViper(
		"",
		"LOKEX",
	)

	v.SetConfigType("yaml")

	require.NoError(
		t,
		v.ReadConfig(
			strings.NewReader(`
manifest:
  generate:
    out: ./from-config.json
`),
		),
	)

	cmd := newGenerateTestCommand()

	cfg := &GenerateConfig{}

	err := LoadGenerateConfig(
		v,
		cmd,
		cfg,
	)
	require.NoError(t, err)

	require.Equal(
		t,
		"{name}.{lang}.{ext}",
		*cfg.NamePattern,
	)

	require.Equal(
		t,
		"./from-env.json",
		*cfg.Out,
	)
}

func TestLoadGenerateConfig_FlagOverridesEnvAndConfig(t *testing.T) {
	t.Setenv(
		"LOKEX_MANIFEST_GENERATE_OUT",
		"./from-env.json",
	)

	v := viper_helpers.NewConfigViper(
		"",
		"LOKEX",
	)

	v.Set(
		"manifest.generate.out",
		"./from-config.json",
	)

	cmd := newGenerateTestCommand()

	require.NoError(
		t,
		cmd.ParseFlags([]string{
			"--out=./from-cli.json",
		}),
	)

	cfg := &GenerateConfig{}

	err := LoadGenerateConfig(
		v,
		cmd,
		cfg,
	)
	require.NoError(t, err)

	require.Equal(
		t,
		"./from-cli.json",
		*cfg.Out,
	)
}

func newGenerateTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "generate",
	}

	params.BindFlags(
		cmd,
		generateParamSpecs,
	)

	return cmd
}
