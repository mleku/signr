package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var defaultCmd = &cobra.Command{
	Use:   "default",
	Short: "set the default key to sign with",
	Long: `sets the default key to sign with if not specified for the sign command.

if the following CLI argument starts with an @ it is interpreted to be the key fingerprint
`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			PrintErr(
				"ERROR: default key must be named.\n\nhere are the options:\n\n")
			listkeysCmd.Run(cmd, args)
			PrintErr("\n")
			os.Exit(1)
		}

		grid, _, err := GetList(nil)
		if err != nil {
			PrintErr(
				"ERROR: '%s'\n\n", err)
			os.Exit(1)
		}

		for i := range grid {
			for j := range grid[i] {

				if args[0] == grid[i][j] {

					defaultKey = grid[i][0]

					viper.Set("default", defaultKey)

					err = viper.WriteConfig()
					if err != nil {
						PrintErr("failed to update config: '%v'\n", err)
						return
					}
					PrintErr("key %s %s now default\n",
						grid[i][0], grid[i][1])
					return
				}
			}
		}
	},
}

func init() {
	setCmd.AddCommand(defaultCmd)
}
