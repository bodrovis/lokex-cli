package upload

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
	"github.com/spf13/cobra"
)

func TestPrintBatchItemResult(t *testing.T) {
	t.Run("prints failed item", func(t *testing.T) {
		cmd := &cobra.Command{}
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		item := lokexupload.BatchUploadResultItem{
			Index:   2,
			SrcPath: "./locales/de.json",
			Err:     errors.New("boom"),
		}

		printBatchItemResult(cmd, item, false)

		got := out.String()
		want := "Upload failed: index=2 src=\"./locales/de.json\" err=boom\n"
		if got != want {
			t.Fatalf("unexpected output:\n got: %q\nwant: %q", got, want)
		}
	})

	t.Run("prints started item when poll is false", func(t *testing.T) {
		cmd := &cobra.Command{}
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		item := lokexupload.BatchUploadResultItem{
			Index:     0,
			SrcPath:   "./locales/en.json",
			ProcessID: "process-123",
		}

		printBatchItemResult(cmd, item, false)

		got := out.String()
		want := "Upload started: index=0 src=\"./locales/en.json\" process_id=process-123\n"
		if got != want {
			t.Fatalf("unexpected output:\n got: %q\nwant: %q", got, want)
		}
	})

	t.Run("prints completed item when poll is true", func(t *testing.T) {
		cmd := &cobra.Command{}
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		item := lokexupload.BatchUploadResultItem{
			Index:     1,
			SrcPath:   "./locales/fr.json",
			ProcessID: "process-456",
		}

		printBatchItemResult(cmd, item, true)

		got := out.String()
		want := "Upload completed: index=1 src=\"./locales/fr.json\" process_id=process-456\n"
		if got != want {
			t.Fatalf("unexpected output:\n got: %q\nwant: %q", got, want)
		}
	})
}

func TestPrintBatchSummary(t *testing.T) {
	cmd := &cobra.Command{}
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	printBatchSummary(cmd, 3, 2, 1)

	got := out.String()
	want := "Batch summary: total=3 success=2 failed=1\n"
	if got != want {
		t.Fatalf("unexpected output:\n got: %q\nwant: %q", got, want)
	}
}

func TestPrintBatchUploadResult(t *testing.T) {
	t.Run("prints started items and summary when poll is false", func(t *testing.T) {
		cmd := &cobra.Command{}
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		result := lokexupload.BatchUploadResult{
			Items: []lokexupload.BatchUploadResultItem{
				{
					Index:     0,
					SrcPath:   "./locales/en.json",
					ProcessID: "process-1",
				},
				{
					Index:   1,
					SrcPath: "./locales/de.json",
					Err:     errors.New("failed upload"),
				},
			},
		}

		printBatchUploadResult(cmd, result, false)

		got := out.String()

		if !strings.Contains(got, `Upload started: index=0 src="./locales/en.json" process_id=process-1`) {
			t.Fatalf("unexpected output: %q", got)
		}
		if !strings.Contains(got, `Upload failed: index=1 src="./locales/de.json" err=failed upload`) {
			t.Fatalf("unexpected output: %q", got)
		}
		if !strings.Contains(got, "Batch summary: total=2 success=1 failed=1") {
			t.Fatalf("unexpected output: %q", got)
		}
	})

	t.Run("prints completed items and summary when poll is true", func(t *testing.T) {
		cmd := &cobra.Command{}
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		result := lokexupload.BatchUploadResult{
			Items: []lokexupload.BatchUploadResultItem{
				{
					Index:     0,
					SrcPath:   "./locales/en.json",
					ProcessID: "process-1",
				},
				{
					Index:     1,
					SrcPath:   "./locales/fr.json",
					ProcessID: "process-2",
				},
			},
		}

		printBatchUploadResult(cmd, result, true)

		got := out.String()

		if !strings.Contains(got, `Upload completed: index=0 src="./locales/en.json" process_id=process-1`) {
			t.Fatalf("unexpected output: %q", got)
		}
		if !strings.Contains(got, `Upload completed: index=1 src="./locales/fr.json" process_id=process-2`) {
			t.Fatalf("unexpected output: %q", got)
		}
		if !strings.Contains(got, "Batch summary: total=2 success=2 failed=0") {
			t.Fatalf("unexpected output: %q", got)
		}
	})

	t.Run("prints summary for empty result", func(t *testing.T) {
		cmd := &cobra.Command{}
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)

		result := lokexupload.BatchUploadResult{}

		printBatchUploadResult(cmd, result, false)

		got := out.String()
		want := "Batch summary: total=0 success=0 failed=0\n"
		if got != want {
			t.Fatalf("unexpected output:\n got: %q\nwant: %q", got, want)
		}
	})
}
