package viper_helpers

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestBindEnvKeys_BindsAllKeys(t *testing.T) {
	v := NewConfigViper("", "LOKEX")

	t.Setenv("LOKEX_TOKEN", "test-token")
	t.Setenv("LOKEX_PROJECT_ID", "123")

	err := BindEnvKeys(v, []string{"token", "project-id"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got := v.GetString("token"); got != "test-token" {
		t.Fatalf("unexpected token: got %q, want %q", got, "test-token")
	}

	if got := v.GetString("project-id"); got != "123" {
		t.Fatalf("unexpected project-id: got %q, want %q", got, "123")
	}
}

func TestBindEnvKeys_EmptyKeys(t *testing.T) {
	v := viper.New()

	err := BindEnvKeys(v, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestBindEnvKeys_UsesEnvKeyReplacer(t *testing.T) {
	v := viper.New()
	v.SetEnvPrefix("LOKEX")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	v.AutomaticEnv()

	t.Setenv("LOKEX_UPLOAD_FILENAME", "messages.json")

	err := BindEnvKeys(v, []string{"upload.filename"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got := v.GetString("upload.filename"); got != "messages.json" {
		t.Fatalf("unexpected upload.filename: got %q, want %q", got, "messages.json")
	}
}
