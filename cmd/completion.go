package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates zsh completion scripts",
	Long: `To load completion run

. <(znet completion)

To configure your zsh shell to load completions for each session add to your zsh

# ~/.zshrc or ~/.profile
. <(znet completion)
`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenZshCompletion(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
