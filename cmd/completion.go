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
  # Source on the fly
  source <(occo completion bash)

  # Or install permanently
  occo completion bash | sudo tee /etc/bash_completion.d/occo > /dev/null

Zsh:
  # Option 1: Source on the fly (recommended)
  # Add to end of ~/.zshrc (after compinit):
  source <(occo completion zsh)

  # Option 2: Save to completions dir
  occo completion zsh > ~/.zsh/completions/_occo
  # Ensure ~/.zshrc has: fpath=(~/.zsh/completions $fpath)
  # Clear completion cache and restart:
  rm -f ~/.zcompdump && exec zsh

Fish:
  occo completion fish > ~/.config/fish/completions/occo.fish

PowerShell:
  occo completion powershell >> $PROFILE

For more details, see: https://github.com/qbicsoftware/occo#shell-completion
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
		}
		return fmt.Errorf("unsupported shell: %s", args[0])
	},
}

func init() {
	completionCmd.Flags().BoolP("help", "h", false, "Help for completion")
	rootCmd.AddCommand(completionCmd)
}
