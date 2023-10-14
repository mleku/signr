package cmd

import (
	"github.com/mleku/signr/pkg/signr"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen <name>",
	Short: "Generate a new nostr key",
	Long:  `Generates a new key saves it into the application data directory as part of the keychain.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			signr.Fatal("key name is required")
		}
		cfg.Gen(args[0])
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}
