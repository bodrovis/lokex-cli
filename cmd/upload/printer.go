package upload

import (
	"strings"

	"github.com/spf13/cobra"

	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
)

func printUploadResult(cmd *cobra.Command, result string, poll bool) {
	result = strings.TrimSpace(result)

	if result == "" {
		if poll {
			cmd.Println("Upload completed (process ID unknown)")
			return
		}
		cmd.Println("Upload started (process ID unknown)")
		return
	}

	if poll {
		cmd.Printf("Upload completed: %s\n", result)
		return
	}

	cmd.Printf("Upload started: %s\n", result)
}

func printBatchUploadResult(
	cmd *cobra.Command,
	result lokexupload.BatchUploadResult,
	poll bool,
) {
	var successCount int
	var failedCount int

	for _, item := range result.Items {
		if item.Err != nil {
			failedCount++
		} else {
			successCount++
		}

		printBatchItemResult(cmd, item, poll)
	}

	printBatchSummary(
		cmd,
		len(result.Items),
		successCount,
		failedCount,
	)
}

func printBatchItemResult(
	cmd *cobra.Command,
	item lokexupload.BatchUploadResultItem,
	poll bool,
) {
	if item.Err != nil {
		cmd.Printf(
			"Upload failed: index=%d src=%q err=%v\n",
			item.Index,
			item.SrcPath,
			item.Err,
		)

		return
	}

	if poll {
		cmd.Printf(
			"Upload completed: index=%d src=%q process_id=%s\n",
			item.Index,
			item.SrcPath,
			item.ProcessID,
		)

		return
	}

	cmd.Printf(
		"Upload started: index=%d src=%q process_id=%s\n",
		item.Index,
		item.SrcPath,
		item.ProcessID,
	)
}

func printBatchSummary(cmd *cobra.Command, total, success, failed int) {
	cmd.Printf(
		"Batch summary: total=%d success=%d failed=%d\n",
		total,
		success,
		failed,
	)
}
