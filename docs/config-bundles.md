# Config Bundles

Config bundles are versioned, schema-validated OpenCode configuration packages distributed via GitHub releases.

## What is a Bundle?

A bundle is a tar archive containing:
- `opencode-bundle.manifest.json` — Bundle metadata
- Preset JSON files (e.g., `opencode.mixed.json`, `opencode.openai.json`)
- Optional: `.opencode/schemas/` for session artifacts

## Available Bundles

| Bundle | Description |
|--------|-------------|
| [qbicsoftware/opencode-config-bundle](https://github.com/qbicsoftware/opencode-config-bundle) | Official bundle with multiple presets |

## Using a Bundle

### 1. Register a Source

```sh
oc source add qbicsoftware/opencode-config-bundle --name qbic
```

Source formats:
- `owner/repo` — e.g., `qbicsoftware/opencode-config-bundle`
- `github.com/owner/repo`
- `https://github.com/owner/repo/releases/tag/v1.0.0`

### 2. List Available Presets

```sh
oc preset list --sources
```

### 3. Apply a Preset

```sh
oc bundle apply qbic --version v1.0.0 --preset mixed --project-root .
oc bundle apply qbic --preset mixed --project-root .
```

**Prompt files:** If the preset declares `prompt_files` in the manifest, those files are automatically installed to `.opencode/prompts/` in your project. This allows agent configs to reference them using `{file:.opencode/prompts/filename.md}`.

**Controlling asset installation:**
- `--assets=true` (default): Install prompt files automatically
- `--assets=false`: Skip prompt file installation
- `--force`: Overwrite existing files

If a source publishes multiple versions, the CLI should support:

- explicit version selection with `--version <tag>`
- `latest` for the latest stable release
- interactive version selection in a TTY, ordered newest first and clearly labeling prereleases

## Creating Your Own Bundle

### 1. Create the Manifest

Create `opencode-bundle.manifest.json` at your bundle root:

```json
{
  "manifest_version": "1.0.0",
  "bundle_name": "my-bundle",
  "bundle_version": "v1.0.0",
  "presets": [
    {
      "name": "openai",
      "description": "OpenAI-based configuration",
      "entrypoint": "opencode.openai.json",
      "prompt_files": []
    }
  ]
}
```

### 2. Required Fields

| Field | Description |
|-------|-------------|
| `manifest_version` | Semantic version (e.g., `1.0.0`) |
| `bundle_name` | Stable identifier (lowercase, hyphens allowed) |
| `bundle_version` | Release tag |
| `presets` | Array of preset objects |

### 3. Preset Descriptor Fields

| Field | Description |
|-------|-------------|
| `name` | Stable preset ID |
| `description` | Short description |
| `entrypoint` | Path to preset JSON file |
| `prompt_files` | Array of prompt file paths (can be empty) |

### 4. Add Prompt Files (Optional)

If your preset includes agent configurations that reference custom prompts, include those in your bundle:

```
my-bundle/
├── opencode-bundle.manifest.json
├── opencode.default.json
└── prompts/
    ├── coder.md
    └── reviewer.md
```

Reference them in your preset config:

```json
{
  "agent": {
    "coder": {
      "prompt": "{file:.opencode/prompts/coder.md}"
    }
  }
}
```

### 5. Publish as GitHub Release

1. Create a GitHub release. Stable releases and prereleases are both valid bundle sources.
2. Attach a `.tar.gz` bundle archive containing:
   - `opencode-bundle.manifest.json` at root, or under a single top-level directory
   - All preset files referenced by the manifest
   - All prompt files referenced in `prompt_files` arrays
   - `.opencode/schemas/` when your bundle ships the canonical session schemas
3. Attach a matching `-checksums.txt` file using SHA-256.
4. Do not rely on GitHub's auto-generated source tarballs as the primary consumable bundle artifact.

See `docs/specs/bundle-contract.md` for the full GitHub-release distribution contract and a copy-pasteable GitHub Actions example workflow.

### 6. Register Your Bundle

```sh
oc source add your-username/your-bundle --name mybundle
oc bundle apply mybundle --preset <preset-name> --project-root .
