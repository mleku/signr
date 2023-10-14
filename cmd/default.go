package cmd

import (
	"github.com/mleku/signr/pkg/signr"
	"github.com/spf13/cobra"
)

var defaultCmd = &cobra.Command{
	Use:   "default",
	Short: "set the default key to sign with",
	Long: `sets the default key to sign with if not specified for the sign command.

if the following CLI argument starts with an @ it is interpreted to be the key fingerprint.

either fingerprint or key name can be used to identify the key intended.
`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {

			signr.PrintErr(
				"ERROR: default key must be named.\n\n" +
					"here are the options:\n\n")

			listkeysCmd.Run(cmd, args)
			signr.Fatal("\n")
		}

		if err := cfg.SetDefault(args[0]);err != nil {
			signr.PrintErr("%s\n", err)
		}

	},
}

func init() {

	setCmd.AddCommand(defaultCmd)
}
