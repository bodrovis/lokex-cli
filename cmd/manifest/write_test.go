package manifest

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/bodrovis/lokex-cli/internal/uploadmanifest"
	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestWriteGeneratedManifest_Stdout(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
	}

	var out bytes.Buffer
	cmd.SetOut(&out)

	mf := uploadmanifest.File{
		Items: []uploadmanifest.Item{
			{
				SrcPath: "../locales/en/common.json",
				Params: lokexupload.UploadParams{
					"filename": "common.json",
					"lang_iso": "en",
				},
			},
		},
	}

	err := writeGeneratedManifest(
		cmd,
		"-",
		mf,
	)
	require.NoError(t, err)

	got := out.String()

	require.Contains(
		t,
		got,
		"\n  \"items\":",
	)

	require.Contains(
		t,
		got,
		"\"filename\": \"common.json\"",
	)

	require.Contains(
		t,
		got,
		"\"lang_iso\": \"en\"",
	)

	require.Contains(
		t,
		got,
		"\"src_path\": \"../locales/en/common.json\"",
	)

	require.NotContains(
		t,
		got,
		"Generated upload manifest:",
	)

	require.True(
		t,
		got[len(got)-1] == '\n',
	)
}

func TestWriteGeneratedManifest_File(t *testing.T) {
	dir := t.TempDir()

	outPath := filepath.Join(
		dir,
		"manifests",
		"upload.json",
	)

	cmd := &cobra.Command{
		Use: "test",
	}

	var out bytes.Buffer
	cmd.SetOut(&out)

	mf := uploadmanifest.File{
		Items: []uploadmanifest.Item{
			{
				SrcPath: "../locales/de/common.json",
				Params: lokexupload.UploadParams{
					"filename": "common.json",
					"lang_iso": "de",
				},
			},
		},
	}

	err := writeGeneratedManifest(
		cmd,
		outPath,
		mf,
	)
	require.NoError(t, err)

	require.FileExists(t, outPath)

	data, err := os.ReadFile(outPath)
	require.NoError(t, err)

	got := string(data)

	require.Contains(
		t,
		got,
		"\n  \"items\":",
	)

	require.Contains(
		t,
		got,
		"\"lang_iso\": \"de\"",
	)

	require.True(
		t,
		got[len(got)-1] == '\n',
	)

	require.Equal(
		t,
		"Generated upload manifest: "+outPath+" (1 items)\n",
		out.String(),
	)
}

func TestWriteGeneratedManifest_PrintsItemCount(t *testing.T) {
	dir := t.TempDir()

	outPath := filepath.Join(
		dir,
		"manifest.json",
	)

	cmd := &cobra.Command{
		Use: "test",
	}

	var out bytes.Buffer
	cmd.SetOut(&out)

	mf := uploadmanifest.File{
		Items: []uploadmanifest.Item{
			{
				SrcPath: "en.json",
				Params: lokexupload.UploadParams{
					"filename": "messages.json",
					"lang_iso": "en",
				},
			},
			{
				SrcPath: "de.json",
				Params: lokexupload.UploadParams{
					"filename": "messages.json",
					"lang_iso": "de",
				},
			},
		},
	}

	err := writeGeneratedManifest(
		cmd,
		outPath,
		mf,
	)
	require.NoError(t, err)

	require.Contains(
		t,
		out.String(),
		"(2 items)",
	)
}

type failingWriter struct {
	err error
}

func (w failingWriter) Write([]byte) (int, error) {
	return 0, w.err
}

func TestWriteGeneratedManifest_StdoutError(t *testing.T) {
	wantErr := errors.New("write failed")

	cmd := &cobra.Command{
		Use: "test",
	}

	cmd.SetOut(
		failingWriter{
			err: wantErr,
		},
	)

	mf := uploadmanifest.File{
		Items: []uploadmanifest.Item{
			{
				SrcPath: "en.json",
				Params: lokexupload.UploadParams{
					"filename": "messages.json",
					"lang_iso": "en",
				},
			},
		},
	}

	err := writeGeneratedManifest(
		cmd,
		"-",
		mf,
	)

	require.ErrorIs(t, err, wantErr)

	require.Contains(
		t,
		err.Error(),
		"write upload manifest to stdout",
	)
}
