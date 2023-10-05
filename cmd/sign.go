package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign <file> [key name]",
	Short: "Generate a signature on a file",
	Long: `Generates a signature for the provided file, using either the default key, or if specified, another key.

Use the filename '-' to indicate that the data is being piped in with < or |.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sign called")
	},
}

func init() {
	rootCmd.AddCommand(signCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// signCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// signCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
