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

var genCmd = &cobra.Command{
	Use:   "gen <name>",
	Short: "Generate a new nostr key",
	Long: `Generates a new key saves it into the application data directory as part of the keychain.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			_, _ = fmt.Fprintln(os.Stderr, "key name is required")
			os.Exit(1)
		}

		sec, pub, err := GenKeyPair()
		if err != nil {
			PrintErr(
				"error generating key: '%s'", err)
			os.Exit(1)
		}
		secBytes := sec.Serialize()
		pubBytes := schnorr.SerializePubKey(pub)
		if verbose {
			PrintErr(
				"generated key pair:\n\nhex:\n\tsecret: %s\n\tpublic: %s\n\n",
				hex.EncodeToString(secBytes),
				hex.EncodeToString(pubBytes),
			)
		}
		nsec, _ := nostr.SecretKeyToString(sec)
		npub, _ := nostr.PublicKeyToString(pub)
		if verbose {
			PrintErr(
				"nostr:\n\tsecret: %s\n\tpublic: %s\n\n", nsec,
				npub)
		}

		if err = Save(args[0], secBytes, npub); err != nil {
			PrintErr(
				"error saving keys: %v", err)
		}
	},
}

func GenKeyPair() (sec *secp.SecretKey, pub *secp.PublicKey, err error) {

	sec, err = secp.GenerateSecretKey()
	if err != nil {
		PrintErr("error generating key: '%s'", err)
		return
	}

	pub = sec.PubKey()

	return
}

func init() {
	rootCmd.AddCommand(genCmd)
}
