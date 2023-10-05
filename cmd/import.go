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
	"strings"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import <secret key> <name>",
	Short: "Import a secret key",
	Long: `Import a secret key provided in hexadecimal, nostr nsec or bip39 word key formats. 

Each format has a distinctive prefix that enables them to be automatically recognised.

Commonly hexadecimal format nostr keys are not prefixed correctly with '0x' so be aware this is required or it will be rejected.

Word keys will be tried in english first and then all the other languages if the key does not start with 'nsec' or '0x'

When importing a word key, there must be 25 words, the last word is the name used to refer to the key.
`,
	Run: func(cmd *cobra.Command, args []string) {
		argLen := len(args)
		if argLen == 1 {
			_, _ = fmt.Fprintln(os.Stderr,
				"a key name is required after the secret key")
			os.Exit(1)
		}
		if argLen == 24 {
			_, _ = fmt.Fprintln(os.Stderr,
				"a key name is required after the 24 words of the bip39 word key")
			os.Exit(1)
		}
		var sec *secp.SecretKey
		var err error
		switch {
		case argLen == 25:
			fmt.Println(len(args), args[:24])
			secBytes := bip39.MnemonicToSeed(strings.Join(args[:24], " "), "")
			fmt.Printf("hex:\n\tsecret: 0x%s\n\tpublic: 0x%s\n",
				hex.EncodeToString(secBytes),
			)
			sec = secp.SecKeyFromBytes(secBytes)
		case strings.HasPrefix(args[0], nostr.SecHRP):
			sec, err = nostr.DecodeSecretKey(args[0])
		default:
			var secBytes []byte
			secBytes, err = hex.DecodeString(args[0])
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr,
					"key should be in hex or is somehow mangled, cannot decode: "+err.Error())
				os.Exit(1)
			}
			sec = secp.SecKeyFromBytes(secBytes)
		}
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if sec == nil {
			_, _ = fmt.Fprintln(os.Stderr,
				"input did not match any known formats")
			os.Exit(1)
		}
		pub := sec.PubKey()
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

func init() {
	rootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
