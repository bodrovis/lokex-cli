package uploadmanifest

import (
	"encoding/json/jsontext"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	path := filepath.Join(
		dir,
		"manifest.json",
	)

	require.NoError(
		t,
		os.WriteFile(
			path,
			[]byte(`{
				"items": [
					{
						"src_path": "../locales/en/common.json",
						"params": {
							"filename": "common.json",
							"lang_iso": "en"
						}
					}
				]
			}`),
			0o644,
		),
	)

	mf, err := Load(path)

	require.NoError(t, err)
	require.Len(t, mf.Items, 1)

	require.Equal(
		t,
		"../locales/en/common.json",
		mf.Items[0].SrcPath,
	)

	require.Equal(
		t,
		"common.json",
		mf.Items[0].Params["filename"],
	)

	require.Equal(
		t,
		"en",
		mf.Items[0].Params["lang_iso"],
	)
}

func TestLoad_PreservesJSONNumbers(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	path := filepath.Join(
		dir,
		"manifest.json",
	)

	require.NoError(
		t,
		os.WriteFile(
			path,
			[]byte(`{
				"items": [
					{
						"params": {
							"filename": "common.json",
							"lang_iso": "en",
							"filter_task_id": 9223372036854775807
						}
					}
				]
			}`),
			0o644,
		),
	)

	mf, err := Load(path)

	require.NoError(t, err)
	require.Len(t, mf.Items, 1)

	got := mf.Items[0].Params["filter_task_id"]

	require.Equal(
		t,
		jsontext.Value("9223372036854775807"),
		got,
	)
}

func TestLoad_MissingFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(
		t.TempDir(),
		"missing.json",
	)

	mf, err := Load(path)

	require.Equal(
		t,
		File{},
		mf,
	)

	require.Error(t, err)
	require.Contains(
		t,
		err.Error(),
		"read manifest file",
	)

	require.Contains(
		t,
		err.Error(),
		path,
	)
}

func TestLoad_InvalidJSON(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	path := filepath.Join(
		dir,
		"manifest.json",
	)

	require.NoError(
		t,
		os.WriteFile(
			path,
			[]byte(`{`),
			0o644,
		),
	)

	mf, err := Load(path)

	require.Equal(
		t,
		File{},
		mf,
	)

	require.Error(t, err)
	require.Contains(
		t,
		err.Error(),
		"parse manifest file",
	)
}

func TestLoad_EmptyItems(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	path := filepath.Join(
		dir,
		"manifest.json",
	)

	require.NoError(
		t,
		os.WriteFile(
			path,
			[]byte(`{
				"items": []
			}`),
			0o644,
		),
	)

	mf, err := Load(path)

	require.Equal(
		t,
		File{},
		mf,
	)

	require.EqualError(
		t,
		err,
		fmt.Sprintf(
			"manifest file %q contains no items",
			path,
		),
	)
}

func TestLoad_MissingItems(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	path := filepath.Join(
		dir,
		"manifest.json",
	)

	require.NoError(
		t,
		os.WriteFile(
			path,
			[]byte(`{}`),
			0o644,
		),
	)

	mf, err := Load(path)

	require.Equal(
		t,
		File{},
		mf,
	)

	require.EqualError(
		t,
		err,
		fmt.Sprintf(
			"manifest file %q contains no items",
			path,
		),
	)
}
