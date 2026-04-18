package params

import (
	"errors"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type testFlags struct {
	Values []string
}

type testCfg struct {
	Values []string
}

type testReq struct {
	Values []string
}

func TestBindFlags_BindsAllSpecsInOrderAndDisablesSorting(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	flags := &testFlags{}

	var called []string

	specs := []ParamSpec[testFlags, testCfg, *testReq]{
		{
			FlagName: "first",
			BindFlag: func(cmd *cobra.Command, flags *testFlags) {
				called = append(called, "first")
				flags.Values = append(flags.Values, "first")
			},
			ApplyDefault:  func(*cobra.Command, *testFlags, *testCfg) {},
			LoadFromViper: func(*viper.Viper, *testCfg) {},
		},
		{
			FlagName: "second",
			BindFlag: func(cmd *cobra.Command, flags *testFlags) {
				called = append(called, "second")
				flags.Values = append(flags.Values, "second")
			},
			ApplyDefault:  func(*cobra.Command, *testFlags, *testCfg) {},
			LoadFromViper: func(*viper.Viper, *testCfg) {},
		},
	}

	BindFlags(cmd, flags, specs)

	if cmd.Flags().SortFlags {
		t.Fatal("expected SortFlags to be false")
	}

	wantCalled := []string{"first", "second"}
	if !reflect.DeepEqual(called, wantCalled) {
		t.Fatalf("unexpected bind call order: got %v, want %v", called, wantCalled)
	}

	if !reflect.DeepEqual(flags.Values, wantCalled) {
		t.Fatalf("unexpected flags values: got %v, want %v", flags.Values, wantCalled)
	}
}

func TestApplyDefaults_ReturnsImmediatelyWhenCfgIsNil(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	flags := &testFlags{}

	called := false

	specs := []ParamSpec[testFlags, testCfg, *testReq]{
		{
			BindFlag: func(*cobra.Command, *testFlags) {},
			ApplyDefault: func(*cobra.Command, *testFlags, *testCfg) {
				called = true
			},
			LoadFromViper: func(*viper.Viper, *testCfg) {},
		},
	}

	ApplyDefaults(cmd, flags, nil, specs)

	if called {
		t.Fatal("expected ApplyDefault not to be called when cfg is nil")
	}
}

func TestApplyDefaults_AppliesAllSpecsInOrder(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	flags := &testFlags{}
	cfg := &testCfg{}

	var called []string

	specs := []ParamSpec[testFlags, testCfg, *testReq]{
		{
			BindFlag: func(*cobra.Command, *testFlags) {},
			ApplyDefault: func(cmd *cobra.Command, flags *testFlags, cfg *testCfg) {
				called = append(called, "first")
				flags.Values = append(flags.Values, "first")
				cfg.Values = append(cfg.Values, "first")
			},
			LoadFromViper: func(*viper.Viper, *testCfg) {},
		},
		{
			BindFlag: func(*cobra.Command, *testFlags) {},
			ApplyDefault: func(cmd *cobra.Command, flags *testFlags, cfg *testCfg) {
				called = append(called, "second")
				flags.Values = append(flags.Values, "second")
				cfg.Values = append(cfg.Values, "second")
			},
			LoadFromViper: func(*viper.Viper, *testCfg) {},
		},
	}

	ApplyDefaults(cmd, flags, cfg, specs)

	want := []string{"first", "second"}

	if !reflect.DeepEqual(called, want) {
		t.Fatalf("unexpected ApplyDefault call order: got %v, want %v", called, want)
	}
	if !reflect.DeepEqual(flags.Values, want) {
		t.Fatalf("unexpected flags values: got %v, want %v", flags.Values, want)
	}
	if !reflect.DeepEqual(cfg.Values, want) {
		t.Fatalf("unexpected cfg values: got %v, want %v", cfg.Values, want)
	}
}

func TestLoadFromViper_LoadsAllSpecsInOrder(t *testing.T) {
	v := viper.New()
	cfg := &testCfg{}

	var called []string

	specs := []ParamSpec[testFlags, testCfg, *testReq]{
		{
			BindFlag:     func(*cobra.Command, *testFlags) {},
			ApplyDefault: func(*cobra.Command, *testFlags, *testCfg) {},
			LoadFromViper: func(v *viper.Viper, cfg *testCfg) {
				called = append(called, "first")
				cfg.Values = append(cfg.Values, "first")
			},
		},
		{
			BindFlag:     func(*cobra.Command, *testFlags) {},
			ApplyDefault: func(*cobra.Command, *testFlags, *testCfg) {},
			LoadFromViper: func(v *viper.Viper, cfg *testCfg) {
				called = append(called, "second")
				cfg.Values = append(cfg.Values, "second")
			},
		},
	}

	LoadFromViper(v, cfg, specs)

	want := []string{"first", "second"}

	if !reflect.DeepEqual(called, want) {
		t.Fatalf("unexpected LoadFromViper call order: got %v, want %v", called, want)
	}
	if !reflect.DeepEqual(cfg.Values, want) {
		t.Fatalf("unexpected cfg values: got %v, want %v", cfg.Values, want)
	}
}

func TestConfigKeys_ReturnsKeysInOrder(t *testing.T) {
	specs := []ParamSpec[testFlags, testCfg, *testReq]{
		{ConfigKey: "upload.filename"},
		{ConfigKey: "upload.lang_iso"},
		{ConfigKey: "upload.replace_modified"},
	}

	got := ConfigKeys(specs)
	want := []string{"upload.filename", "upload.lang_iso", "upload.replace_modified"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected config keys: got %v, want %v", got, want)
	}
}

func TestApplyToRequest_UsesEmptyDefaultsWhenNil(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	flags := &testFlags{}
	req := &testReq{}

	called := false

	specs := []ParamSpec[testFlags, testCfg, *testReq]{
		{
			ApplyToRequest: func(cmd *cobra.Command, flags *testFlags, defaults *testCfg, req *testReq) error {
				called = true
				if defaults == nil {
					t.Fatal("expected defaults to be non-nil")
				}
				defaults.Values = append(defaults.Values, "set-inside")
				req.Values = append(req.Values, "applied")
				return nil
			},
		},
	}

	err := ApplyToRequest(cmd, flags, nil, req, specs)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !called {
		t.Fatal("expected ApplyToRequest callback to be called")
	}

	wantReq := []string{"applied"}
	if !reflect.DeepEqual(req.Values, wantReq) {
		t.Fatalf("unexpected req values: got %v, want %v", req.Values, wantReq)
	}
}

func TestApplyToRequest_SkipsNilCallbacks(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	flags := &testFlags{}
	defaults := &testCfg{}
	req := &testReq{}

	called := false

	specs := []ParamSpec[testFlags, testCfg, *testReq]{
		{},
		{
			ApplyToRequest: func(cmd *cobra.Command, flags *testFlags, defaults *testCfg, req *testReq) error {
				called = true
				req.Values = append(req.Values, "called")
				return nil
			},
		},
	}

	err := ApplyToRequest(cmd, flags, defaults, req, specs)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !called {
		t.Fatal("expected non-nil ApplyToRequest callback to be called")
	}

	wantReq := []string{"called"}
	if !reflect.DeepEqual(req.Values, wantReq) {
		t.Fatalf("unexpected req values: got %v, want %v", req.Values, wantReq)
	}
}

func TestApplyToRequest_StopsOnFirstError(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	flags := &testFlags{}
	defaults := &testCfg{}
	req := &testReq{}

	expectedErr := errors.New("boom")
	var called []string

	specs := []ParamSpec[testFlags, testCfg, *testReq]{
		{
			ApplyToRequest: func(cmd *cobra.Command, flags *testFlags, defaults *testCfg, req *testReq) error {
				called = append(called, "first")
				return expectedErr
			},
		},
		{
			ApplyToRequest: func(cmd *cobra.Command, flags *testFlags, defaults *testCfg, req *testReq) error {
				called = append(called, "second")
				return nil
			},
		},
	}

	err := ApplyToRequest(cmd, flags, defaults, req, specs)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	wantCalled := []string{"first"}
	if !reflect.DeepEqual(called, wantCalled) {
		t.Fatalf("unexpected ApplyToRequest call order: got %v, want %v", called, wantCalled)
	}
}
