package preset

import (
	"os"
	"path/filepath"
	"testing"
)

// TestValidPresets tests that all expected presets are valid
func TestValidPresets(t *testing.T) {
	presets := ValidPresets()
	expected := []string{"mixed", "openai", "big-pickle", "minimax", "kimi"}

	if len(presets) != len(expected) {
		t.Errorf("expected %d presets, got %d", len(expected), len(presets))
	}

	for _, e := range expected {
		found := false
		for _, p := range presets {
			if p == e {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected preset %q not found", e)
		}
	}
}

// TestPresetFileName tests preset file name generation
func TestPresetFileName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"mixed", "opencode.mixed.json"},
		{"openai", "opencode.openai.json"},
		{"big-pickle", "opencode.big-pickle.json"},
		{"minimax", "opencode.minimax.json"},
		{"kimi", "opencode.kimi.json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PresetFileName(tt.name); got != tt.want {
				t.Errorf("PresetFileName(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

// TestCopyPreset tests preset copying functionality
func TestCopyPreset(t *testing.T) {
	// Create a temp file to copy
	tmpDir, err := os.MkdirTemp("", "oc-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	srcFile := filepath.Join(tmpDir, "source.json")
	destFile := filepath.Join(tmpDir, "dest.json")

	// Write test content
	if err := os.WriteFile(srcFile, []byte(`{"test": true}`), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	// Test copy without force (new file)
	if err := CopyPreset(srcFile, destFile, false); err != nil {
		t.Errorf("CopyPreset() error = %v", err)
	}

	// Verify content
	content, err := os.ReadFile(destFile)
	if err != nil {
		t.Errorf("failed to read dest file: %v", err)
	}
	if string(content) != `{"test": true}` {
		t.Errorf("content mismatch: got %s", string(content))
	}

	// Test copy with force (existing file)
	if err := CopyPreset(srcFile, destFile, true); err != nil {
		t.Errorf("CopyPreset() with force error = %v", err)
	}

	// Test copy without force (existing file should fail)
	if err := CopyPreset(srcFile, destFile, false); err == nil {
		t.Error("CopyPreset() should fail when file exists and force=false")
	}
}
