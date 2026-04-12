package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completions",
	Long: `Generate shell completion scripts for supported shells.

Bash:
  # Source completion in your shell
  source <(oc completion bash)

  # Or install permanently (requires sudo for /etc, or use ~/.local/):
  oc completion bash | sudo tee /etc/bash_completion.d/oc > /dev/null

Zsh:
  # Save completion file
  oc completion zsh > ~/.zsh/completions/_oc

  # Add to ~/.zshrc (MUST be after compinit):
  # Note: The completion file requires compinit to be loaded first.
  # If you see "compdef: command not found", source completion after compinit.
  # Example: source <(oc completion zsh)

  # Alternative: source on the fly (works automatically)
  source <(oc completion zsh)

  # Clear completion cache and restart:
  rm -f ~/.zcompdump && exec zsh

Fish:
  oc completion fish > ~/.config/fish/completions/oc.fish

PowerShell:
  # Add to your PowerShell profile
  oc completion powershell >> $PROFILE

For more details, see: https://github.com/sven1103-agent/opencode-config-cli#shell-completion
`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			return fmt.Errorf("unsupported shell: %s", args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
