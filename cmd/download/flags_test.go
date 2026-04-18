package download

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestNewFlags_Defaults(t *testing.T) {
	t.Parallel()

	flags := newFlags()
	if flags == nil {
		t.Fatal("expected non-nil flags")
	}

	if flags.Out != "./locales" {
		t.Fatalf("unexpected Out: got %q, want %q", flags.Out, "./locales")
	}

	if flags.Format != "" {
		t.Fatalf("expected empty Format, got %q", flags.Format)
	}

	if flags.Async {
		t.Fatal("expected Async to be false by default")
	}
}

func TestBindFlags_BindsRepresentativeValues(t *testing.T) {
	t.Parallel()

	cmd := &cobra.Command{Use: "test"}
	flags := newFlags()

	bindFlags(cmd, flags)

	err := cmd.Flags().Parse([]string{
		"--out=/tmp/lokex",
		"--format=json",
		"--async",
		"--original-filenames",
		"--bundle-structure=%LANG_ISO%.json",
		"--filter-langs=en,fr",
		"--include-comments",
		"--export-sort=first_added",
		"--language-mapping=[{\"lang_iso\":\"en\",\"custom_iso\":\"en_US\"}]",
		"--filter-task-id=123",
	})
	if err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	if flags.Out != "/tmp/lokex" {
		t.Fatalf("unexpected Out: got %q", flags.Out)
	}
	if flags.Format != "json" {
		t.Fatalf("unexpected Format: got %q", flags.Format)
	}
	if !flags.Async {
		t.Fatal("expected Async to be true")
	}
	if !flags.OriginalFilenames {
		t.Fatal("expected OriginalFilenames to be true")
	}
	if flags.BundleStructure != "%LANG_ISO%.json" {
		t.Fatalf("unexpected BundleStructure: got %q", flags.BundleStructure)
	}
	if len(flags.FilterLangs) != 2 || flags.FilterLangs[0] != "en" || flags.FilterLangs[1] != "fr" {
		t.Fatalf("unexpected FilterLangs: got %#v", flags.FilterLangs)
	}
	if !flags.IncludeComments {
		t.Fatal("expected IncludeComments to be true")
	}
	if flags.ExportSort != "first_added" {
		t.Fatalf("unexpected ExportSort: got %q", flags.ExportSort)
	}
	if flags.LanguageMappingJSON != `[{"lang_iso":"en","custom_iso":"en_US"}]` {
		t.Fatalf("unexpected LanguageMappingJSON: got %q", flags.LanguageMappingJSON)
	}
	if flags.FilterTaskID != 123 {
		t.Fatalf("unexpected FilterTaskID: got %d", flags.FilterTaskID)
	}
}

func TestBindFlags_PreservesDefaultsWhenNoArgs(t *testing.T) {
	t.Parallel()

	cmd := &cobra.Command{Use: "test"}
	flags := newFlags()

	bindFlags(cmd, flags)

	if err := cmd.Flags().Parse(nil); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	if flags.Out != "./locales" {
		t.Fatalf("unexpected Out: got %q, want %q", flags.Out, "./locales")
	}
	if flags.Format != "" {
		t.Fatalf("expected empty Format, got %q", flags.Format)
	}
	if flags.Async {
		t.Fatal("expected Async to remain false")
	}
	if flags.FilterTaskID != 0 {
		t.Fatalf("expected FilterTaskID to remain 0, got %d", flags.FilterTaskID)
	}
	if flags.FilterLangs != nil {
		t.Fatalf("expected FilterLangs to remain nil, got %#v", flags.FilterLangs)
	}
}

func TestBindFlags_DisablesFlagSorting(t *testing.T) {
	t.Parallel()

	cmd := &cobra.Command{Use: "test"}
	flags := newFlags()

	cmd.Flags().SortFlags = true
	bindFlags(cmd, flags)

	if cmd.Flags().SortFlags {
		t.Fatal("expected SortFlags to be false")
	}
}

func TestBindFlags_BoolCanBeExplicitlySetFalse(t *testing.T) {
	t.Parallel()

	cmd := &cobra.Command{Use: "test"}
	flags := newFlags()

	bindFlags(cmd, flags)

	if err := cmd.Flags().Parse([]string{"--compact=false"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	if flags.Compact {
		t.Fatal("expected Compact to be false")
	}

	if !cmd.Flags().Changed("compact") {
		t.Fatal("expected compact flag to be marked as changed")
	}
}

func TestBindFlags_MarksRepresentativeFlagsAsChanged(t *testing.T) {
	t.Parallel()

	cmd := &cobra.Command{Use: "test"}
	flags := newFlags()

	bindFlags(cmd, flags)

	if err := cmd.Flags().Parse([]string{
		"--format=json",
		"--filter-langs=en,fr",
		"--filter-task-id=123",
	}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	for _, name := range []string{"format", "filter-langs", "filter-task-id"} {
		if !cmd.Flags().Changed(name) {
			t.Fatalf("expected flag %q to be marked as changed", name)
		}
	}
}

func TestBindFlags_BindsInt64ValueBeyondInt32Range(t *testing.T) {
	t.Parallel()

	cmd := &cobra.Command{Use: "test"}
	flags := newFlags()

	bindFlags(cmd, flags)

	const want int64 = 5000000000

	if err := cmd.Flags().Parse([]string{"--filter-task-id=5000000000"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	if flags.FilterTaskID != want {
		t.Fatalf("unexpected FilterTaskID: got %d, want %d", flags.FilterTaskID, want)
	}

	if !cmd.Flags().Changed("filter-task-id") {
		t.Fatal("expected filter-task-id to be marked as changed")
	}
}

func TestBindFlags_BindsEverySpecFlag(t *testing.T) {
	t.Parallel()

	cmd := &cobra.Command{Use: "test"}
	flags := newFlags()

	bindFlags(cmd, flags)

	for _, spec := range downloadParamSpecs {
		if spec.FlagName == "" {
			t.Fatal("encountered spec with empty FlagName")
		}

		flag := cmd.Flags().Lookup(spec.FlagName)
		if flag == nil {
			t.Fatalf("expected flag %q to be bound from spec", spec.FlagName)
		}
	}
}

func TestBindFlags_RepresentativeFlagDefaultsMatchNewFlags(t *testing.T) {
	t.Parallel()

	cmd := &cobra.Command{Use: "test"}
	flags := newFlags()

	bindFlags(cmd, flags)

	tests := []struct {
		name string
		want string
	}{
		{name: "out", want: "./locales"},
		{name: "format", want: ""},
		{name: "bundle-structure", want: ""},
		{name: "language-mapping", want: ""},
	}

	for _, tt := range tests {
		flag := cmd.Flags().Lookup(tt.name)
		if flag == nil {
			t.Fatalf("expected flag %q to exist", tt.name)
		}
		if flag.DefValue != tt.want {
			t.Fatalf("unexpected default for %q: got %q, want %q", tt.name, flag.DefValue, tt.want)
		}
	}
}
