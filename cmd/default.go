package cmd

import (
	"github.com/mleku/signr/pkg/signr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var defaultCmd = &cobra.Command{
	Use:   "default",
	Short: "set the default key to sign with",
	Long: `sets the default key to sign with if not specified for the sign command.

if the following CLI argument starts with an @ it is interpreted to be the key fingerprint
`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {

			signr.PrintErr(
				"ERROR: default key must be named.\n\n" +
					"here are the options:\n\n")

			listkeysCmd.Run(cmd, args)
			signr.Fatal("\n")
		}

		grid, _, err := signr.GetList(cfg, nil)
		if err != nil {

			signr.Fatal("ERROR: '%s'\n\n", err)
		}

		for _, row := range grid {

			for j := range row {

				if args[0] == row[j] {

					cfg.DefaultKey = row[0]

					viper.Set("default", cfg.DefaultKey)

					if err = viper.WriteConfig(); err != nil {

						signr.Fatal("failed to update config: '%v'\n", err)
					}

					signr.PrintErr("key %s %s now default\n", row[0], row[1])
					return
				}
			}
		}
	},
}

func init() {

	setCmd.AddCommand(defaultCmd)
}
