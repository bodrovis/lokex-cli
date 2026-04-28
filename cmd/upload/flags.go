package upload

import (
	params "github.com/bodrovis/lokex-cli/internal/params"
	"github.com/spf13/cobra"
)

type Flags struct {
	Filename string
	SrcPath  string
	Data     string
	LangISO  string
	Poll     bool

	Manifest string

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
	return &Flags{}
}

func bindFlags(cmd *cobra.Command, flags *Flags) {
	params.BindFlags(cmd, flags, uploadParamSpecs)
}
