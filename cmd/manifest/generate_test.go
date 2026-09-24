package manifest

import (
	"testing"

	"github.com/bodrovis/lokex-cli/internal/appstate"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestValidateGenerateConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     *GenerateConfig
		wantErr string
	}{
		{
			name:    "nil config",
			cfg:     nil,
			wantErr: "manifest generate config is nil",
		},
		{
			name: "valid",
			cfg: &GenerateConfig{
				Paths:           new([]string{"./locales"}),
				NamePattern:     new("{lang}/{name}.{ext}"),
				FilenamePattern: new("{name}.{ext}"),
				Out:             new("./lokex-manifest.json"),
			},
		},
		{
			name: "missing path",
			cfg: &GenerateConfig{
				NamePattern:     new("{lang}/{name}.{ext}"),
				FilenamePattern: new("{name}.{ext}"),
				Out:             new("./lokex-manifest.json"),
			},
			wantErr: "path is required",
		},
		{
			name: "empty path list",
			cfg: &GenerateConfig{
				Paths:           new([]string{}),
				NamePattern:     new("{lang}/{name}.{ext}"),
				FilenamePattern: new("{name}.{ext}"),
				Out:             new("./lokex-manifest.json"),
			},
			wantErr: "path is required",
		},
		{
			name: "missing name pattern",
			cfg: &GenerateConfig{
				Paths:           new([]string{"./locales"}),
				FilenamePattern: new("{name}.{ext}"),
				Out:             new("./lokex-manifest.json"),
			},
			wantErr: "name-pattern is required",
		},
		{
			name: "whitespace name pattern",
			cfg: &GenerateConfig{
				Paths:           new([]string{"./locales"}),
				NamePattern:     new("   "),
				FilenamePattern: new("{name}.{ext}"),
				Out:             new("./lokex-manifest.json"),
			},
			wantErr: "name-pattern is required",
		},
		{
			name: "whitespace only paths",
			cfg: &GenerateConfig{
				Paths:           new([]string{"   ", "\t"}),
				NamePattern:     new("{lang}/{name}.{ext}"),
				FilenamePattern: new("{name}.{ext}"),
				Out:             new("./lokex-manifest.json"),
			},
			wantErr: "path is required",
		},
		{
			name: "missing filename pattern",
			cfg: &GenerateConfig{
				Paths:       new([]string{"./locales"}),
				NamePattern: new("{lang}/{name}.{ext}"),
				Out:         new("./lokex-manifest.json"),
			},
			wantErr: "filename-pattern is required",
		},
		{
			name: "whitespace filename pattern",
			cfg: &GenerateConfig{
				Paths:           new([]string{"./locales"}),
				NamePattern:     new("{lang}/{name}.{ext}"),
				FilenamePattern: new("   "),
				Out:             new("./lokex-manifest.json"),
			},
			wantErr: "filename-pattern is required",
		},
		{
			name: "missing out",
			cfg: &GenerateConfig{
				Paths:           new([]string{"./locales"}),
				NamePattern:     new("{lang}/{name}.{ext}"),
				FilenamePattern: new("{name}.{ext}"),
			},
			wantErr: "out is required",
		},
		{
			name: "whitespace out",
			cfg: &GenerateConfig{
				Paths:           new([]string{"./locales"}),
				NamePattern:     new("{lang}/{name}.{ext}"),
				FilenamePattern: new("{name}.{ext}"),
				Out:             new("   "),
			},
			wantErr: "out is required",
		},
		{
			name: "valid with empty and non-empty paths",
			cfg: &GenerateConfig{
				Paths: new([]string{
					"   ",
					"./locales",
				}),
				NamePattern:     new("{lang}/{name}.{ext}"),
				FilenamePattern: new("{name}.{ext}"),
				Out:             new("./lokex-manifest.json"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateGenerateConfig(tt.cfg)

			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.wantErr)
		})
	}
}

func TestGenerateCommand_PreRunE_DoesNotRequireProjectAccess(t *testing.T) {
	v := viper.New()

	v.Set(
		"manifest.generate.path",
		[]string{"./locales"},
	)

	v.Set(
		"manifest.generate.name-pattern",
		"{lang}/{name}.{ext}",
	)

	state := &appstate.State{
		Viper: v,
	}

	cmd := newGenerateCommand(state)

	require.NoError(
		t,
		cmd.PreRunE(cmd, nil),
	)
}
