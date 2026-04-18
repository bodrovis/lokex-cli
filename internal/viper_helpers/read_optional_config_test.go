package viper_helpers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadOptionalConfig_IgnoresMissingImplicitConfig(t *testing.T) {
	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(origWD)
	}()

	emptyWD := t.TempDir()
	if err := os.Chdir(emptyWD); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}

	origUserHomeDir := userHomeDir
	userHomeDir = func() (string, error) {
		return t.TempDir(), nil
	}
	defer func() {
		userHomeDir = origUserHomeDir
	}()

	v := NewConfigViper("", "LOKEX")

	if err := ReadOptionalConfig(v, ""); err != nil {
		t.Fatalf("expected missing implicit config to be ignored, got: %v", err)
	}
}

func TestReadOptionalConfig_ReturnsErrorForMissingExplicitConfig(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "missing.yaml")
	v := NewConfigViper(configFile, "LOKEX")

	err := ReadOptionalConfig(v, configFile)
	if err == nil {
		t.Fatal("expected error for missing explicit config file")
	}
}

func TestReadOptionalConfig_ReturnsErrorForInvalidImplicitConfig(t *testing.T) {
	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(origWD)
	}()

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}

	origUserHomeDir := userHomeDir
	userHomeDir = func() (string, error) {
		return t.TempDir(), nil
	}
	defer func() {
		userHomeDir = origUserHomeDir
	}()

	if err := os.WriteFile(filepath.Join(tmpDir, "lokex.yaml"), []byte(":\nbad yaml"), 0o644); err != nil {
		t.Fatalf("failed to write invalid config file: %v", err)
	}

	v := NewConfigViper("", "LOKEX")

	err = ReadOptionalConfig(v, "")
	if err == nil {
		t.Fatal("expected parse error for invalid implicit config")
	}
}

func TestReadOptionalConfig_ReturnsErrorForInvalidExplicitConfig(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "broken.yaml")
	if err := os.WriteFile(configFile, []byte(":\nbad yaml"), 0o644); err != nil {
		t.Fatalf("failed to write invalid config file: %v", err)
	}

	v := NewConfigViper(configFile, "LOKEX")

	err := ReadOptionalConfig(v, configFile)
	if err == nil {
		t.Fatal("expected parse error for invalid explicit config")
	}
}
