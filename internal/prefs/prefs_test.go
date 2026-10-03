package prefs

import (
	"os"
	"path/filepath"
	"testing"
)

// isolate points the config dir at a temp directory for the test.
func isolate(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("HOME", dir)
}

func TestLoadEmpty(t *testing.T) {
	isolate(t)
	if got := Model(); got != "" {
		t.Errorf("Model() = %q, want empty", got)
	}
}

func TestSetAndLoadModel(t *testing.T) {
	isolate(t)

	if err := SetModel("openai/gpt-6-sol"); err != nil {
		t.Fatalf("SetModel: %v", err)
	}
	if got := Model(); got != "openai/gpt-6-sol" {
		t.Errorf("Model() = %q, want openai/gpt-6-sol", got)
	}

	// The file should exist on disk.
	dir, _ := os.UserConfigDir()
	if _, err := os.Stat(filepath.Join(dir, "moki", "prefs.json")); err != nil {
		t.Errorf("prefs file not written: %v", err)
	}
}

func TestSetModelOverwrites(t *testing.T) {
	isolate(t)

	_ = SetModel("a/one")
	_ = SetModel("b/two")
	if got := Model(); got != "b/two" {
		t.Errorf("Model() = %q, want b/two", got)
	}
}
