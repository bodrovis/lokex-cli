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

func newFlags() *Flags {
	return &Flags{
		ContextTimeout: 150 * time.Second,
	}
}

func bindFlags(cmd *cobra.Command, flags *Flags) {
	cmd.Flags().SortFlags = false

	// Runtime / command behavior
	cmd.Flags().BoolVar(
		&flags.Poll,
		"poll",
		flags.Poll,
		"Wait until Lokalise finishes processing the upload",
	)

	cmd.Flags().DurationVar(
		&flags.ContextTimeout,
		"context-timeout",
		flags.ContextTimeout,
		"Overall command timeout (e.g. 30s, 2m). 0 disables the timeout",
	)

	cmd.Flags().StringVar(
		&flags.SrcPath,
		"src-path",
		flags.SrcPath,
		"Local path to read file contents from (optional)",
	)

	// Required / primary input
	cmd.Flags().StringVar(
		&flags.Filename,
		"filename",
		flags.Filename,
		"Filename sent to Lokalise (required)",
	)

	cmd.Flags().StringVar(
		&flags.LangISO,
		"lang-iso",
		flags.LangISO,
		"Language code of the translations in the file (required)",
	)

	// File source / input mode
	cmd.Flags().StringVar(
		&flags.Data,
		"data",
		flags.Data,
		"Base64-encoded file contents (optional; if set, file is not read from disk)",
	)

	cmd.Flags().StringVar(
		&flags.Format,
		"format",
		flags.Format,
		"File format (e.g. json, strings, xml)",
	)

	// Tagging / statuses
	cmd.Flags().StringSliceVar(
		&flags.Tags,
		"tags",
		flags.Tags,
		"Tags to apply to keys",
	)

	cmd.Flags().BoolVar(
		&flags.TagInsertedKeys,
		"tag-inserted-keys",
		flags.TagInsertedKeys,
		"Add tags to inserted keys",
	)

	cmd.Flags().BoolVar(
		&flags.TagUpdatedKeys,
		"tag-updated-keys",
		flags.TagUpdatedKeys,
		"Add tags to updated keys",
	)

	cmd.Flags().BoolVar(
		&flags.TagSkippedKeys,
		"tag-skipped-keys",
		flags.TagSkippedKeys,
		"Add tags to skipped keys",
	)

	cmd.Flags().StringSliceVar(
		&flags.CustomTranslationStatusIDs,
		"custom-translation-status-ids",
		flags.CustomTranslationStatusIDs,
		"Custom translation status IDs to add",
	)

	cmd.Flags().BoolVar(
		&flags.CustomTranslationStatusInsertedKeys,
		"custom-translation-status-inserted-keys",
		flags.CustomTranslationStatusInsertedKeys,
		"Add custom statuses to inserted keys",
	)

	cmd.Flags().BoolVar(
		&flags.CustomTranslationStatusUpdatedKeys,
		"custom-translation-status-updated-keys",
		flags.CustomTranslationStatusUpdatedKeys,
		"Add custom statuses to updated keys",
	)

	cmd.Flags().BoolVar(
		&flags.CustomTranslationStatusSkippedKeys,
		"custom-translation-status-skipped-keys",
		flags.CustomTranslationStatusSkippedKeys,
		"Add custom statuses to skipped keys",
	)

	// Import behavior
	cmd.Flags().BoolVar(
		&flags.ConvertPlaceholders,
		"convert-placeholders",
		flags.ConvertPlaceholders,
		"Convert placeholders to Lokalise universal placeholders",
	)

	cmd.Flags().BoolVar(
		&flags.DetectICUPlurals,
		"detect-icu-plurals",
		flags.DetectICUPlurals,
		"Automatically detect and parse ICU plurals",
	)

	cmd.Flags().BoolVar(
		&flags.ReplaceModified,
		"replace-modified",
		flags.ReplaceModified,
		"Replace modified translations from the uploaded file",
	)

	cmd.Flags().BoolVar(
		&flags.SlashNToLinebreak,
		"slashn-to-linebreak",
		flags.SlashNToLinebreak,
		"Replace \\n with a real line break",
	)

	cmd.Flags().BoolVar(
		&flags.KeysToValues,
		"keys-to-values",
		flags.KeysToValues,
		"Replace values with key names",
	)

	cmd.Flags().BoolVar(
		&flags.DistinguishByFile,
		"distinguish-by-file",
		flags.DistinguishByFile,
		"Allow same key names to coexist across different filenames",
	)

	cmd.Flags().BoolVar(
		&flags.ApplyTM,
		"apply-tm",
		flags.ApplyTM,
		"Apply 100% translation memory matches",
	)

	cmd.Flags().BoolVar(
		&flags.UseAutomations,
		"use-automations",
		flags.UseAutomations,
		"Run automations for this upload",
	)

	cmd.Flags().BoolVar(
		&flags.HiddenFromContributors,
		"hidden-from-contributors",
		flags.HiddenFromContributors,
		"Mark newly created keys as hidden from contributors",
	)

	cmd.Flags().BoolVar(
		&flags.CleanupMode,
		"cleanup-mode",
		flags.CleanupMode,
		"Delete keys/translations not present in the uploaded file",
	)

	cmd.Flags().BoolVar(
		&flags.SkipDetectLangISO,
		"skip-detect-lang-iso",
		flags.SkipDetectLangISO,
		"Skip automatic language detection by filename",
	)

	// Task-related / advanced
	cmd.Flags().Int64Var(
		&flags.FilterTaskID,
		"filter-task-id",
		flags.FilterTaskID,
		"Apply import results as a part of a task (offline_xliff only)",
	)
}
