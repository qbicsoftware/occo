package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sven1103-agent/opencode-config-cli/internal/bundle"
	"github.com/sven1103-agent/opencode-config-cli/internal/source"
)

// setupTestProject creates a temporary project directory
func setupTestProject(t *testing.T) string {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "opencode-test-project-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Create .opencode directory
	opencodeDir := filepath.Join(tempDir, ".opencode")
	if err := os.MkdirAll(opencodeDir, 0755); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("failed to create .opencode directory: %v", err)
	}

	return tempDir
}

// TestBundleApplyNoSource tests applying bundle without a source
func TestBundleApplyNoSource(t *testing.T) {
	// Save original flag values
	origPreset := bundlePreset
	origProjectRoot := bundleProjectRoot
	origForce := bundleForce
	origDryRun := bundleDryRun
	origOutput := bundleOutput
	defer func() {
		bundlePreset = origPreset
		bundleProjectRoot = origProjectRoot
		bundleForce = origForce
		bundleDryRun = origDryRun
		bundleOutput = origOutput
	}()

	// Test with nonexistent source
	bundlePreset = "test"
	bundleProjectRoot = "."
	bundleDryRun = false

	err := runBundleApply("nonexistent-id")
	if err == nil {
		t.Error("runBundleApply() expected error for nonexistent source")
	}
}

// TestBundleApplyMissingPreset tests applying with missing preset flag
func TestBundleApplyMissingPreset(t *testing.T) {
	origPreset := bundlePreset
	defer func() { bundlePreset = origPreset }()

	bundlePreset = ""

	err := runBundleApply("abc12345")
	if err == nil {
		t.Error("runBundleApply() expected error when preset is missing")
	}
}

// TestBundleStatusNoProvenance tests status command with no provenance
func TestBundleStatusNoProvenance(t *testing.T) {
	origProjectRoot := bundleProjectRoot
	defer func() { bundleProjectRoot = origProjectRoot }()

	// Use a temp directory with no provenance
	tempDir, err := os.MkdirTemp("", "opencode-test-noprovenance-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	bundleProjectRoot = tempDir

	err = runBundleStatus()
	if err == nil {
		t.Error("runBundleStatus() expected error when no provenance exists")
	}
}

// TestBundleStatusWithProvenance tests status command with provenance
func TestBundleStatusWithProvenance(t *testing.T) {
	origProjectRoot := bundleProjectRoot
	defer func() { bundleProjectRoot = origProjectRoot }()

	// Create temp project with provenance
	tempDir := setupTestProject(t)
	defer os.RemoveAll(tempDir)

	prov := &bundle.Provenance{
		SourceID:      "test-id",
		SourceName:    "test-source",
		SourceType:    "local-directory",
		BundleVersion: "v1.0.0",
		PresetName:    "test",
		Entrypoint:    "test.json",
		AppliedAt:     "2026-03-31T00:00:00Z",
	}
	if err := bundle.SaveProvenance(tempDir, prov, false); err != nil {
		t.Fatalf("failed to save provenance: %v", err)
	}

	bundleProjectRoot = tempDir

	err := runBundleStatus()
	if err != nil {
		t.Errorf("runBundleStatus() error = %v", err)
	}
}

// TestBundleUpdateNonGitHub tests update command with non-github source
func TestBundleUpdateNonGitHub(t *testing.T) {
	// This test requires a source in the registry
	// For now just verify it returns error for non-github source

	err := runBundleUpdate("nonexistent")
	if err == nil {
		t.Error("runBundleUpdate() expected error for nonexistent source")
	}
}

func TestBundleApplyPassesVersionForGitHubSources(t *testing.T) {
	restore := saveRegistry(t)
	defer restore()

	registry, _ := source.LoadRegistry()
	registry.Sources = []source.Source{{
		ID:       "github1",
		Location: "qbicsoftware/opencode-config-bundle",
		Type:     source.SourceTypeGitHubRelease,
		Name:     "qbic",
	}}
	if err := source.SaveRegistry(registry); err != nil {
		t.Fatalf("failed to save registry: %v", err)
	}

	projectRoot := setupTestProject(t)
	defer os.RemoveAll(projectRoot)

	origPreset := bundlePreset
	origProjectRoot := bundleProjectRoot
	origVersion := bundleVersion
	origResolver := bundleResolveToLocal
	defer func() {
		bundlePreset = origPreset
		bundleProjectRoot = origProjectRoot
		bundleVersion = origVersion
		bundleResolveToLocal = origResolver
	}()

	bundlePreset = "test"
	bundleProjectRoot = projectRoot
	bundleVersion = "v1.2.3"
	bundleResolveToLocal = func(sourceType, sourceLocation, versionTag string) (string, func(), error) {
		if sourceType != "github-release" {
			t.Fatalf("sourceType = %q, want github-release", sourceType)
		}
		if sourceLocation != "qbicsoftware/opencode-config-bundle" {
			t.Fatalf("sourceLocation = %q", sourceLocation)
		}
		if versionTag != "v1.2.3" {
			t.Fatalf("versionTag = %q, want v1.2.3", versionTag)
		}

		bundleRoot := t.TempDir()
		manifest := `{"manifest_version":"1.0.0","bundle_name":"qbic","bundle_version":"v1.2.3","presets":[{"name":"test","entrypoint":"test.json"}]}`
		if err := os.WriteFile(filepath.Join(bundleRoot, "opencode-bundle.manifest.json"), []byte(manifest), 0644); err != nil {
			return "", nil, err
		}
		if err := os.WriteFile(filepath.Join(bundleRoot, "test.json"), []byte(`{"agents":[]}`), 0644); err != nil {
			return "", nil, err
		}
		return bundleRoot, func() {}, nil
	}

	if err := runBundleApply("github1"); err != nil {
		t.Fatalf("runBundleApply() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(projectRoot, "opencode.json")); err != nil {
		t.Fatalf("expected opencode.json to be written: %v", err)
	}
}

// TestBundleApplyFlags tests that bundle apply flags are properly configured
func TestBundleApplyFlags(t *testing.T) {
	if bundleApplyCmd.Flags().Lookup("preset") == nil {
		t.Error("preset flag should exist on bundle apply command")
	}
	if bundleApplyCmd.Flags().Lookup("project-root") == nil {
		t.Error("project-root flag should exist on bundle apply command")
	}
	if bundleApplyCmd.Flags().Lookup("force") == nil {
		t.Error("force flag should exist on bundle apply command")
	}
	if bundleApplyCmd.Flags().Lookup("dry-run") == nil {
		t.Error("dry-run flag should exist on bundle apply command")
	}
}

// TestBundleStatusFlags tests that bundle status flags are properly configured
func TestBundleStatusFlags(t *testing.T) {
	if bundleStatusCmd.Flags().Lookup("project-root") == nil {
		t.Error("project-root flag should exist on bundle status command")
	}
}

// TestBundleUpdateFlags tests that bundle update flags are properly configured
func TestBundleUpdateFlags(t *testing.T) {
	if bundleUpdateCmd.Flags().Lookup("yes") == nil {
		t.Error("yes flag should exist on bundle update command")
	}
}

func TestBundleApplyVersionFlagExists(t *testing.T) {
	if bundleApplyCmd.Flags().Lookup("version") == nil {
		t.Fatal("version flag should exist on bundle apply command")
	}
}

func TestBundleApplyRejectsVersionForLocalSources(t *testing.T) {
	restore := saveRegistry(t)
	defer restore()

	bundleDir := t.TempDir()
	manifest := `{"manifest_version":"1.0.0","bundle_name":"local","bundle_version":"v1.0.0","presets":[{"name":"test","entrypoint":"test.json"}]}`
	if err := os.WriteFile(filepath.Join(bundleDir, "opencode-bundle.manifest.json"), []byte(manifest), 0644); err != nil {
		t.Fatalf("failed to write manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bundleDir, "test.json"), []byte(`{"agents":[]}`), 0644); err != nil {
		t.Fatalf("failed to write preset: %v", err)
	}

	registry, _ := source.LoadRegistry()
	registry.Sources = []source.Source{{
		ID:       "local1",
		Location: bundleDir,
		Type:     source.SourceTypeLocalDirectory,
		Name:     "local",
	}}
	if err := source.SaveRegistry(registry); err != nil {
		t.Fatalf("failed to save registry: %v", err)
	}

	projectRoot := setupTestProject(t)
	defer os.RemoveAll(projectRoot)

	origPreset := bundlePreset
	origProjectRoot := bundleProjectRoot
	origVersion := bundleVersion
	defer func() {
		bundlePreset = origPreset
		bundleProjectRoot = origProjectRoot
		bundleVersion = origVersion
	}()

	bundlePreset = "test"
	bundleProjectRoot = projectRoot
	bundleVersion = "v1.2.3"

	err := runBundleApply("local1")
	if err == nil {
		t.Fatal("runBundleApply() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "--version is only supported for github-release sources") {
		t.Fatalf("runBundleApply() error = %v", err)
	}
}
