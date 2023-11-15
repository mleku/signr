package signr

import (
	"encoding/hex"
	"fmt"

	"github.com/mleku/ec/schnorr"
	secp "github.com/mleku/ec/secp"
	"github.com/mleku/signr/pkg/nostr"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show details of a nostr key",
	Long: `prints out the hex secret and public key and npub/nsec for use elsewhere.
`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			s.Fatal("key name is required")
		}
		var keyToShow string
		var err error
		var keySlice []string
		keySlice, err = s.GetKeyPairNames()
		if err != nil {
			s.Fatal("error: %s\n", err)
		}
		var found bool
		for _, k := range keySlice {
			if k == args[0] {
				found, keyToShow = true, k
			}
		}
		if !found {
			s.Fatal("'%s' key not found", args[1])
		}
		var sec *secp.SecretKey
		sec, err = s.GetKey(keyToShow, Pass)
		if err != nil {
			s.Fatal("error: %s\n", err)
		}
		pub := sec.PubKey()
		var hexSec, hexPub, npub, nsec string
		secBytes := sec.Serialize()
		hexSec = hex.EncodeToString(secBytes)
		hexPub = hex.EncodeToString(schnorr.SerializePubKey(pub))
		nsec, _ = nostr.SecretKeyToNsec(sec)
		npub, _ = nostr.PublicKeyToNpub(pub)
		fmt.Printf("secret key:     %s\npublic key:     %s\n", nsec, npub)
		fmt.Printf("hex secret key: %s\nhex public key: %s\n", hexSec, hexPub)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
