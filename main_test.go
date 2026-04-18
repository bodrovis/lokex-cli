package main

import "testing"

func TestRun_Smoke(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no args",
			args:    nil,
			wantErr: false,
		},
		{
			name:    "help",
			args:    []string{"--help"},
			wantErr: false,
		},
		{
			name:    "invalid command",
			args:    []string{"definitely-not-a-command"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := run(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("run(%v) error = %v, wantErr = %v", tt.args, err, tt.wantErr)
			}
		})
	}
}
