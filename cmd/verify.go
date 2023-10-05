package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify <signature/sigfile> <signed file>",
	Short: "Verify that a signature matches the provided file",
	Long: `Verifies that the signature on a given file is correct.

Signature can be given on the commandline, if it starts with 'signr' and contains the expected series of '_' separated fields, or otherwise is interpreted as a file containing such a signature.

The signed file can be given as '-' which indicates to read from the standard input as the file is being piped using '<' or '|'.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("verify called")
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// verifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
