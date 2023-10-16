package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a named key",
	Long: `delete a key - for safety reasons, it does not delete the key, but instead changes its name.

to actually delete a key, you must manually delete it in the filesystem. the files are written read only so the filesystem will double check you want to do it.
`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args)<1 {

			s.Fatal("key name is required\n")
		}

		err := s.Delete(args[0])
		if err != nil {
			s.PrintErr("error deleting key '%s': %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
