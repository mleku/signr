package signr

import (
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values from the CLI",
	Long:  `Configuration values can be changed via the CLI, see the subcommands of this command`,
	Run: func(cmd *cobra.Command, args []string) {
		s.Err("ERROR: no options given.\n\n")
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
