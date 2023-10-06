package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/mleku/ec/schnorr"
	secp "github.com/mleku/ec/secp"
	"github.com/mleku/signr/pkg/nostr"
	"github.com/spf13/cobra"
	"os"
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
			_, _ = fmt.Fprintln(os.Stderr,
				"a key name is required after the secret key")
			os.Exit(1)
		}

		var sec *secp.SecretKey
		var err error
		if strings.HasPrefix(args[0], nostr.SecHRP) {
			sec, err = nostr.DecodeSecretKey(args[0])
		} else {
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
		if verbose {
			fmt.Printf("hex:\n\tsecret: %s\n\tpublic: %s\n",
				hex.EncodeToString(secBytes),
				hex.EncodeToString(pubBytes),
			)
		}
		nsec, _ := nostr.SecretKeyToString(sec)
		npub, _ := nostr.PublicKeyToString(pub)
		if verbose {
			fmt.Printf("nostr:\n\tsecret: %s\n\tpublic: %s\n\n",
				nsec, npub)
		}
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
