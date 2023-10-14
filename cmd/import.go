package cmd

import (
	"github.com/mleku/signr/pkg/signr"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import <secret key> <name>",
	Short: "Import a secret key",
	Long: `Import a secret key provided in hexadecimal and nostr nsec formats.
`,
	Run: func(cmd *cobra.Command, args []string) {

		argLen := len(args)
		if argLen == 1 {

			signr.Fatal("a key name is required after the secret key")
		}
		if err := cfg.Import(args[0], args[1]); err != nil {
			signr.Fatal("ERROR: while importing: '%s'\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
