// Package preset provides functionality for handling OpenCode presets.
package preset

import (
	"fmt"
	"os"
	"path/filepath"
)

// ValidPresets returns the list of valid preset names.
func ValidPresets() []string {
	return []string{
		"mixed",
		"openai",
		"big-pickle",
		"minimax",
		"kimi",
	}
}

// PresetFileName returns the filename for a given preset name.
func PresetFileName(name string) string {
	return fmt.Sprintf("opencode.%s.json", name)
}

// FindPreset searches for a preset file in common locations.
// Returns the full path if found, or an error if not found.
func FindPreset(name string) (string, error) {
	// Check for bundled presets relative to the binary
	// In development, check the repo root
	possiblePaths := []string{
		filepath.Join(".", PresetFileName(name)),
		filepath.Join("..", PresetFileName(name)),
	}

	// Check executable's directory
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		possiblePaths = append(possiblePaths, filepath.Join(execDir, PresetFileName(name)))
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("preset not found: %s", name)
}

// CopyPreset copies a preset file to the specified destination.
func CopyPreset(srcPath, destPath string, force bool) error {
	// Check if destination exists
	if !force {
		if _, err := os.Stat(destPath); err == nil {
			return fmt.Errorf("output file exists: %s (use --force to overwrite)", destPath)
		}
	}

	// Ensure destination directory exists
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	// Copy the file
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read preset file: %w", err)
	}

	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write preset file: %w", err)
	}

	return nil
}
