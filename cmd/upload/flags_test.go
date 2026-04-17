package upload

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

	if flags.ContextTimeout != 150*time.Second {
		t.Fatalf("unexpected ContextTimeout: got %v, want %v", flags.ContextTimeout, 150*time.Second)
	}

	if flags.Filename != "" {
		t.Fatalf("expected empty Filename, got %q", flags.Filename)
	}
	if flags.SrcPath != "" {
		t.Fatalf("expected empty SrcPath, got %q", flags.SrcPath)
	}
	if flags.Data != "" {
		t.Fatalf("expected empty Data, got %q", flags.Data)
	}
	if flags.LangISO != "" {
		t.Fatalf("expected empty LangISO, got %q", flags.LangISO)
	}
	if flags.Poll {
		t.Fatal("expected Poll to be false")
	}
	if flags.Format != "" {
		t.Fatalf("expected empty Format, got %q", flags.Format)
	}
	if flags.FilterTaskID != 0 {
		t.Fatalf("expected zero FilterTaskID, got %d", flags.FilterTaskID)
	}
}

func TestBindFlags_BindsRepresentativeValues(t *testing.T) {
	t.Parallel()

	cmd := &cobra.Command{Use: "test"}
	flags := newFlags()

	bindFlags(cmd, flags)

	err := cmd.Flags().Parse([]string{
		"--poll",
		"--context-timeout=45s",
		"--src-path=./locales/en.json",
		"--filename=admin/main.json",
		"--lang-iso=en",
		"--data=ZGF0YQ==",
		"--format=json",
		"--tags=mobile,web",
		"--convert-placeholders",
		"--apply-tm",
		"--hidden-from-contributors",
		"--custom-translation-status-ids=1,2",
		"--filter-task-id=123",
	})
	if err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	if !flags.Poll {
		t.Fatal("expected Poll to be true")
	}
	if flags.ContextTimeout != 45*time.Second {
		t.Fatalf("unexpected ContextTimeout: got %v", flags.ContextTimeout)
	}
	if flags.SrcPath != "./locales/en.json" {
		t.Fatalf("unexpected SrcPath: got %q", flags.SrcPath)
	}
	if flags.Filename != "admin/main.json" {
		t.Fatalf("unexpected Filename: got %q", flags.Filename)
	}
	if flags.LangISO != "en" {
		t.Fatalf("unexpected LangISO: got %q", flags.LangISO)
	}
	if flags.Data != "ZGF0YQ==" {
		t.Fatalf("unexpected Data: got %q", flags.Data)
	}
	if flags.Format != "json" {
		t.Fatalf("unexpected Format: got %q", flags.Format)
	}
	if len(flags.Tags) != 2 || flags.Tags[0] != "mobile" || flags.Tags[1] != "web" {
		t.Fatalf("unexpected Tags: got %#v", flags.Tags)
	}
	if !flags.ConvertPlaceholders {
		t.Fatal("expected ConvertPlaceholders to be true")
	}
	if !flags.ApplyTM {
		t.Fatal("expected ApplyTM to be true")
	}
	if !flags.HiddenFromContributors {
		t.Fatal("expected HiddenFromContributors to be true")
	}
	if len(flags.CustomTranslationStatusIDs) != 2 ||
		flags.CustomTranslationStatusIDs[0] != "1" ||
		flags.CustomTranslationStatusIDs[1] != "2" {
		t.Fatalf("unexpected CustomTranslationStatusIDs: got %#v", flags.CustomTranslationStatusIDs)
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

	if flags.ContextTimeout != 150*time.Second {
		t.Fatalf("unexpected ContextTimeout: got %v, want %v", flags.ContextTimeout, 150*time.Second)
	}
	if flags.Poll {
		t.Fatal("expected Poll to remain false")
	}
	if flags.Filename != "" {
		t.Fatalf("expected empty Filename, got %q", flags.Filename)
	}
	if flags.LangISO != "" {
		t.Fatalf("expected empty LangISO, got %q", flags.LangISO)
	}
	if flags.Data != "" {
		t.Fatalf("expected empty Data, got %q", flags.Data)
	}
	if flags.Format != "" {
		t.Fatalf("expected empty Format, got %q", flags.Format)
	}
	if flags.Tags != nil {
		t.Fatalf("expected nil Tags, got %#v", flags.Tags)
	}
	if flags.CustomTranslationStatusIDs != nil {
		t.Fatalf("expected nil CustomTranslationStatusIDs, got %#v", flags.CustomTranslationStatusIDs)
	}
	if flags.FilterTaskID != 0 {
		t.Fatalf("expected zero FilterTaskID, got %d", flags.FilterTaskID)
	}
}

func TestBindFlags_PreservesPreconfiguredValues(t *testing.T) {
	t.Parallel()

	cmd := &cobra.Command{Use: "test"}
	flags := newFlags()
	flags.Poll = true
	flags.ContextTimeout = 90 * time.Second
	flags.SrcPath = "./seed.json"
	flags.Filename = "seed.json"
	flags.LangISO = "fr"
	flags.Data = "seed-data"
	flags.Format = "xml"
	flags.Tags = []string{"existing"}
	flags.FilterTaskID = 55

	bindFlags(cmd, flags)

	if err := cmd.Flags().Parse(nil); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	if !flags.Poll {
		t.Fatal("expected Poll to preserve preconfigured value")
	}
	if flags.ContextTimeout != 90*time.Second {
		t.Fatalf("unexpected ContextTimeout: got %v", flags.ContextTimeout)
	}
	if flags.SrcPath != "./seed.json" {
		t.Fatalf("unexpected SrcPath: got %q", flags.SrcPath)
	}
	if flags.Filename != "seed.json" {
		t.Fatalf("unexpected Filename: got %q", flags.Filename)
	}
	if flags.LangISO != "fr" {
		t.Fatalf("unexpected LangISO: got %q", flags.LangISO)
	}
	if flags.Data != "seed-data" {
		t.Fatalf("unexpected Data: got %q", flags.Data)
	}
	if flags.Format != "xml" {
		t.Fatalf("unexpected Format: got %q", flags.Format)
	}
	if len(flags.Tags) != 1 || flags.Tags[0] != "existing" {
		t.Fatalf("unexpected Tags: got %#v", flags.Tags)
	}
	if flags.FilterTaskID != 55 {
		t.Fatalf("unexpected FilterTaskID: got %d", flags.FilterTaskID)
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
	flags.HiddenFromContributors = true

	bindFlags(cmd, flags)

	if err := cmd.Flags().Parse([]string{"--hidden-from-contributors=false"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	if flags.HiddenFromContributors {
		t.Fatal("expected HiddenFromContributors to be false")
	}

	if !cmd.Flags().Changed("hidden-from-contributors") {
		t.Fatal("expected hidden-from-contributors to be marked as changed")
	}
}
