package download

import (
	"strings"
	"testing"
)

func TestDownloadParamSpecs_ConfigKeysAreNotBlank(t *testing.T) {
	for _, spec := range downloadParamSpecs {
		if strings.TrimSpace(spec.ConfigKey) == "" {
			t.Fatalf("spec for flag %q has blank ConfigKey", spec.FlagName)
		}
	}
}

func TestDownloadParamSpecs_ConfigKeysAreUnique(t *testing.T) {
	seen := map[string]bool{}

	for _, spec := range downloadParamSpecs {
		key := strings.TrimSpace(spec.ConfigKey)
		if key == "" {
			continue
		}
		if seen[key] {
			t.Fatalf("duplicate config key %q", key)
		}
		seen[key] = true
	}
}
