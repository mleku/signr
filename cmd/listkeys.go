package cmd

import (
	"github.com/mleku/signr/pkg/signr"
	"strings"

	"github.com/spf13/cobra"
)

var listkeysCmd = &cobra.Command{
	Use:   "listkeys",
	Short: "List the keys in the keychain",
	Long: `List the keys in the keychain with the name and fingerprint.
`,
	Run: func(cmd *cobra.Command, args []string) {

		grid, encrypted, err :=
			signr.GetList(cfg, [][]string{{"name", "pubkey fingerprint"}})
		if err != nil {

			signr.Fatal("error getting list: '%v'\n\n", err)
		}

		// get the maximum width of the columns
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

		// split header from the grid
		header, tail := grid[0], grid[1:]

		// add separators to the columns.
		grid = append([][]string{header},
			[]string{
				strings.Repeat("-", maxLen1) + " ",
				strings.Repeat("-", maxLen2),
			},
		)

		// add the rows after the spacers
		grid = append(grid, tail...)
		maxLen1++

		signr.PrintErr("keys in keychain: (* = password protected)\n\n")

		for _, row := range grid {

			// show default item
			isDefault := "          "
			if row[0] == cfg.DefaultKey {

				isDefault = " (default)"
			}

			crypted := " "
			if _, ok := encrypted[row[0]]; ok {
				crypted = "*"
			}

			row[0] += strings.Repeat(" ", maxLen1-len(row[0]))

			signr.PrintErr(
				"  %s %s%s\n", crypted, row[0], row[1]+isDefault)
		}

	},
}

func init() {

	rootCmd.AddCommand(listkeysCmd)
}
