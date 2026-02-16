package internal

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultConfigDir(t *testing.T) {
	configDir := defaultConfigDir()

	if strings.HasPrefix(configDir, string(filepath.Separator)+".config") {
		t.Fatalf("defaultConfigDir() returned root-level path %q", configDir)
	}

	wantSuffix := filepath.Join(".config", "ci-thief")
	if !strings.HasSuffix(configDir, wantSuffix) {
		t.Fatalf("defaultConfigDir() = %q, want suffix %q", configDir, wantSuffix)
	}
}
