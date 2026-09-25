package manifest

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/bodrovis/lokex-cli/internal/appstate"
	"github.com/bodrovis/lokex-cli/internal/uploadmanifest"
	"github.com/spf13/cobra"
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

func TestRunGenerateCommand(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	sourceDir := filepath.Join(
		root,
		"locales",
		"en",
	)

	require.NoError(
		t,
		os.MkdirAll(
			sourceDir,
			0o755,
		),
	)

	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(
				sourceDir,
				"common.json",
			),
			[]byte(`{}`),
			0o644,
		),
	)

	outPath := filepath.Join(
		root,
		"manifests",
		"upload.json",
	)

	cmd := &cobra.Command{
		Use: "test",
	}

	var stdout bytes.Buffer
	cmd.SetOut(&stdout)

	cfg := &GenerateConfig{
		Paths: new(
			[]string{
				filepath.Join(
					root,
					"locales",
				),
			},
		),
		NamePattern: new(
			"{lang}/{name}.{ext}",
		),
		FilenamePattern: new(
			"{name}.{ext}",
		),
		LanguageMapping: new(
			`[
				{
					"original_language_iso": "en",
					"custom_language_iso": "en-US"
				}
			]`,
		),
		Out: new(outPath),
	}

	err := runGenerateCommand(
		cmd,
		cfg,
	)
	require.NoError(t, err)

	require.FileExists(
		t,
		outPath,
	)

	mf, err := uploadmanifest.Load(
		outPath,
	)
	require.NoError(t, err)
	require.Len(t, mf.Items, 1)

	require.Equal(
		t,
		"common.json",
		mf.Items[0].Params["filename"],
	)

	require.Equal(
		t,
		"en-US",
		mf.Items[0].Params["lang_iso"],
	)

	require.Equal(
		t,
		filepath.ToSlash(
			filepath.Join(
				"..",
				"locales",
				"en",
				"common.json",
			),
		),
		mf.Items[0].SrcPath,
	)

	require.Equal(
		t,
		"Generated upload manifest: "+outPath+" (1 items)\n",
		stdout.String(),
	)
}

func TestRunGenerateCommand_Stdout(t *testing.T) {
	root := t.TempDir()

	t.Chdir(root)

	sourceDir := filepath.Join(
		root,
		"de",
	)

	require.NoError(
		t,
		os.MkdirAll(
			sourceDir,
			0o755,
		),
	)

	require.NoError(
		t,
		os.WriteFile(
			filepath.Join(
				sourceDir,
				"messages.json",
			),
			[]byte(`{}`),
			0o644,
		),
	)

	cmd := &cobra.Command{
		Use: "test",
	}

	var stdout bytes.Buffer
	cmd.SetOut(&stdout)

	cfg := &GenerateConfig{
		Paths: new(
			[]string{
				root,
			},
		),
		NamePattern: new(
			"{lang}/{name}.{ext}",
		),
		FilenamePattern: new(
			"{name}.{ext}",
		),
		Out: new("-"),
	}

	err := runGenerateCommand(
		cmd,
		cfg,
	)
	require.NoError(t, err)

	require.Contains(
		t,
		stdout.String(),
		`"filename": "messages.json"`,
	)

	require.Contains(
		t,
		stdout.String(),
		`"lang_iso": "de"`,
	)

	require.Contains(
		t,
		stdout.String(),
		`"src_path": "de/messages.json"`,
	)

	require.NotContains(
		t,
		stdout.String(),
		"Generated upload manifest:",
	)
}

func TestRunGenerateCommand_InvalidLanguageMapping(
	t *testing.T,
) {
	t.Parallel()

	cmd := &cobra.Command{
		Use: "test",
	}

	cfg := &GenerateConfig{
		LanguageMapping: new(`{`),
		Out:             new("-"),
	}

	err := runGenerateCommand(
		cmd,
		cfg,
	)

	require.Error(t, err)

	require.Contains(
		t,
		err.Error(),
		"parse language-mapping",
	)
}

func TestRunGenerateCommand_GenerateError(
	t *testing.T,
) {
	t.Parallel()

	cmd := &cobra.Command{
		Use: "test",
	}

	cfg := &GenerateConfig{
		Paths: new(
			[]string{
				filepath.Join(
					t.TempDir(),
					"missing",
				),
			},
		),
		NamePattern: new(
			"{lang}.{ext}",
		),
		FilenamePattern: new(
			"{lang}.{ext}",
		),
		Out: new("-"),
	}

	err := runGenerateCommand(
		cmd,
		cfg,
	)

	require.Error(t, err)

	require.Contains(
		t,
		err.Error(),
		"inspect path",
	)
}
