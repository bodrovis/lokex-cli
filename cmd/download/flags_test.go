package download

import (
	"testing"
	"time"

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

	if flags.ContextTimeout != 150*time.Second {
		t.Fatalf("unexpected ContextTimeout: got %v, want %v", flags.ContextTimeout, 150*time.Second)
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
		"--context-timeout=45s",
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
	if flags.ContextTimeout != 45*time.Second {
		t.Fatalf("unexpected ContextTimeout: got %v", flags.ContextTimeout)
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
	if flags.ContextTimeout != 150*time.Second {
		t.Fatalf("unexpected ContextTimeout: got %v, want %v", flags.ContextTimeout, 150*time.Second)
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
