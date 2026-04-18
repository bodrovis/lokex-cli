package upload

import (
	"reflect"
	"testing"

	lokexupload "github.com/bodrovis/lokex/v2/client/upload"
	"github.com/spf13/cobra"
)

func boolPtr(v bool) *bool    { return &v }
func int64Ptr(v int64) *int64 { return &v }

func TestReqString_SetsNonEmptyString(t *testing.T) {
	req := lokexupload.UploadParams{}
	flags := &Flags{Data: "YmFzZTY0"}

	fn := reqString("data", func(f *Flags) string {
		return f.Data
	})

	err := fn(&cobra.Command{}, flags, nil, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got := req["data"]; got != "YmFzZTY0" {
		t.Fatalf("unexpected value: got %#v", got)
	}
}

func TestReqString_DoesNotSetEmptyString(t *testing.T) {
	req := lokexupload.UploadParams{}
	flags := &Flags{}

	fn := reqString("data", func(f *Flags) string {
		return f.Data
	})

	err := fn(&cobra.Command{}, flags, nil, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, ok := req["data"]; ok {
		t.Fatal("expected empty string not to be set")
	}
}

func TestReqTrimmedString_SetsTrimmedValue(t *testing.T) {
	req := lokexupload.UploadParams{}
	flags := &Flags{LangISO: "  en  "}

	fn := reqTrimmedString("lang_iso", func(f *Flags) string {
		return f.LangISO
	})

	err := fn(&cobra.Command{}, flags, nil, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got := req["lang_iso"]; got != "en" {
		t.Fatalf("unexpected value: got %#v, want %q", got, "en")
	}
}

func TestReqTrimmedString_SetsEmptyStringAfterTrim(t *testing.T) {
	req := lokexupload.UploadParams{}
	flags := &Flags{LangISO: "   "}

	fn := reqTrimmedString("lang_iso", func(f *Flags) string {
		return f.LangISO
	})

	err := fn(&cobra.Command{}, flags, nil, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, ok := req["lang_iso"]
	if !ok {
		t.Fatal("expected key to be set")
	}
	if got != "" {
		t.Fatalf("unexpected value: got %#v, want empty string", got)
	}
}

func TestReqNormalizedFilename_NormalizesSlashesAndTrims(t *testing.T) {
	req := lokexupload.UploadParams{}
	flags := &Flags{Filename: `  dir\subdir\messages.json  `}

	fn := reqNormalizedFilename("filename", func(f *Flags) string {
		return f.Filename
	})

	err := fn(&cobra.Command{}, flags, nil, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got := req["filename"]; got != "dir/subdir/messages.json" {
		t.Fatalf("unexpected value: got %#v, want %q", got, "dir/subdir/messages.json")
	}
}

func TestReqNormalizedFilename_SetsEmptyStringWhenBlank(t *testing.T) {
	req := lokexupload.UploadParams{}
	flags := &Flags{Filename: "   "}

	fn := reqNormalizedFilename("filename", func(f *Flags) string {
		return f.Filename
	})

	err := fn(&cobra.Command{}, flags, nil, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, ok := req["filename"]
	if !ok {
		t.Fatal("expected key to be set")
	}
	if got != "" {
		t.Fatalf("unexpected value: got %#v, want empty string", got)
	}
}

func TestReqStringSlice_SetsNonEmptySlice(t *testing.T) {
	req := lokexupload.UploadParams{}
	flags := &Flags{Tags: []string{"mobile", "backend"}}

	fn := reqStringSlice("tags", func(f *Flags) []string {
		return f.Tags
	})

	err := fn(&cobra.Command{}, flags, nil, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, ok := req["tags"]
	if !ok {
		t.Fatal("expected key to be set")
	}

	want := []string{"mobile", "backend"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected value: got %#v, want %#v", got, want)
	}
}

func TestReqStringSlice_DoesNotSetEmptySlice(t *testing.T) {
	req := lokexupload.UploadParams{}
	flags := &Flags{}

	fn := reqStringSlice("tags", func(f *Flags) []string {
		return f.Tags
	})

	err := fn(&cobra.Command{}, flags, nil, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, ok := req["tags"]; ok {
		t.Fatal("expected empty slice not to be set")
	}
}

func TestReqBoolWithDefault_SetsFlagValueWhenFlagChanged(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	flags := &Flags{}
	cmd.Flags().BoolVar(&flags.ConvertPlaceholders, "convert-placeholders", flags.ConvertPlaceholders, "")

	if err := cmd.Flags().Set("convert-placeholders", "true"); err != nil {
		t.Fatalf("failed to set flag: %v", err)
	}

	req := lokexupload.UploadParams{}

	fn := reqBoolWithDefault(
		"convert-placeholders",
		"convert_placeholders",
		func(f *Flags) bool { return f.ConvertPlaceholders },
		func(c *UploadConfig) *bool { return c.ConvertPlaceholders },
	)

	err := fn(cmd, flags, &UploadConfig{ConvertPlaceholders: boolPtr(false)}, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got := req["convert_placeholders"]; got != true {
		t.Fatalf("unexpected value: got %#v, want true", got)
	}
}

func TestReqBoolWithDefault_SetsDefaultWhenFlagNotChanged(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	flags := &Flags{}
	cmd.Flags().BoolVar(&flags.ConvertPlaceholders, "convert-placeholders", flags.ConvertPlaceholders, "")

	req := lokexupload.UploadParams{}

	fn := reqBoolWithDefault(
		"convert-placeholders",
		"convert_placeholders",
		func(f *Flags) bool { return f.ConvertPlaceholders },
		func(c *UploadConfig) *bool { return c.ConvertPlaceholders },
	)

	err := fn(cmd, flags, &UploadConfig{ConvertPlaceholders: boolPtr(true)}, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got := req["convert_placeholders"]; got != true {
		t.Fatalf("unexpected value: got %#v, want true", got)
	}
}

func TestReqBoolWithDefault_DoesNotSetWhenFlagNotChangedAndNoDefault(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	flags := &Flags{}
	cmd.Flags().BoolVar(&flags.ConvertPlaceholders, "convert-placeholders", flags.ConvertPlaceholders, "")

	req := lokexupload.UploadParams{}

	fn := reqBoolWithDefault(
		"convert-placeholders",
		"convert_placeholders",
		func(f *Flags) bool { return f.ConvertPlaceholders },
		func(c *UploadConfig) *bool { return c.ConvertPlaceholders },
	)

	err := fn(cmd, flags, &UploadConfig{}, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, ok := req["convert_placeholders"]; ok {
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

	req := lokexupload.UploadParams{}

	fn := reqInt64WithDefault(
		"filter-task-id",
		"filter_task_id",
		func(f *Flags) int64 { return f.FilterTaskID },
		func(c *UploadConfig) *int64 { return c.FilterTaskID },
	)

	err := fn(cmd, flags, &UploadConfig{FilterTaskID: int64Ptr(999)}, req)
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

	req := lokexupload.UploadParams{}

	fn := reqInt64WithDefault(
		"filter-task-id",
		"filter_task_id",
		func(f *Flags) int64 { return f.FilterTaskID },
		func(c *UploadConfig) *int64 { return c.FilterTaskID },
	)

	err := fn(cmd, flags, &UploadConfig{FilterTaskID: int64Ptr(456)}, req)
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

	req := lokexupload.UploadParams{}

	fn := reqInt64WithDefault(
		"filter-task-id",
		"filter_task_id",
		func(f *Flags) int64 { return f.FilterTaskID },
		func(c *UploadConfig) *int64 { return c.FilterTaskID },
	)

	err := fn(cmd, flags, &UploadConfig{}, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, ok := req["filter_task_id"]; ok {
		t.Fatal("expected key not to be set")
	}
}
