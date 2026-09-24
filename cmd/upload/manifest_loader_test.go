package upload

import (
	"encoding/json/jsontext"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadManifestFile(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "manifest.json")

		content := `{
	  "items": [
	    {
	      "params": {
	        "filename": "locales/en.json",
	        "lang_iso": "en"
	      }
	    },
	    {
	      "params": {
	        "filename": "locales/de.json",
	        "lang_iso": "de"
	      },
	      "src_path": "./de.json"
	    }
	  ]
	}`

		require.NoError(
			t,
			os.WriteFile(path, []byte(content), 0o644),
		)

		got, err := loadManifestFile(path)
		require.NoError(t, err)

		require.Len(t, got.Items, 2)

		assert.Equal(
			t,
			"locales/en.json",
			got.Items[0].Params["filename"],
		)

		assert.Equal(
			t,
			"de",
			got.Items[1].Params["lang_iso"],
		)

		assert.Equal(
			t,
			"./de.json",
			got.Items[1].SrcPath,
		)
	})

	t.Run("preserves numeric params as jsontext.Value", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "manifest.json")

		content := `{
	  "items": [
	    {
	      "params": {
	        "filename": "locales/en.json",
	        "lang_iso": "en",
	        "filter_task_id": 1234567890123456789
	      },
	      "src_path": "./en.json"
	    }
	  ]
	}`

		require.NoError(
			t,
			os.WriteFile(path, []byte(content), 0o644),
		)

		got, err := loadManifestFile(path)
		require.NoError(t, err)
		require.Len(t, got.Items, 1)

		raw := got.Items[0].Params["filter_task_id"]

		num, ok := raw.(jsontext.Value)
		require.Truef(
			t,
			ok,
			"expected filter_task_id to be jsontext.Value, got %T (%#v)",
			raw,
			raw,
		)

		assert.Equal(t, jsontext.KindNumber, num.Kind())
		assert.Equal(t, "1234567890123456789", num.String())
	})

	t.Run("file does not exist", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "missing.json")

		_, err := loadManifestFile(path)

		require.Error(t, err)
		assert.Contains(t, err.Error(), `read manifest file "`)
		assert.Contains(t, err.Error(), "missing.json")
	})

	t.Run("invalid json", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "manifest.json")

		require.NoError(t, os.WriteFile(path, []byte(`{
	  "items": [
	    {
	      "params": {
	        "filename": "locales/en.json",
	        "lang_iso": "en"
	      }
	    }
	`), 0o644))

		_, err := loadManifestFile(path)

		require.Error(t, err)
		assert.Contains(t, err.Error(), `parse manifest file "`)
	})
}

func TestLoadManifestFile_NoItems(t *testing.T) {
	tests := map[string]string{
		"empty items": `{
			"items": []
		}`,
		"missing items field": `{
			"foo": "bar"
		}`,
	}

	for name, content := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "manifest.json")

			require.NoError(
				t,
				os.WriteFile(path, []byte(content), 0o644),
			)

			_, err := loadManifestFile(path)

			require.EqualError(
				t,
				err,
				fmt.Sprintf(
					"manifest file %q contains no items",
					path,
				),
			)
		})
	}
}
