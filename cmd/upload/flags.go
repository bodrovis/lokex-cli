package upload

import (
	"time"

	"github.com/spf13/cobra"
)

type Flags struct {
	Filename       string
	SrcPath        string
	Data           string
	LangISO        string
	Poll           bool
	ContextTimeout time.Duration

	ConvertPlaceholders                 bool
	DetectICUPlurals                    bool
	Tags                                []string
	TagInsertedKeys                     bool
	TagUpdatedKeys                      bool
	TagSkippedKeys                      bool
	ReplaceModified                     bool
	SlashNToLinebreak                   bool
	KeysToValues                        bool
	DistinguishByFile                   bool
	ApplyTM                             bool
	UseAutomations                      bool
	HiddenFromContributors              bool
	CleanupMode                         bool
	CustomTranslationStatusIDs          []string
	CustomTranslationStatusInsertedKeys bool
	CustomTranslationStatusUpdatedKeys  bool
	CustomTranslationStatusSkippedKeys  bool
	SkipDetectLangISO                   bool
	Format                              string
	FilterTaskID                        int64
}

func bindFlags(cmd *cobra.Command, flags *Flags) {
	cmd.Flags().SortFlags = false

	// Runtime / command behavior
	cmd.Flags().BoolVar(
		&flags.Poll,
		"poll",
		false,
		"Wait until Lokalise finishes processing the upload",
	)

	cmd.Flags().DurationVar(
		&flags.ContextTimeout,
		"context-timeout",
		150*time.Second,
		"Overall command timeout (e.g. 30s, 2m). 0 disables the timeout",
	)

	cmd.Flags().StringVar(
		&flags.SrcPath,
		"src-path",
		"",
		"Local path to read file contents from (optional)",
	)

	// Required / primary input
	cmd.Flags().StringVar(
		&flags.Filename,
		"filename",
		"",
		"Filename sent to Lokalise (required)",
	)

	cmd.Flags().StringVar(
		&flags.LangISO,
		"lang-iso",
		"",
		"Language code of the translations in the file (required)",
	)

	// File source / input mode
	cmd.Flags().StringVar(
		&flags.Data,
		"data",
		"",
		"Base64-encoded file contents (optional; if set, file is not read from disk)",
	)

	cmd.Flags().StringVar(
		&flags.Format,
		"format",
		"",
		"File format (e.g. json, strings, xml)",
	)

	// Tagging / statuses
	cmd.Flags().StringSliceVar(
		&flags.Tags,
		"tags",
		nil,
		"Tags to apply to keys",
	)

	cmd.Flags().BoolVar(
		&flags.TagInsertedKeys,
		"tag-inserted-keys",
		false,
		"Add tags to inserted keys",
	)

	cmd.Flags().BoolVar(
		&flags.TagUpdatedKeys,
		"tag-updated-keys",
		false,
		"Add tags to updated keys",
	)

	cmd.Flags().BoolVar(
		&flags.TagSkippedKeys,
		"tag-skipped-keys",
		false,
		"Add tags to skipped keys",
	)

	cmd.Flags().StringSliceVar(
		&flags.CustomTranslationStatusIDs,
		"custom-translation-status-ids",
		nil,
		"Custom translation status IDs to add",
	)

	cmd.Flags().BoolVar(
		&flags.CustomTranslationStatusInsertedKeys,
		"custom-translation-status-inserted-keys",
		false,
		"Add custom statuses to inserted keys",
	)

	cmd.Flags().BoolVar(
		&flags.CustomTranslationStatusUpdatedKeys,
		"custom-translation-status-updated-keys",
		false,
		"Add custom statuses to updated keys",
	)

	cmd.Flags().BoolVar(
		&flags.CustomTranslationStatusSkippedKeys,
		"custom-translation-status-skipped-keys",
		false,
		"Add custom statuses to skipped keys",
	)

	// Import behavior
	cmd.Flags().BoolVar(
		&flags.ConvertPlaceholders,
		"convert-placeholders",
		false,
		"Convert placeholders to Lokalise universal placeholders",
	)

	cmd.Flags().BoolVar(
		&flags.DetectICUPlurals,
		"detect-icu-plurals",
		false,
		"Automatically detect and parse ICU plurals",
	)

	cmd.Flags().BoolVar(
		&flags.ReplaceModified,
		"replace-modified",
		false,
		"Replace modified translations from the uploaded file",
	)

	cmd.Flags().BoolVar(
		&flags.SlashNToLinebreak,
		"slashn-to-linebreak",
		false,
		"Replace \\n with a real line break",
	)

	cmd.Flags().BoolVar(
		&flags.KeysToValues,
		"keys-to-values",
		false,
		"Replace values with key names",
	)

	cmd.Flags().BoolVar(
		&flags.DistinguishByFile,
		"distinguish-by-file",
		false,
		"Allow same key names to coexist across different filenames",
	)

	cmd.Flags().BoolVar(
		&flags.ApplyTM,
		"apply-tm",
		false,
		"Apply 100% translation memory matches",
	)

	cmd.Flags().BoolVar(
		&flags.UseAutomations,
		"use-automations",
		false,
		"Run automations for this upload",
	)

	cmd.Flags().BoolVar(
		&flags.HiddenFromContributors,
		"hidden-from-contributors",
		false,
		"Mark newly created keys as hidden from contributors",
	)

	cmd.Flags().BoolVar(
		&flags.CleanupMode,
		"cleanup-mode",
		false,
		"Delete keys/translations not present in the uploaded file",
	)

	cmd.Flags().BoolVar(
		&flags.SkipDetectLangISO,
		"skip-detect-lang-iso",
		false,
		"Skip automatic language detection by filename",
	)

	// Task-related / advanced
	cmd.Flags().Int64Var(
		&flags.FilterTaskID,
		"filter-task-id",
		0,
		"Apply import results as a part of a task (offline_xliff only)",
	)
}
