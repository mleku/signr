package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

var listkeysCmd = &cobra.Command{
	Use:   "listkeys",
	Short: "List the keys in the keychain",
	Long: `List the keys in the keychain with the name and fingerprint.
`,
	Run: func(cmd *cobra.Command, args []string) {

		grid, encrypted, err := s.GetList([][]string{{"name", "fingerprint"}})
		if err != nil {
			s.Fatal("error getting list: '%v'\n\n", err)
		}
		defaultStr := make(map[bool]string)
		defaultStr[true] = " (default)"
		defaultStr[false] = strings.Repeat(" ", len(defaultStr[true]))
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
		s.Log("keys in keychain: (* = password protected)\n\n")
		cryptedStr := make(map[bool]string)
		cryptedStr[true] = " "
		cryptedStr[false] = "*"
		for _, row := range grid {
			_, clear := encrypted[row[0]]
			s.Log("  %s %s %s\n",
				cryptedStr[!clear],
				PadToLength(row[0], maxLen1),
				row[1]+defaultStr[row[0] == s.DefaultKey],
			)
		}

	},
}

func PadToLength(text string, length int) string {
	pad := length - len(text)
	return text + strings.Repeat(" ", pad)
}

func init() {
	rootCmd.AddCommand(listkeysCmd)
}
