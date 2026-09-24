package uploadmanifest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiscoverFiles_SingleFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	path := filepath.Join(dir, "messages.json")

	require.NoError(
		t,
		os.WriteFile(
			path,
			[]byte(`{}`),
			0o644,
		),
	)

	got, err := DiscoverFiles(
		[]string{path},
	)

	require.NoError(t, err)

	require.Equal(
		t,
		[]DiscoveredFile{
			{
				Path:         path,
				RelativePath: "messages.json",
			},
		},
		got,
	)
}

func TestDiscoverFiles_DirectoryRecursively(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	enDir := filepath.Join(root, "en")
	deDir := filepath.Join(root, "de", "nested")

	require.NoError(
		t,
		os.MkdirAll(enDir, 0o755),
	)

	require.NoError(
		t,
		os.MkdirAll(deDir, 0o755),
	)

	enFile := filepath.Join(
		enDir,
		"common.json",
	)

	deFile := filepath.Join(
		deDir,
		"messages.json",
	)

	require.NoError(
		t,
		os.WriteFile(
			enFile,
			[]byte(`{}`),
			0o644,
		),
	)

	require.NoError(
		t,
		os.WriteFile(
			deFile,
			[]byte(`{}`),
			0o644,
		),
	)

	got, err := DiscoverFiles(
		[]string{root},
	)

	require.NoError(t, err)

	require.Equal(
		t,
		[]DiscoveredFile{
			{
				Path:         deFile,
				RelativePath: "de/nested/messages.json",
			},
			{
				Path:         enFile,
				RelativePath: "en/common.json",
			},
		},
		got,
	)
}

func TestDiscoverFiles_MultiplePathsAreSorted(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	first := filepath.Join(
		root,
		"a.json",
	)

	second := filepath.Join(
		root,
		"z.json",
	)

	require.NoError(
		t,
		os.WriteFile(
			first,
			[]byte(`{}`),
			0o644,
		),
	)

	require.NoError(
		t,
		os.WriteFile(
			second,
			[]byte(`{}`),
			0o644,
		),
	)

	got, err := DiscoverFiles(
		[]string{
			second,
			first,
		},
	)

	require.NoError(t, err)

	require.Equal(
		t,
		[]DiscoveredFile{
			{
				Path:         first,
				RelativePath: "a.json",
			},
			{
				Path:         second,
				RelativePath: "z.json",
			},
		},
		got,
	)
}

func TestDiscoverFiles_MissingPath(t *testing.T) {
	t.Parallel()

	path := filepath.Join(
		t.TempDir(),
		"does-not-exist",
	)

	got, err := DiscoverFiles(
		[]string{path},
	)

	require.Nil(t, got)

	require.Error(t, err)

	require.Contains(
		t,
		err.Error(),
		"inspect path",
	)

	require.Contains(
		t,
		err.Error(),
		path,
	)
}

func TestDiscoverFiles_NoPaths(t *testing.T) {
	t.Parallel()

	got, err := DiscoverFiles(nil)

	require.NoError(t, err)
	require.Empty(t, got)
}
