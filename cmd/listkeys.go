package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listkeysCmd represents the listkeys command
var listkeysCmd = &cobra.Command{
	Use:   "listkeys",
	Short: "List the keys in the keychain",
	Long: `List the keys in the keychain with the name and fingerprint.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("listkeys called")
	},
}

func init() {
	rootCmd.AddCommand(listkeysCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listkeysCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listkeysCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
