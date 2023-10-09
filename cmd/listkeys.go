package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// listkeysCmd represents the listkeys command
var listkeysCmd = &cobra.Command{
	Use:   "listkeys",
	Short: "List the keys in the keychain",
	Long: `List the keys in the keychain with the name and fingerprint.
`,
	Run: func(cmd *cobra.Command, args []string) {
		grid, encrypted, err :=
			GetList([][]string{{"name", "pubkey fingerprint"}})
		if err != nil {
			printErr(
				"error getting list: '%v'\n\n", err)
			os.Exit(1)
		}
		var maxLen1, maxLen2 int
		for i := range grid {
			l := len(grid[i][0])
			if l > maxLen1 {
				maxLen1 = l
			}
			l = len(grid[i][1])
			if l > maxLen2 {
				maxLen2 = l
			}
		}
		header, tail := grid[0], grid[1:]
		grid = append([][]string{header},
			[]string{
				strings.Repeat("-", maxLen1) + " ",
				strings.Repeat("-", maxLen2),
			})
		grid = append(grid, tail...)
		maxLen1++
		printErr(
			"keys in keychain: (* = password protected)\n\n")
		for i := range grid {
			isDefault := "          "
			if grid[i][0] == defaultKey {
				isDefault = " (default)"
			}
			crypted := " "
			if _, ok := encrypted[grid[i][0]]; ok {
				crypted = "*"
			}
			grid[i][0] = grid[i][0] + strings.Repeat(" ",
				maxLen1-len(grid[i][0]))
			printErr(
				"  %s %s%s\n", crypted, grid[i][0], grid[i][1]+isDefault)
		}

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
