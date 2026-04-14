# Installation

This guide covers how to install Occo on macOS and Linux.

## Prerequisites

- A Unix-like operating system (macOS or Linux)
- Go 1.21+ (for `go install` method)
- `curl` and `tar` (for manual download)

## Install Methods

### Method 1: Go Install (Recommended)

If you have Go installed:

```sh
go install github.com/qbicsoftware/occo@latest
```

This installs the `occo` binary to `$GOPATH/bin` (or `$HOME/go/bin`). Make sure this directory is in your PATH.

**Note:** `@latest` follows Go module version selection. It does not mean "latest GitHub prerelease".

**Install a specific version:**

```sh
go install github.com/qbicsoftware/occo@v1.0.0
```

### Method 2: Manual Download

Download the correct tarball for your platform from [GitHub Releases](https://github.com/qbicsoftware/occo/releases):

```sh
# Example: macOS ARM64
VERSION=v1.0.0
curl -L "https://github.com/qbicsoftware/occo/releases/download/${VERSION}/occo_${VERSION#v}_darwin_arm64.tar.gz" | tar xz
mv occo ~/.local/bin/
```

Set `VERSION` on its own line before running `curl`. Do not collapse this into `VERSION=... curl ...`, because `${VERSION}` is expanded by the shell before that temporary assignment takes effect.

**Available platforms:**
- `darwin_amd64` — macOS Intel
- `darwin_arm64` — macOS Apple Silicon
- `linux_amd64` — Linux x86_64
- `linux_arm64` — Linux ARM64

**Add to PATH:**

```sh
export PATH="$HOME/.local/bin:$PATH"
```

To make this permanent, add the line above to your shell profile (`~/.zshrc` or `~/.bashrc`).

## Verify Installation

```sh
occo version
```

You should see output like `occo version 1.0.0`.

## PATH Issues

If you see "command not found: occo" after installation:

1. **For go install**: Ensure `$HOME/go/bin` is in your PATH
2. **For manual download**: Ensure `~/.local/bin` is in your PATH

```sh
# Add to ~/.zshrc or ~/.bashrc
export PATH="$HOME/go/bin:$HOME/.local/bin:$PATH"
```

## Upgrading

Re-run `go install` with the desired version:

```sh
go install github.com/qbicsoftware/occo@latest
```

Or download a new release manually from [GitHub Releases](https://github.com/qbicsoftware/occo/releases).

## Configuration Location

The CLI stores registered sources in a JSON file. By default it follows the XDG Base Directory specification:

- **Default location**: `~/.config/occo/sources.json`
- **Custom location**: Set `XDG_CONFIG_HOME` environment variable to override

### Legacy Migration

If you used an older version of the CLI (pre-v1.0.0), your sources may be in:
- `~/.config/opencode-helper/sources.json` (shell script version)
- `~/Library/Application Support/opencode-helper/sources.json` (pre-XDG Go CLI)

These will be automatically migrated to the new XDG standard location on first run.

## Shell Completion

Enable tab completion for your shell:

```sh
# Bash: source on the fly, or install permanently
source <(occo completion bash)
occo completion bash | sudo tee /etc/bash_completion.d/occo > /dev/null

# Zsh: source on the fly (recommended)
source <(occo completion zsh)

# Or save to completions dir:
# occo completion zsh > ~/.zsh/completions/_occo
# Add to ~/.zshrc: fpath=(~/.zsh/completions $fpath)
# Clear cache: rm -f ~/.zcompdump && exec zsh

# Fish
occo completion fish > ~/.config/fish/completions/occo.fish
```

Run `occo completion --help` for full instructions.

## Next Steps

- [Configure a config bundle](config-bundles.md)
- [CLI Reference](cli-reference.md)
