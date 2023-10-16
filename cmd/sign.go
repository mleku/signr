package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	Hex, OnlySig bool
	Pass, Custom string
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign [flags] <file/hash> [key name]",
	Short: "Generate a signature on a file or hash",
	Long: `Generates a signature for the provided file, using either the default key, or if specified, another key.

Use the filename '-' to indicate that the data is being piped in with < or |.

If filename is in fact a 64 character hexadecimal value, it will be signed on without the nonce and public key, and yield only a bech32 signature, or hexadecimal with the --hex flag.

The sigonly option is the same as the hex option except the output signature is in nostr nsig Bech32 format.
`,
	Run: func(cmd *cobra.Command, args []string) {

		signature, err := s.Sign(args, Pass, Custom, Hex, OnlySig)
		if err != nil {

			s.Err("ERROR: while signing: %s\n", err)
		} else {

			fmt.Println(signature)
		}
		return

	},
}

func init() {

	signCmd.PersistentFlags().StringVarP(&Pass, "pass", "p", "",
		"password to unlock the key - for better security, use the " +
		"environment variable")
	signCmd.PersistentFlags().StringVarP(&Custom, "custom", "k", "",
		"custom additional namespace")
	signCmd.PersistentFlags().BoolVarP(&Hex, "hex", "x", false,
		"print signature in hex - this also applies the same effect as sigonly")
	signCmd.PersistentFlags().BoolVarP(&OnlySig, "sigonly", "s", false,
		"print only signature - note: this also omits the adding of a " +
		"nonce as a verifier could not know it")
	rootCmd.AddCommand(signCmd)
}
