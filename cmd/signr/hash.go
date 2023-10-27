package signr

import (
	"encoding/hex"
	"fmt"

	"github.com/mleku/signr/pkg/signr"
	"github.com/spf13/cobra"
)

var hashCmd = &cobra.Command{
	Use:   "hash <filename>",
	Short: "return the SHA256 hash of a file or from stdin",
	Long: `hashes a named file and prints the hexadecimal 32 byte SHA256 hash.

use the filename '-' to indicate the file is being piped in via | or > or similar.
`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			s.Fatal("file name is required")
		}
		hash, wasHash, err := signr.Hash(args[0])
		if err != nil {
			s.Fatal("error hashing file: %s\n", err)
		}
		if wasHash {
			s.Fatal("file '%s' does not exist and is also a valid 32 byte"+
				" hexadecimal value\n", args[0])
		}
		fmt.Println(hex.EncodeToString(hash))
	},
}

func init() {
	rootCmd.AddCommand(hashCmd)
}
