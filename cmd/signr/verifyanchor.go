package signr

import (
	"github.com/spf13/cobra"
)

// anchorCmd represents the import command
var verifyanchorCmd = &cobra.Command{
	Use:   "verifyanchor  [--custom/-k NAME-SPACE] <NPUB> <MERKLE/hash> <NSIG>",
	Short: "validate an on-chain anchor incription",
	Long: `verify that an anchor inscription is valid.

The combination of NPUB MERKLE and NSIG, as produced by the anchor command of signr can be verified to be correct using verifyanchor.

The correct signing material is generated from the NPUB and MERKLE hash values, and then from the hash of this material, one can verify the NPUB and MERKLE match the NSIG.

For security against cross protocol attacks, ensure that you provide the correct --custom namespace string used to generate the anchor. By chance it can happen that distinct combinations of private key and hash can derive the same signature from separate protocols. The custom field helps ensure that such collisions do not compromise the security of users of your protocol.
`+customInfo,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	verifyanchorCmd.PersistentFlags().StringVarP(&Custom, "custom", "k", "",
		"custom namespace separator")
	rootCmd.AddCommand(verifyanchorCmd)
}
