package download

import (
	"reflect"
	"strings"
	"testing"

	lokexdownload "github.com/bodrovis/lokex/v2/client/download"
	"github.com/spf13/cobra"
)

func TestReqString_SetsNonEmptyString(t *testing.T) {
	req := lokexdownload.DownloadParams{}
	flags := &Flags{BundleStructure: "bundle/%LANG_ISO%.json"}

	fn := reqString("bundle_structure", func(f *Flags) string {
		return f.BundleStructure
	})

	err := fn(&cobra.Command{}, flags, nil, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got := req["bundle_structure"]; got != "bundle/%LANG_ISO%.json" {
		t.Fatalf("unexpected value: got %#v", got)
	}
}

func TestReqString_DoesNotSetEmptyString(t *testing.T) {
	req := lokexdownload.DownloadParams{}
	flags := &Flags{}

	fn := reqString("bundle_structure", func(f *Flags) string {
		return f.BundleStructure
	})

	err := fn(&cobra.Command{}, flags, nil, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, ok := req["bundle_structure"]; ok {
		t.Fatal("expected empty string not to be set")
	}
}

func TestReqDirectString_SetsValueEvenWhenEmpty(t *testing.T) {
	req := lokexdownload.DownloadParams{}
	flags := &Flags{}

	fn := reqDirectString("format", func(f *Flags) string {
		return f.Format
	})

	err := fn(&cobra.Command{}, flags, nil, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, ok := req["format"]
	if !ok {
		t.Fatal("expected key to be set")
	}
	if got != "" {
		t.Fatalf("unexpected value: got %#v, want empty string", got)
	}
}

func TestReqStringSlice_SetsNonEmptySlice(t *testing.T) {
	req := lokexdownload.DownloadParams{}
	flags := &Flags{FilterLangs: []string{"en", "fr"}}

	fn := reqStringSlice("filter_langs", func(f *Flags) []string {
		return f.FilterLangs
	})

	err := fn(&cobra.Command{}, flags, nil, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, ok := req["filter_langs"]
	if !ok {
		t.Fatal("expected key to be set")
	}

	want := []string{"en", "fr"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected value: got %#v, want %#v", got, want)
	}
}

func TestReqStringSlice_DoesNotSetEmptySlice(t *testing.T) {
	req := lokexdownload.DownloadParams{}
	flags := &Flags{}

	fn := reqStringSlice("filter_langs", func(f *Flags) []string {
		return f.FilterLangs
	})

	err := fn(&cobra.Command{}, flags, nil, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, ok := req["filter_langs"]; ok {
		t.Fatal("expected empty slice not to be set")
	}
}

func TestReqBoolWithDefault_SetsFlagValueWhenFlagChanged(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	flags := &Flags{}
	cmd.Flags().BoolVar(&flags.Async, "async", flags.Async, "")

	if err := cmd.Flags().Set("async", "true"); err != nil {
		t.Fatalf("failed to set flag: %v", err)
	}

	req := lokexdownload.DownloadParams{}

	fn := reqBoolWithDefault(
		"async",
		"async",
		func(f *Flags) bool { return f.Async },
		func(c *DownloadConfig) *bool { return c.Async },
	)

	err := fn(cmd, flags, &DownloadConfig{Async: new(false)}, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got := req["async"]; got != true {
		t.Fatalf("unexpected value: got %#v, want true", got)
	}
}

func TestReqBoolWithDefault_SetsDefaultWhenFlagNotChanged(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	flags := &Flags{}
	cmd.Flags().BoolVar(&flags.Async, "async", flags.Async, "")

	req := lokexdownload.DownloadParams{}

	fn := reqBoolWithDefault(
		"async",
		"async",
		func(f *Flags) bool { return f.Async },
		func(c *DownloadConfig) *bool { return c.Async },
	)

	err := fn(cmd, flags, &DownloadConfig{Async: new(true)}, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got := req["async"]; got != true {
		t.Fatalf("unexpected value: got %#v, want true", got)
	}
}

func TestReqBoolWithDefault_DoesNotSetWhenFlagNotChangedAndNoDefault(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	flags := &Flags{}
	cmd.Flags().BoolVar(&flags.Async, "async", flags.Async, "")

	req := lokexdownload.DownloadParams{}

	fn := reqBoolWithDefault(
		"async",
		"async",
		func(f *Flags) bool { return f.Async },
		func(c *DownloadConfig) *bool { return c.Async },
	)

	err := fn(cmd, flags, &DownloadConfig{}, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, ok := req["async"]; ok {
		t.Fatal("expected key not to be set")
	}
}

func TestReqInt64WithDefault_SetsFlagValueWhenFlagChanged(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	flags := &Flags{}
	cmd.Flags().Int64Var(&flags.FilterTaskID, "filter-task-id", flags.FilterTaskID, "")

	if err := cmd.Flags().Set("filter-task-id", "123"); err != nil {
		t.Fatalf("failed to set flag: %v", err)
	}

	req := lokexdownload.DownloadParams{}

	fn := reqInt64WithDefault(
		"filter-task-id",
		"filter_task_id",
		func(f *Flags) int64 { return f.FilterTaskID },
		func(c *DownloadConfig) *int64 { return c.FilterTaskID },
	)

	err := fn(cmd, flags, &DownloadConfig{FilterTaskID: new(int64(999))}, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got := req["filter_task_id"]; got != int64(123) {
		t.Fatalf("unexpected value: got %#v, want 123", got)
	}
}

func TestReqInt64WithDefault_SetsDefaultWhenFlagNotChanged(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	flags := &Flags{}
	cmd.Flags().Int64Var(&flags.FilterTaskID, "filter-task-id", flags.FilterTaskID, "")

	req := lokexdownload.DownloadParams{}

	fn := reqInt64WithDefault(
		"filter-task-id",
		"filter_task_id",
		func(f *Flags) int64 { return f.FilterTaskID },
		func(c *DownloadConfig) *int64 { return c.FilterTaskID },
	)

	err := fn(cmd, flags, &DownloadConfig{FilterTaskID: new(int64(456))}, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got := req["filter_task_id"]; got != int64(456) {
		t.Fatalf("unexpected value: got %#v, want 456", got)
	}
}

func TestReqInt64WithDefault_DoesNotSetWhenFlagNotChangedAndNoDefault(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	flags := &Flags{}
	cmd.Flags().Int64Var(&flags.FilterTaskID, "filter-task-id", flags.FilterTaskID, "")

	req := lokexdownload.DownloadParams{}

	fn := reqInt64WithDefault(
		"filter-task-id",
		"filter_task_id",
		func(f *Flags) int64 { return f.FilterTaskID },
		func(c *DownloadConfig) *int64 { return c.FilterTaskID },
	)

	err := fn(cmd, flags, &DownloadConfig{}, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, ok := req["filter_task_id"]; ok {
		t.Fatal("expected key not to be set")
	}
}

func TestReqLanguageMapping_IgnoresEmptyValue(t *testing.T) {
	req := lokexdownload.DownloadParams{}
	flags := &Flags{}

	fn := reqLanguageMapping("language-mapping", func(f *Flags) string {
		return f.LanguageMappingJSON
	})

	err := fn(&cobra.Command{}, flags, nil, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, ok := req["language_mapping"]; ok {
		t.Fatal("expected key not to be set")
	}
}

func TestReqLanguageMapping_ParsesAndSetsValue(t *testing.T) {
	req := lokexdownload.DownloadParams{}
	flags := &Flags{
		LanguageMappingJSON: `[{"original_language_iso":"en","custom_language_iso":"en-US"}]`,
	}

	fn := reqLanguageMapping("language-mapping", func(f *Flags) string {
		return f.LanguageMappingJSON
	})

	err := fn(&cobra.Command{}, flags, nil, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, ok := req["language_mapping"]
	if !ok {
		t.Fatal("expected key to be set")
	}

	want := []map[string]any{
		{
			"original_language_iso": "en",
			"custom_language_iso":   "en-US",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected value: got %#v, want %#v", got, want)
	}
}

func TestReqLanguageMapping_ReturnsWrappedErrorForInvalidJSON(t *testing.T) {
	req := lokexdownload.DownloadParams{}
	flags := &Flags{
		LanguageMappingJSON: `{"broken":true}`,
	}

	fn := reqLanguageMapping("language-mapping", func(f *Flags) string {
		return f.LanguageMappingJSON
	})

	err := fn(&cobra.Command{}, flags, nil, req)
	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "parse --language-mapping") {
		t.Fatalf("expected wrapped flag name in error, got %v", err)
	}
}

func TestParseLanguageMapping_ValidJSON(t *testing.T) {
	got, err := parseLanguageMapping(`[{"original_language_iso":"en","custom_language_iso":"en-US"}]`)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := []map[string]any{
		{
			"original_language_iso": "en",
			"custom_language_iso":   "en-US",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected value: got %#v, want %#v", got, want)
	}
}

func TestParseLanguageMapping_InvalidJSON(t *testing.T) {
	_, err := parseLanguageMapping(`{"not":"an array"}`)
	if err == nil {
		t.Fatal("expected error")
	}
}
