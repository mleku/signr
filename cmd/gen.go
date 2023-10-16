package cmd

import (
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen <name>",
	Short: "Generate a new nostr key",
	Long:  `Generates a new key saves it into the application data directory as part of the keychain.

the name should be relevant to the purpose of the key, cannot contain any non-printable characters or white space, and any such name will be rejected.
`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			s.Fatal("key name is required")
		}

		keyName := args[0]

		var err error
		keyName, err = s.Sanitize(keyName)
		if err != nil {
			s.Fatal("key name failed sanitizing: %s\n", keyName)
		}

		if keyName == args[0] {

			s.Fatal("key name input is invalid - sanitizing it changed it\n")
		}

		s.Gen(keyName)
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}
