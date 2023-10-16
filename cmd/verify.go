package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PubKey string

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify <file> <signature file/signature>",
	Short: "check that a file matches a signature",
	Long: `checks that a signature is valid for a given file.

the public key is embedded in the signature should also match the key expected to be used on the signature.

use the filename '-' to indicate the file is being piped in via stdin.

if the signature is a signature-only, whether as a parameter or in the referenced file, the pubkey must be provided by parameter or environment variable.

if the consuming protocol requires an additional custom namespace, and was used when making a signature only, it must be passed in to construct the correct signing material for the message hash to check the signature against.		
`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			s.Fatal("ERROR: at minimum a file and a keyfile name must be specified\n\n")
		}

		filename := args[0]
		sigOrSigFile := args[1]

		s.Log("pubkey input: %s\n", PubKey)

		// if we didn't get a pubkey string, scan env for it.
		if PubKey == "" {

			PubKey = viper.GetString("pubkey")
			s.Log("pubkey from env: %s\n", PubKey)

		}

		var valid bool
		var err error
		valid, err = s.Verify(filename, sigOrSigFile, PubKey, Custom)

		if err != nil {

			s.Fatal("error verifying signature: %s\n", err)
		}

		var validity = map[bool]string{true: "VALID", false: "INVALID"}

		fmt.Println(validity[valid])
	},
}

func init() {

	verifyCmd.PersistentFlags().StringVarP(&PubKey, "pubkey", "p", "",
		"public key to check with if custom protocol omits it from the output")

	verifyCmd.PersistentFlags().StringVarP(&Custom, "custom", "k", "",
		"custom additional namespace")

	rootCmd.AddCommand(verifyCmd)
}
