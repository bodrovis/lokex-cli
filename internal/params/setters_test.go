package params

import (
	"reflect"
	"testing"
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
		if got != "json" {
			t.Fatalf("expected trimmed value %q, got %v", "json", got)
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

	t.Run("trims out empty and whitespace-only values", func(t *testing.T) {
		t.Parallel()

		params := testParams{}
		value := []string{"en", "", "   ", "fr"}

		SetStringSlice(params, "filter_langs", value)

		got, ok := params["filter_langs"]
		if !ok {
			t.Fatal("expected key to be set")
		}

		want := []string{"en", "fr"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("expected %v, got %v", want, got)
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

	t.Run("does not set slice with only empty values", func(t *testing.T) {
		t.Parallel()

		params := testParams{}
		SetStringSlice(params, "filter_langs", []string{"", "   "})

		if _, ok := params["filter_langs"]; ok {
			t.Fatal("expected key not to be set")
		}
	})
}
