package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
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
	Example: ". <(znet completion)",
	Run: func(cmd *cobra.Command, args []string) {
		err := rootCmd.GenZshCompletion(os.Stdout)
		if err != nil {
			log.Error(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
