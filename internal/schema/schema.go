// Package schema provides functionality for handling OpenCode schemas.
package schema

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidatePath ensures the target path stays within the allowed base directory
// to prevent path traversal attacks.
func ValidatePath(baseDir, targetPath string) error {
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return fmt.Errorf("failed to resolve base directory: %w", err)
	}
	absTarget, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("failed to resolve target path: %w", err)
	}

	absBase = filepath.Clean(absBase)
	absTarget = filepath.Clean(absTarget)

	if !strings.HasPrefix(absTarget, absBase+string(filepath.Separator)) {
		if absTarget != absBase {
			return fmt.Errorf("path traversal detected: %s is not within %s", targetPath, baseDir)
		}
	}

	return nil
}

// SchemaFile represents a schema file.
type SchemaFile struct {
	Name     string
	Filename string
}

// Schemas returns the list of schema files to install.
func Schemas() []SchemaFile {
	return []SchemaFile{
		{Name: "handoff", Filename: "handoff.schema.json"},
		{Name: "result", Filename: "result.schema.json"},
	}
}

// SchemaDir returns the directory name for schemas.
func SchemaDir() string {
	return "schemas"
}

// FindSchema searches for a schema file in common locations.
// Returns the full path if found, or an error if not found.
func FindSchema(name string) (string, error) {
	schemaFile := ""
	for _, s := range Schemas() {
		if s.Name == name {
			schemaFile = s.Filename
			break
		}
	}
	if schemaFile == "" {
		return "", fmt.Errorf("unknown schema: %s", name)
	}

	// Check for bundled schemas in common locations
	possiblePaths := []string{
		filepath.Join(".opencode", "schemas", schemaFile),
		filepath.Join("..", ".opencode", "schemas", schemaFile),
		filepath.Join("..", "..", ".opencode", "schemas", schemaFile),
	}

	// Check executable's directory
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		possiblePaths = append(possiblePaths, filepath.Join(execDir, ".opencode", "schemas", schemaFile))
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("schema not found: %s", name)
}

// InstallSchema copies a schema file to the target directory.
func InstallSchema(srcPath, destDir string, force bool) error {
	// Get filename from source
	filename := filepath.Base(srcPath)
	destPath := filepath.Join(destDir, filename)

	// Ensure destination directory exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	// Check if destination exists
	if !force {
		if _, err := os.Stat(destPath); err == nil {
			return fmt.Errorf("schema already exists: %s (use --force to overwrite)", destPath)
		}
	}

	// Copy the file
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write schema file: %w", err)
	}

	return nil
}

// InstallAll installs all schema files to the target directory.
func InstallAll(targetDir string, force bool) error {
	// Validate target directory to prevent path traversal
	if err := ValidatePath(targetDir, targetDir); err != nil {
		return err
	}

	schemasDir := filepath.Join(targetDir, SchemaDir())

	for _, schema := range Schemas() {
		srcPath, err := FindSchema(schema.Name)
		if err != nil {
			return err
		}

		if err := InstallSchema(srcPath, schemasDir, force); err != nil {
			return err
		}
	}

	return nil
}
