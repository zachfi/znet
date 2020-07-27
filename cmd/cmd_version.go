package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Show znet version",
	Long:    "Displays the current built version of the znet command.",
	Example: "znet version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(rootCmd.Use + " " + Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
