package viper_helpers

import (
	"reflect"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestApplyConfigValue_String_DirectType(t *testing.T) {
	v := viper.New()
	v.Set("key", "value")

	var got string
	called := false

	ApplyConfigValue[string](v, "key", func(val string) {
		called = true
		got = val
	})

	if !called {
		t.Fatal("expected setter to be called")
	}
	if got != "value" {
		t.Fatalf("unexpected value: got %q, want %q", got, "value")
	}
}

func TestApplyConfigValue_Bool_DirectType(t *testing.T) {
	v := viper.New()
	v.Set("key", true)

	var got bool
	called := false

	ApplyConfigValue[bool](v, "key", func(val bool) {
		called = true
		got = val
	})

	if !called {
		t.Fatal("expected setter to be called")
	}
	if got != true {
		t.Fatalf("unexpected value: got %v, want true", got)
	}
}

func TestApplyConfigValue_Bool_FallbackGetter(t *testing.T) {
	v := viper.New()
	v.Set("key", "true")

	var got bool
	called := false

	ApplyConfigValue[bool](v, "key", func(val bool) {
		called = true
		got = val
	})

	if !called {
		t.Fatal("expected setter to be called")
	}
	if got != true {
		t.Fatalf("unexpected value: got %v, want true", got)
	}
}

func TestApplyConfigValue_Int64_DirectType(t *testing.T) {
	v := viper.New()
	v.Set("key", int64(42))

	var got int64
	called := false

	ApplyConfigValue[int64](v, "key", func(val int64) {
		called = true
		got = val
	})

	if !called {
		t.Fatal("expected setter to be called")
	}
	if got != 42 {
		t.Fatalf("unexpected value: got %d, want 42", got)
	}
}

func TestApplyConfigValue_Int64_FallbackGetter(t *testing.T) {
	v := viper.New()
	v.Set("key", "42")

	var got int64
	called := false

	ApplyConfigValue[int64](v, "key", func(val int64) {
		called = true
		got = val
	})

	if !called {
		t.Fatal("expected setter to be called")
	}
	if got != 42 {
		t.Fatalf("unexpected value: got %d, want 42", got)
	}
}

func TestApplyConfigValue_Int_DirectType(t *testing.T) {
	v := viper.New()
	v.Set("key", int(42))

	var got int
	called := false

	ApplyConfigValue[int](v, "key", func(val int) {
		called = true
		got = val
	})

	if !called {
		t.Fatal("expected setter to be called")
	}
	if got != 42 {
		t.Fatalf("unexpected value: got %d, want 42", got)
	}
}

func TestApplyConfigValue_Duration_DirectType(t *testing.T) {
	v := viper.New()
	v.Set("key", 5*time.Second)

	var got time.Duration
	called := false

	ApplyConfigValue[time.Duration](v, "key", func(val time.Duration) {
		called = true
		got = val
	})

	if !called {
		t.Fatal("expected setter to be called")
	}
	if got != 5*time.Second {
		t.Fatalf("unexpected value: got %v, want %v", got, 5*time.Second)
	}
}

func TestApplyConfigValue_Duration_FallbackGetter(t *testing.T) {
	v := viper.New()
	v.Set("key", "5s")

	var got time.Duration
	called := false

	ApplyConfigValue[time.Duration](v, "key", func(val time.Duration) {
		called = true
		got = val
	})

	if !called {
		t.Fatal("expected setter to be called")
	}
	if got != 5*time.Second {
		t.Fatalf("unexpected value: got %v, want %v", got, 5*time.Second)
	}
}

func TestApplyConfigValue_KeyNotSet_DoesNothing(t *testing.T) {
	v := viper.New()

	called := false

	ApplyConfigValue[string](v, "missing", func(val string) {
		called = true
	})

	if called {
		t.Fatal("expected setter not to be called")
	}
}

func TestApplyConfigValue_UnsupportedType_DoesNothing(t *testing.T) {
	v := viper.New()
	v.Set("key", "123.4")

	called := false

	ApplyConfigValue[float32](v, "key", func(val float32) {
		called = true
	})

	if called {
		t.Fatal("expected setter not to be called for unsupported type")
	}
}

func TestApplyConfigStringSlice_FromStringSlice(t *testing.T) {
	v := viper.New()
	v.Set("key", []string{" one ", "", "two", "   ", " three "})

	var got []string
	called := false

	ApplyConfigStringSlice(v, "key", func(val []string) {
		called = true
		got = val
	})

	want := []string{"one", "two", "three"}

	if !called {
		t.Fatal("expected setter to be called")
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected value: got %v, want %v", got, want)
	}
}

func TestApplyConfigStringSlice_FromAnySlice(t *testing.T) {
	v := viper.New()
	v.Set("key", []any{" one ", 2, true, "", "   ", "three"})

	var got []string
	called := false

	ApplyConfigStringSlice(v, "key", func(val []string) {
		called = true
		got = val
	})

	want := []string{"one", "2", "true", "three"}

	if !called {
		t.Fatal("expected setter to be called")
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected value: got %v, want %v", got, want)
	}
}

func TestApplyConfigStringSlice_FromCommaSeparatedString(t *testing.T) {
	v := viper.New()
	v.Set("key", " one, , two ,three ,, ")

	var got []string
	called := false

	ApplyConfigStringSlice(v, "key", func(val []string) {
		called = true
		got = val
	})

	want := []string{"one", "two", "three"}

	if !called {
		t.Fatal("expected setter to be called")
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected value: got %v, want %v", got, want)
	}
}

func TestApplyConfigStringSlice_KeyNotSet_DoesNothing(t *testing.T) {
	v := viper.New()

	called := false

	ApplyConfigStringSlice(v, "missing", func(val []string) {
		called = true
	})

	if called {
		t.Fatal("expected setter not to be called")
	}
}

func TestApplyConfigStringSlice_EmptyCleanedStringSlice_DoesNothing(t *testing.T) {
	v := viper.New()
	v.Set("key", []string{"", "   ", "\t"})

	called := false

	ApplyConfigStringSlice(v, "key", func(val []string) {
		called = true
	})

	if called {
		t.Fatal("expected setter not to be called")
	}
}

func TestApplyConfigStringSlice_EmptyCleanedAnySlice_DoesNothing(t *testing.T) {
	v := viper.New()
	v.Set("key", []any{"", "   "})

	called := false

	ApplyConfigStringSlice(v, "key", func(val []string) {
		called = true
	})

	if called {
		t.Fatal("expected setter not to be called")
	}
}

func TestApplyConfigStringSlice_EmptyCleanedString_DoesNothing(t *testing.T) {
	v := viper.New()
	v.Set("key", " , ,   , ")

	called := false

	ApplyConfigStringSlice(v, "key", func(val []string) {
		called = true
	})

	if called {
		t.Fatal("expected setter not to be called")
	}
}

func TestApplyConfigStringSlice_UnsupportedType_DoesNothing(t *testing.T) {
	v := viper.New()
	v.Set("key", 123)

	called := false

	ApplyConfigStringSlice(v, "key", func(val []string) {
		called = true
	})

	if called {
		t.Fatal("expected setter not to be called")
	}
}
