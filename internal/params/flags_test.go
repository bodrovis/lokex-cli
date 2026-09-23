package params

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func TestBindFlagSet_BindsAllKindsWithDefaults(t *testing.T) {
	specs := []Spec{
		{
			Name:    "string",
			Kind:    KindString,
			Usage:   "string usage",
			Default: "hello",
		},
		{
			Name:    "bool",
			Kind:    KindBool,
			Usage:   "bool usage",
			Default: true,
		},
		{
			Name:    "slice",
			Kind:    KindStringSlice,
			Usage:   "slice usage",
			Default: []string{"en", "de"},
		},
		{
			Name:    "int",
			Kind:    KindInt,
			Usage:   "int usage",
			Default: 42,
		},
		{
			Name:    "int64",
			Kind:    KindInt64,
			Usage:   "int64 usage",
			Default: int64(123),
		},
		{
			Name:    "duration",
			Kind:    KindDuration,
			Usage:   "duration usage",
			Default: 30 * time.Second,
		},
	}

	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	BindFlagSet(fs, specs)

	require.False(t, fs.SortFlags)

	gotString, err := fs.GetString("string")
	require.NoError(t, err)
	require.Equal(t, "hello", gotString)

	gotBool, err := fs.GetBool("bool")
	require.NoError(t, err)
	require.True(t, gotBool)

	gotSlice, err := fs.GetStringSlice("slice")
	require.NoError(t, err)
	require.Equal(t, []string{"en", "de"}, gotSlice)

	gotInt, err := fs.GetInt("int")
	require.NoError(t, err)
	require.Equal(t, 42, gotInt)

	gotInt64, err := fs.GetInt64("int64")
	require.NoError(t, err)
	require.Equal(t, int64(123), gotInt64)

	gotDuration, err := fs.GetDuration("duration")
	require.NoError(t, err)
	require.Equal(t, 30*time.Second, gotDuration)
}

func TestBindFlagSet_UsesZeroValueWhenDefaultIsNil(t *testing.T) {
	specs := []Spec{
		{
			Name: "string",
			Kind: KindString,
		},
		{
			Name: "bool",
			Kind: KindBool,
		},
		{
			Name: "slice",
			Kind: KindStringSlice,
		},
		{
			Name: "int",
			Kind: KindInt,
		},
		{
			Name: "int64",
			Kind: KindInt64,
		},
		{
			Name: "duration",
			Kind: KindDuration,
		},
	}

	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	BindFlagSet(fs, specs)

	gotString, err := fs.GetString("string")
	require.NoError(t, err)
	require.Equal(t, "", gotString)

	gotBool, err := fs.GetBool("bool")
	require.NoError(t, err)
	require.False(t, gotBool)

	gotSlice, err := fs.GetStringSlice("slice")
	require.NoError(t, err)
	require.Empty(t, gotSlice)

	gotInt, err := fs.GetInt("int")
	require.NoError(t, err)
	require.Equal(t, 0, gotInt)

	gotInt64, err := fs.GetInt64("int64")
	require.NoError(t, err)
	require.Equal(t, int64(0), gotInt64)

	gotDuration, err := fs.GetDuration("duration")
	require.NoError(t, err)
	require.Equal(t, time.Duration(0), gotDuration)
}

func TestBindFlagSet_PreservesExplicitZeroDefaults(t *testing.T) {
	specs := []Spec{
		{
			Name:    "bool",
			Kind:    KindBool,
			Default: false,
		},
		{
			Name:    "int",
			Kind:    KindInt,
			Default: 0,
		},
		{
			Name:    "int64",
			Kind:    KindInt64,
			Default: int64(0),
		},
		{
			Name:    "duration",
			Kind:    KindDuration,
			Default: time.Duration(0),
		},
	}

	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	BindFlagSet(fs, specs)

	gotBool, err := fs.GetBool("bool")
	require.NoError(t, err)
	require.False(t, gotBool)

	gotInt, err := fs.GetInt("int")
	require.NoError(t, err)
	require.Equal(t, 0, gotInt)

	gotInt64, err := fs.GetInt64("int64")
	require.NoError(t, err)
	require.Equal(t, int64(0), gotInt64)

	gotDuration, err := fs.GetDuration("duration")
	require.NoError(t, err)
	require.Equal(t, time.Duration(0), gotDuration)
}

func TestBindFlags_BindsToCommandFlagSet(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
	}

	BindFlags(cmd, []Spec{
		{
			Name:    "format",
			Kind:    KindString,
			Default: "json",
		},
	})

	got, err := cmd.Flags().GetString("format")

	require.NoError(t, err)
	require.Equal(t, "json", got)
}

func TestFlagValue(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
	}

	specs := []Spec{
		{
			Name: "string",
			Kind: KindString,
		},
		{
			Name: "bool",
			Kind: KindBool,
		},
		{
			Name: "slice",
			Kind: KindStringSlice,
		},
		{
			Name: "int",
			Kind: KindInt,
		},
		{
			Name: "int64",
			Kind: KindInt64,
		},
		{
			Name: "duration",
			Kind: KindDuration,
		},
	}

	BindFlags(cmd, specs)

	require.NoError(t, cmd.Flags().Set("string", "hello"))
	require.NoError(t, cmd.Flags().Set("bool", "true"))
	require.NoError(t, cmd.Flags().Set("slice", "en,de"))
	require.NoError(t, cmd.Flags().Set("int", "42"))
	require.NoError(t, cmd.Flags().Set("int64", "123"))
	require.NoError(t, cmd.Flags().Set("duration", "45s"))

	tests := []struct {
		spec Spec
		want any
	}{
		{
			spec: specs[0],
			want: "hello",
		},
		{
			spec: specs[1],
			want: true,
		},
		{
			spec: specs[2],
			want: []string{"en", "de"},
		},
		{
			spec: specs[3],
			want: 42,
		},
		{
			spec: specs[4],
			want: int64(123),
		},
		{
			spec: specs[5],
			want: 45 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.spec.Name, func(t *testing.T) {
			got, err := flagValue(
				cmd,
				tt.spec,
			)

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestBindFlagSet_UnsupportedKindPanics(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	spec := Spec{
		Name: "wat",
		Kind: Kind(255),
	}

	require.PanicsWithValue(
		t,
		`unsupported parameter kind 255 for "wat"`,
		func() {
			BindFlagSet(
				fs,
				[]Spec{spec},
			)
		},
	)
}

func TestFlagValue_UnsupportedKindReturnsError(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
	}

	got, err := flagValue(
		cmd,
		Spec{
			Name: "wat",
			Kind: Kind(255),
		},
	)

	require.Nil(t, got)

	require.EqualError(
		t,
		err,
		`unsupported parameter kind 255 for "wat"`,
	)
}

func TestBindFlagSet_PreservesUsage(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	BindFlagSet(
		fs,
		[]Spec{
			{
				Name:  "format",
				Kind:  KindString,
				Usage: "File format",
			},
		},
	)

	flag := fs.Lookup("format")

	require.NotNil(t, flag)
	require.Equal(t, "File format", flag.Usage)
}
