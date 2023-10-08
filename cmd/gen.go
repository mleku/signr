package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/mleku/ec/schnorr"
	secp "github.com/mleku/ec/secp"
	"github.com/mleku/signr/pkg/nostr"
	"github.com/spf13/cobra"
	"os"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen <name>",
	Short: "Generate a new nostr key",
	Long: `Generates a new key and prints it out in all available formats as
pairs of public and private keys.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			_, _ = fmt.Fprintln(os.Stderr, "key name is required")
			os.Exit(1)
		}

		sec, pub, err := GenKeyPair()
		if err != nil {
			printErr(
				"error generating key: '%s'", err)
			os.Exit(1)
		}
		secBytes := sec.Serialize()
		pubBytes := schnorr.SerializePubKey(pub)
		if verbose {
			printErr(
				"generated key pair:\n\nhex:\n\tsecret: %s\n\tpublic: %s\n\n",
				hex.EncodeToString(secBytes),
				hex.EncodeToString(pubBytes),
			)
		}
		nsec, _ := nostr.SecretKeyToString(sec)
		npub, _ := nostr.PublicKeyToString(pub)
		if verbose {
			printErr(
				"nostr:\n\tsecret: %s\n\tpublic: %s\n\n", nsec,
				npub)
		}

		if err = Save(args[0], secBytes, npub); err != nil {
			printErr(
				"error saving keys: %v", err)
		}
	},
}

func GenKeyPair() (sec *secp.SecretKey, pub *secp.PublicKey, err error) {

	sec, err = secp.GenerateSecretKey()
	if err != nil {
		printErr("error generating key: '%s'", err)
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
