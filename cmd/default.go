package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// defaultCmd represents the default command
var defaultCmd = &cobra.Command{
	Use:   "default",
	Short: "set the default key to sign with",
	Long: `sets the default key to sign with if not specified for the sign command.

if the following CLI argument starts with an @ it is interpreted to be the key fingerprint
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			_, _ = fmt.Fprintf(os.Stderr,
				"ERROR: default key must be named.\n\nhere are the options:\n\n")
			listkeysCmd.Run(cmd, args)
			_, _ = fmt.Fprintf(os.Stderr, "\n")
			os.Exit(1)
		}
		grid, _, err := GetList(nil)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr,
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
						_, _ = fmt.Fprintf(os.Stderr,
							"failed to update config: '%v'\n", err)
						return
					}
					_, _ = fmt.Fprintf(os.Stderr,
						"key %s %s now default\n",
						grid[i][0], grid[i][1])
					return
				}
			}
		}
	},
}

func init() {
	setCmd.AddCommand(defaultCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// defaultCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// defaultCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
