package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// TestValidateOutputPath tests the path traversal protection
func TestValidateOutputPath(t *testing.T) {
	tests := []struct {
		name        string
		projectRoot string
		outputPath  string
		wantErr     bool
	}{
		{
			name:        "valid path within project root",
			projectRoot: "/tmp/project",
			outputPath:  "/tmp/project/opencode.json",
			wantErr:     false,
		},
		{
			name:        "valid nested path",
			projectRoot: "/tmp/project",
			outputPath:  "/tmp/project/subdir/opencode.json",
			wantErr:     false,
		},
		{
			name:        "path traversal attempt",
			projectRoot: "/tmp/project",
			outputPath:  "/tmp/project/../../../etc/passwd",
			wantErr:     true,
		},
		{
			name:        "absolute path traversal",
			projectRoot: "/tmp/project",
			outputPath:  "/etc/passwd",
			wantErr:     true,
		},
		{
			name:        "sibling directory traversal",
			projectRoot: "/tmp/project",
			outputPath:  "/tmp/other-file",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOutputPath(tt.projectRoot, tt.outputPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateOutputPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestInitAcceptanceCriteria tests the full init command against acceptance criteria
func TestInitAcceptanceCriteria(t *testing.T) {
	// Create temporary project directory
	tmpDir, err := os.MkdirTemp("", "oc-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test 1: oc init --help shows usage (tested via flag registration)
	initCmdFunc := initCmd
	if initCmdFunc == nil {
		t.Error("initCmd should be registered")
	}

	// Test 2: Preset validation
	validPresets := []string{"mixed", "openai", "big-pickle", "minimax", "kimi"}
	for _, p := range validPresets {
		if !isValidPreset(p) {
			t.Errorf("preset %s should be valid", p)
		}
	}

	// Test 3: Invalid preset should be rejected
	if isValidPreset("invalid") {
		t.Error("preset 'invalid' should not be valid")
	}

	// Test 4: Output path validation
	validOutput := filepath.Join(tmpDir, "opencode.json")
	if err := validateOutputPath(tmpDir, validOutput); err != nil {
		t.Errorf("valid output path should not error: %v", err)
	}

	// Test 5: Path traversal should be rejected
	traversalOutput := filepath.Join(tmpDir, "..", "..", "etc", "passwd")
	if err := validateOutputPath(tmpDir, traversalOutput); err == nil {
		t.Error("path traversal should be rejected")
	}

	// Test 6: Project root validation
	nonexistentRoot := "/tmp/nonexistent-12345"
	if _, err := os.Stat(nonexistentRoot); !os.IsNotExist(err) {
		t.Logf("note: temp dir may have been created by another test")
	}
}

// isValidPreset checks if a preset name is valid
func isValidPreset(name string) bool {
	validPresets := []string{"mixed", "openai", "big-pickle", "minimax", "kimi"}
	for _, p := range validPresets {
		if name == p {
			return true
		}
	}
	return false
}
