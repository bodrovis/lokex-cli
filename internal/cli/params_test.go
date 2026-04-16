package cli

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

type testParams map[string]any

func TestSetString(t *testing.T) {
	t.Parallel()

	t.Run("sets non-empty trimmed value", func(t *testing.T) {
		t.Parallel()

		params := testParams{}
		SetString(params, "format", " json ")

		got, ok := params["format"]
		if !ok {
			t.Fatal("expected key to be set")
		}
		if got != " json " {
			t.Fatalf("expected original value to be preserved, got %v", got)
		}
	})

	t.Run("does not set empty value", func(t *testing.T) {
		t.Parallel()

		params := testParams{}
		SetString(params, "format", "")
		if _, ok := params["format"]; ok {
			t.Fatal("expected key not to be set")
		}
	})

	t.Run("does not set whitespace-only value", func(t *testing.T) {
		t.Parallel()

		params := testParams{}
		SetString(params, "format", "   ")
		if _, ok := params["format"]; ok {
			t.Fatal("expected key not to be set")
		}
	})
}

func TestSetStringSlice(t *testing.T) {
	t.Parallel()

	t.Run("sets non-empty slice", func(t *testing.T) {
		t.Parallel()

		params := testParams{}
		value := []string{"en", "fr"}

		SetStringSlice(params, "filter_langs", value)

		got, ok := params["filter_langs"]
		if !ok {
			t.Fatal("expected key to be set")
		}
		if !reflect.DeepEqual(got, value) {
			t.Fatalf("expected %v, got %v", value, got)
		}
	})

	t.Run("does not set nil slice", func(t *testing.T) {
		t.Parallel()

		params := testParams{}
		SetStringSlice(params, "filter_langs", nil)

		if _, ok := params["filter_langs"]; ok {
			t.Fatal("expected key not to be set")
		}
	})

	t.Run("does not set empty slice", func(t *testing.T) {
		t.Parallel()

		params := testParams{}
		SetStringSlice(params, "filter_langs", []string{})

		if _, ok := params["filter_langs"]; ok {
			t.Fatal("expected key not to be set")
		}
	})
}

func TestSetChangedBool(t *testing.T) {
	t.Parallel()

	t.Run("sets bool when flag changed to true", func(t *testing.T) {
		t.Parallel()

		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().Bool("compact", false, "")
		if err := cmd.Flags().Parse([]string{"--compact"}); err != nil {
			t.Fatalf("parse flags: %v", err)
		}

		params := testParams{}
		SetChangedBool(cmd, params, "compact", "compact", true)

		got, ok := params["compact"]
		if !ok {
			t.Fatal("expected key to be set")
		}
		if got != true {
			t.Fatalf("expected true, got %v", got)
		}
	})

	t.Run("sets bool when flag changed to false explicitly", func(t *testing.T) {
		t.Parallel()

		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().Bool("compact", true, "")
		if err := cmd.Flags().Parse([]string{"--compact=false"}); err != nil {
			t.Fatalf("parse flags: %v", err)
		}

		params := testParams{}
		SetChangedBool(cmd, params, "compact", "compact", false)

		got, ok := params["compact"]
		if !ok {
			t.Fatal("expected key to be set")
		}
		if got != false {
			t.Fatalf("expected false, got %v", got)
		}
	})

	t.Run("does not set bool when flag not changed", func(t *testing.T) {
		t.Parallel()

		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().Bool("compact", false, "")
		if err := cmd.Flags().Parse(nil); err != nil {
			t.Fatalf("parse flags: %v", err)
		}

		params := testParams{}
		SetChangedBool(cmd, params, "compact", "compact", false)

		if _, ok := params["compact"]; ok {
			t.Fatal("expected key not to be set")
		}
	})
}
