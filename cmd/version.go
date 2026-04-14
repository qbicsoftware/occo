package cmd

import (
	"fmt"

	"github.com/qbicsoftware/occo/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print the version information for occo.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("occo %s\n", version.Version)
	},
}
