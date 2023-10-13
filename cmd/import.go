package cmd

import (
	"encoding/hex"
	"github.com/mleku/ec/schnorr"
	secp "github.com/mleku/ec/secp"
	"github.com/mleku/signr/pkg/nostr"
	"github.com/mleku/signr/pkg/signr"
	"github.com/spf13/cobra"
	"strings"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import <secret key> <name>",
	Short: "Import a secret key",
	Long: `Import a secret key provided in hexadecimal and nostr nsec formats.
`,
	Run: func(cmd *cobra.Command, args []string) {

		argLen := len(args)
		if argLen == 1 {

			signr.Fatal("a key name is required after the secret key")
		}

		var sec *secp.SecretKey
		var err error
		if strings.HasPrefix(args[0], nostr.SecHRP) {

			if sec, err = nostr.DecodeSecretKey(args[0]); err != nil {

				signr.Fatal("ERROR: while decoding key: '%v'\n", err)
			}

		} else {

			var secBytes []byte
			if secBytes, err = hex.DecodeString(args[0]); err != nil {

				signr.Fatal(
					"key is mangled, '%s', cannot decode: '%v'\n", args[0], err)
			}

			sec = secp.SecKeyFromBytes(secBytes)
		}

		if sec == nil {

			signr.Fatal("input did not match any known formats")
		}

		pub := sec.PubKey()
		secBytes := sec.Serialize()

		npub, _ := nostr.PublicKeyToString(pub)

		if cfg.Verbose {

			pubBytes := schnorr.SerializePubKey(pub)

			signr.PrintErr("hex:\n"+
				"\tsecret: %s\n"+
				"\tpublic: %s\n",
				hex.EncodeToString(secBytes),
				hex.EncodeToString(pubBytes),
			)

			nsec, _ := nostr.SecretKeyToString(sec)

			signr.PrintErr("nostr:\n"+
				"\tsecret: %s\n"+
				"\tpublic: %s\n\n",
				nsec, npub)
		}

		if err = signr.Save(cfg, args[1], secBytes, npub); err != nil {

			signr.Fatal("error saving keys: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
