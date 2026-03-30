package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sven1103-agent/opencode-helper/internal/preset"
	"github.com/sven1103-agent/opencode-helper/internal/schema"
)

// validateOutputPath ensures the output path stays within the project root
// to prevent path traversal attacks.
func validateOutputPath(projectRoot, outputPath string) error {
	// Resolve both paths to absolute
	absProjectRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return fmt.Errorf("failed to resolve project root: %w", err)
	}
	absOutputPath, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %w", err)
	}

	// Clean both paths
	absProjectRoot = filepath.Clean(absProjectRoot)
	absOutputPath = filepath.Clean(absOutputPath)

	// Check if output path starts with project root
	if !strings.HasPrefix(absOutputPath, absProjectRoot+string(filepath.Separator)) {
		// Also check if they're equal (edge case for same directory)
		if absOutputPath != absProjectRoot {
			return fmt.Errorf("invalid output path: path traversal detected (output must be within project root)")
		}
	}

	return nil
}

var (
	initProjectRoot string
	initPreset      string
	initOutput      string
	initForce       bool
	initDryRun      bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a project with an OpenCode preset and install schemas",
	Long: `Initialize a project by:
1. Copying a preset file to the project root
2. Installing schemas to .opencode/schemas/

Default preset is "mixed", writing to opencode.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit()
	},
}

func init() {
	initCmd.Flags().StringVar(&initProjectRoot, "project-root", ".", "Project root directory")
	initCmd.Flags().StringVar(&initPreset, "preset", "mixed", "Preset name to use")
	initCmd.Flags().StringVar(&initOutput, "output", "opencode.json", "Output file path")
	initCmd.Flags().BoolVar(&initForce, "force", false, "Overwrite existing files")
	initCmd.Flags().BoolVar(&initDryRun, "dry-run", false, "Show what would be done without doing it")

	rootCmd.AddCommand(initCmd)
}

func runInit() error {
	// Validate preset
	validPresets := preset.ValidPresets()
	valid := false
	for _, p := range validPresets {
		if initPreset == p {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid preset: %s (valid: %v)", initPreset, validPresets)
	}

	// Resolve project root
	projectRoot, err := filepath.Abs(initProjectRoot)
	if err != nil {
		return fmt.Errorf("invalid project root: %w", err)
	}

	// Check if project root exists
	if _, err := os.Stat(projectRoot); os.IsNotExist(err) {
		return fmt.Errorf("project root does not exist: %s", projectRoot)
	}

	// Resolve output path
	outputPath := filepath.Join(projectRoot, initOutput)

	// Validate output path to prevent path traversal
	if err := validateOutputPath(projectRoot, outputPath); err != nil {
		return err
	}

	// Find preset
	presetPath, err := preset.FindPreset(initPreset)
	if err != nil {
		return fmt.Errorf("failed to find preset: %w", err)
	}

	// Dry run mode
	if initDryRun {
		fmt.Printf("dry-run: copy %s -> %s\n", presetPath, outputPath)
		fmt.Printf("dry-run: install schemas to %s/.opencode/schemas/\n", projectRoot)
		return nil
	}

	// Apply preset
	if err := preset.CopyPreset(presetPath, outputPath, initForce); err != nil {
		return fmt.Errorf("failed to apply preset: %w", err)
	}
	fmt.Printf("written: %s\n", outputPath)

	// Install schemas
	opencodeDir := filepath.Join(projectRoot, ".opencode")
	if err := schema.InstallAll(opencodeDir, initForce); err != nil {
		return fmt.Errorf("failed to install schemas: %w", err)
	}
	fmt.Printf("written: %s/.opencode/schemas/\n", projectRoot)

	fmt.Println("done: init complete")

	return nil
}
