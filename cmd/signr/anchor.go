package signr

import (
	"fmt"

	"github.com/spf13/cobra"
)

// anchorCmd represents the import command
var anchorCmd = &cobra.Command{
	Use:   "anchor [--custom/-k NAME-SPACE] <filename/merkle root/hash> [keyname]",
	Short: "generate required info to anchor a hash and signature on chain",
	Long: `from the private key, and a provided hash value, a signature is generated and the necessary elements for an inscription using taproot on Bitcoin are provided in the following format:

WIF NPUB MERKLE NSIG

The WIF is a standard Bitcoin WIF private/secret key in Base58check format, WARNING: unencrypted - be careful with it.

For your convenience, the NPUB, MERKLE and NSIG are presented in hexadecimal so they can be easily concatenated and converted into a binary string for signing to the relevant Taproot address related to the secret key you named.

From the NPUB, MERKLE and NSIG you can use the verifyanchor function to validate that the NSIG is a valid signature on the MERKLE that matches the NPUB.
` + customInfo,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			s.Err("the merkle root/hash is a required parameter")
		}
		if wif, npub, merkle, nsig, err := s.
			GetAnchor(args, Pass, Custom); err != nil {

			s.Err("ERROR: while generating anchor: %s\n", err)
		} else {
			fmt.Printf("%s %s %s %s\n", wif, npub, merkle, nsig)
		}
	},
}

func init() {
	anchorCmd.PersistentFlags().StringVarP(&Custom, "custom", "k", "",
		"custom namespace separator")
	anchorCmd.PersistentFlags().StringVarP(&Pass, "pass", "p", "",
		"password to unlock the key - for better security, use the "+
			"environment variable")
	rootCmd.AddCommand(anchorCmd)
}
