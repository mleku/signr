package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	PubExtension = ".pub"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use: "save <name>",
	Short: "Generate and save a new key, " +
		"with a name into your data directory.",
	Long: `Generates and saves a new key pair into your data directory, with a name format like this: ~/.signr/<name>.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			_, _ = fmt.Fprintln(os.Stderr, "name is required")
			os.Exit(1)
		}

		secPath := filepath.Join(dataDir, args[0])
		pubPath := secPath + PubExtension
		fmt.Fprintf(os.Stderr,
			"saving secret key in '%s', public key in '%s'\n",
			secPath, pubPath)

	},
}

func init() {
	genCmd.AddCommand(saveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// saveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// saveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
