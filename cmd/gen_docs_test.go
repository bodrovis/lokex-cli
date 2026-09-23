package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestNewGenDocsCmd(t *testing.T) {
	root := &cobra.Command{
		Use: "lokex",
	}

	cmd := newGenDocsCmd(root)

	require.NotNil(t, cmd)
	require.Equal(t, "gendocs", cmd.Use)
	require.True(t, cmd.Hidden)
	require.NotNil(t, cmd.RunE)
}

func TestGenerateDocs(t *testing.T) {
	root := &cobra.Command{
		Use:   "lokex",
		Short: "Test CLI",
	}

	root.AddCommand(&cobra.Command{
		Use:   "upload",
		Short: "Upload files",
		Run: func(*cobra.Command, []string) {
		},
	})

	dir := filepath.Join(
		t.TempDir(),
		"docs",
	)

	err := generateDocs(root, dir)
	require.NoError(t, err)

	require.DirExists(t, dir)

	require.FileExists(
		t,
		filepath.Join(dir, "lokex.md"),
	)

	require.FileExists(
		t,
		filepath.Join(dir, "lokex_upload.md"),
	)
}

func TestGenerateDocs_ReturnsMkdirError(t *testing.T) {
	dir := t.TempDir()

	file := filepath.Join(
		dir,
		"not-a-directory",
	)

	require.NoError(
		t,
		os.WriteFile(
			file,
			[]byte("hello"),
			0o644,
		),
	)

	target := filepath.Join(
		file,
		"docs",
	)

	root := &cobra.Command{
		Use: "lokex",
	}

	err := generateDocs(
		root,
		target,
	)

	require.Error(t, err)
}
