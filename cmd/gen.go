package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/islishude/bip39"
	"github.com/mleku/ec/schnorr"
	secp "github.com/mleku/ec/secp"
	"github.com/mleku/signr/pkg/bip39langs"
	"github.com/mleku/signr/pkg/nostr"
	"github.com/spf13/cobra"
	"os"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate a new nostr key",
	Long: `Generates a new key and prints it out in all available formats as
pairs of public and private keys.`,
	Run: func(cmd *cobra.Command, args []string) {
		sec, pub, err := GenKeyPair()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr,
				"error generating key: '%s'", err)
			os.Exit(1)
		}
		secBytes := sec.Serialize()
		pubBytes := schnorr.SerializePubKey(pub)
		fmt.Printf("hex:\n\tsecret: %s\n\tpublic: %s\n",
			hex.EncodeToString(secBytes),
			hex.EncodeToString(pubBytes),
		)

		nsec, _ := nostr.SecretKeyToString(sec)
		npub, _ := nostr.PublicKeyToString(pub)
		fmt.Printf("nostr:\n\tsecret: %s\n\tpublic: %s\n", nsec, npub)
		var mnem string
		lang := rootCmd.PersistentFlags().Lookup("lang").Value.String()
		mnem, _ = bip39.NewMnemonicByEntropy(secBytes, bip39langs.Map[lang])
		fmt.Printf("bip39:\n\t%s\n", mnem)
	},
}

func GenKeyPair() (sec *secp.SecretKey, pub *secp.PublicKey, err error) {
	sec, err = secp.GenerateSecretKey()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error generating key: '%s'", err)
		return
	}
	pub = sec.PubKey()

	return
}

func init() {
	rootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
