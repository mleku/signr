package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values from the CLI",
	Long:  `Configuration values can be changed via the CLI, see the subcommands of this command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stderr, "no options given")
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
