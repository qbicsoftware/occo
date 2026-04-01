package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [shell]",
	Short: "Generate shell completion scripts",
	Long: `Generate the autocompletion script for oc for the specified shell.

To load completions:

Bash:

  $ source <(oc completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ oc completion bash > /etc/bash_completion.d/oc
  # macOS:
  $ oc completion bash > /usr/local/etc/bash_completion.d/oc

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ oc completion zsh > "${fpath[1]}/_oc"

  # You will need to start a new shell for this setup to take effect.

Fish:

  $ oc completion fish | source

  # To load completions for each session, execute once:
  $ oc completion fish > ~/.config/fish/completions/oc.fish

PowerShell:

  $ oc completion powershell | Out-String | Invoke-Expression

  # To load completions for every session, add the output of the above command
  # to your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
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
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
