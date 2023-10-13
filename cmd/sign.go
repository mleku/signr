package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/minio/sha256-simd"
	"github.com/mleku/ec/schnorr"
	secp "github.com/mleku/ec/secp"
	"github.com/mleku/signr/pkg/nostr"
	"github.com/mleku/signr/pkg/signr"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

var (
	PrintAsHex   bool
	PrintOnlySig bool
	Pass         string
	Custom       string
	PubKey       string
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign [flags] <file/hash> [key name]",
	Short: "Generate a signature on a file or hash",
	Long: `Generates a signature for the provided file, using either the default key, or if specified, another key.

Use the filename '-' to indicate that the data is being piped in with < or |.

If filename is in fact a 64 character hexadecimal value, it will be signed on without the nonce and public key, and yield only a bech32 signature, or hexadecimal with the --hex flag.
`,
	Run: func(cmd *cobra.Command, args []string) {

		signingKey := cfg.DefaultKey
		filename := args[0]

		var err error

		switch {
		case len(args) < 1:

			signr.PrintErr("ERROR: at minimum a file to be signed needs to " +
				"be specified\n\n")
			os.Exit(1)

		case len(args) > 1:

			var keySlice []string
			keySlice, err = signr.GetKeyPairNames(cfg)
			if err != nil {

				signr.PrintErr("ERROR: '%s'\n", err)
				os.Exit(1)
			}

			var found bool
			for _, k := range keySlice {

				if k == args[1] {
					found, signingKey = true, k
				}
			}

			if !found {

				signr.PrintErr("'%s' key not found\n", args[1])

				listkeysCmd.Run(cmd, nil)

				return
			}
		}

		signingStrings := signr.GetDefaultSigningStrings()

		signingStrings = signr.AddCustom(signingStrings, Custom)

		var skipRandomness bool

		// if the command line contains a raw hash we assume that a simple
		// signature on this is intended. it will still use the namespacing, the
		// pubkey and any custom string, but not the nonce. it is assumed that
		// the protocol generating the hash has accounted for sufficient
		// entropy.
		var sum []byte
		if len(filename) == 64 {

			sum, err = hex.DecodeString(filename)
			if err == nil {
				skipRandomness = true
			}

		}

		if !skipRandomness {
			// add the signature nonce
			nonce := make([]byte, 8)

			_, err = rand.Read(nonce)
			if err != nil {
				signr.PrintErr("ERROR: '%s'\n\n", err)
				os.Exit(1)
			}

			signingStrings = append(signingStrings, hex.EncodeToString(nonce))

			var f io.ReadCloser

			switch {
			case filename == "-":

				// read from stdin
				f = os.Stdin

			default:

				f, err = os.Open(filename)
				if err != nil {
					signr.PrintErr(
						"ERROR: unable to open file: '%s'\n\n", err)
					os.Exit(1)
				}

				defer func(f io.ReadCloser) {
					err := f.Close()
					if err != nil {
						signr.PrintErr("ERROR: '%s'\n", err)
						os.Exit(1)
					}
				}(f)

			}
			h := sha256.New()

			if _, err := io.Copy(h, f); err != nil {
				signr.PrintErr(
					"ERROR: unable to read file to generate hash: '%s'\n\n",
					err)
				os.Exit(1)
			}

			sum = h.Sum(nil)

		}

		// add the public key. This must always be present as it isolates
		// the namespace of even intra-protocol signing.
		var pkb []byte
		pkb, err = signr.ReadFile(cfg.DataDir, signingKey + "." + signr.PubExt)
		if err != nil {

			signr.Fatal("ERROR: '%s'\n", err)
		}

		// the keychain stores secrets as hex but the pubkeys in nostr npub.
		signingStrings = append(signingStrings, strings.TrimSpace(string(pkb)))

		// append the checksum.
		signingStrings = append(signingStrings, hex.EncodeToString(sum))

		// construct the signing material.
		message := strings.Join(signingStrings, "_")

		if cfg.Verbose {
			signr.PrintErr("signing on message: %s\n", message)
		}

		messageHash := sha256.Sum256([]byte(message))

		var key *secp.SecretKey
		key, err = signr.GetKey(cfg.DataDir, signingKey, Pass, cmd)
		if err != nil {
			signr.Fatal("ERROR: '%s'\n", err)
		}

		var sig *schnorr.Signature
		sig, err = schnorr.Sign(key, messageHash[:])
		var sigStr string

		if skipRandomness {

			if PrintAsHex {

				sigStr = hex.EncodeToString(sig.Serialize())

			} else {

				sigStr, err = nostr.EncodeSignature(sig)
				if err != nil {

					signr.PrintErr("ERROR: while formatting signature: '%s'\n",
						err)
					return
				}
			}
			fmt.Println(sigStr)
			return

		} else if PrintOnlySig {

			if PrintAsHex {

				sigStr = hex.EncodeToString(sig.Serialize())

			} else {

				sigStr, err = nostr.EncodeSignature(sig)
				if err != nil {

					signr.PrintErr("ERROR: while formatting signature: '%s'\n",
						err)
					return
				}
			}

			fmt.Println(sigStr)
			return
		}
		sigStr, err = signr.FormatSig(signingStrings, sig)
		if err != nil {

			signr.Fatal("ERROR: '%s'\n", err)
		}

		fmt.Println(sigStr)
	},
}

func init() {
	signCmd.PersistentFlags().StringVar(&Pass, "pass", "",
		"password to unlock the key")
	signCmd.PersistentFlags().StringVar(&Custom, "custom", "",
		"custom additional namespace")
	signCmd.PersistentFlags().BoolVar(&PrintAsHex, "hex", false,
		"print signature in hex")
	signCmd.PersistentFlags().BoolVar(&PrintOnlySig, "sigonly", false,
		"print only signature")
	rootCmd.AddCommand(signCmd)
}
