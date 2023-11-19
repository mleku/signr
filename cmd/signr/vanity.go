package signr

import (
	"github.com/mleku/signr/pkg/signr"
	"github.com/spf13/cobra"
	"strings"
)

var vanityCmd = &cobra.Command{
	Use:   "vanity <string> <name> [begin|contain|end]",
	Short: "Generate a new vanity nostr key",
	Long: `Generates a new nostr key that begins with, contains, or ends with a given string.

Vanity keys are a kind of proof of work. They take time to generate and are easier to identify for humans.

The longer the <string> the longer it will take to generate it.

If the final position spec is omitted, the search will look for the beginning.
`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			s.Fatal("key match string and name is required\n")
		}
		keyName := args[1]
		str := args[0]
		var err error
		keyName, err = s.Sanitize(keyName)
		if err != nil {
			s.Fatal("key name failed sanitizing: %s\n", keyName)
		}
		if keyName != args[1] {
			s.Fatal("key name input is invalid - sanitizing it changed it, was '%s' changed to '%s'\n",
				args[0], keyName)
		}
		var where signr.Position
		var canonical string
		where = signr.PositionBeginning
		if len(args) >= 3 {
			canonical = strings.ToLower(args[2])
			switch {
			case strings.HasPrefix(canonical, "begin"):
				where = signr.PositionBeginning
			case strings.Contains(canonical, "contain"):
				where = signr.PositionContains
			case strings.HasSuffix(canonical, "end"):
				where = signr.PositionEnding
			}
		}
		if err = s.Vanity(str, keyName, where); err != nil {
			s.Fatal("error: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(vanityCmd)
}
