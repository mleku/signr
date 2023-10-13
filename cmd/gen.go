package cmd

import (
	"encoding/hex"
	"github.com/mleku/ec/schnorr"
	"github.com/mleku/signr/pkg/nostr"
	"github.com/mleku/signr/pkg/signr"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen <name>",
	Short: "Generate a new nostr key",
	Long:  `Generates a new key saves it into the application data directory as part of the keychain.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {

			signr.Fatal("key name is required")
		}

		sec, pub, err := signr.GenKeyPair()
		if err != nil {

			signr.Fatal("error generating key: '%s'", err)
		}

		secBytes := sec.Serialize()

		npub, _ := nostr.PublicKeyToString(pub)

		if cfg.Verbose {

			pubBytes := schnorr.SerializePubKey(pub)

			signr.PrintErr(
				"generated key pair:\n"+
					"\nhex:\n"+
					"\tsecret: %s\n"+
					"\tpublic: %s\n\n",
				hex.EncodeToString(secBytes),
				hex.EncodeToString(pubBytes),
			)

			nsec, _ := nostr.SecretKeyToString(sec)
			signr.PrintErr(
				"nostr:\n"+
					"\tsecret: %s\n"+
					"\tpublic: %s\n\n", nsec,
				npub)
		}

		if err = signr.Save(cfg, args[0], secBytes, npub); err != nil {

			signr.PrintErr("error saving keys: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}
