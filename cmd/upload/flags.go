package upload

import (
	"time"

	params "github.com/bodrovis/lokex-cli/internal/params"
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
	params.BindFlags(cmd, flags, uploadParamSpecs)
}
