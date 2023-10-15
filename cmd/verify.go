package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/minio/sha256-simd"
	"github.com/mleku/ec/schnorr"
	secp "github.com/mleku/ec/secp"
	"github.com/mleku/signr/pkg/nostr"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

var PubKey string

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify <file> <signature file/signature>",
	Short: "check that a file matches a signature",
	Long: `checks that a signature is valid for a given file.

the public key is embedded in the signature should also match the key expected to be used on the signature.

use the filename '-' to indicate the file is being piped in via stdin.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			cfg.Err("ERROR: at minimum a file and a keyfile name must be specified\n\n")
			os.Exit(1)
		}

		filename := args[0]
		sigOrSigFile := args[1]

		var signingStrings []string
		var err error

		// if the parameter appears to be a signature, use it directly
		if strings.HasPrefix(sigOrSigFile, "signr_0_SHA256_SCHNORR_") {

			signingStrings = strings.Split(sigOrSigFile, "_")

		} else {

			var data []byte

			data, err = os.ReadFile(sigOrSigFile)
			if err != nil {
				cfg.Err(
					"ERROR: reading file '%s': %v\n", sigOrSigFile, err)
				return
			}

			signingStrings = strings.Split(string(data), "_")
		}

		var f io.ReadCloser
		switch {
		case filename == "-":

			// read from stdin
			f = os.Stdin

		default:

			// read from the named file
			f, err = os.Open(filename)
			if err != nil {
				cfg.Fatal(
					"ERROR: unable to open file: '%s'\n\n", err)
			}

			defer func(f io.ReadCloser) {
				err := f.Close()
				if err != nil {
					cfg.Fatal("ERROR: closing file '%s'\n", err)
				}
			}(f)

		}

		h := sha256.New()

		// feed the file data through the hasher
		if _, err := io.Copy(h, f); err != nil {
			cfg.Fatal(
				"ERROR: unable to read file to generate hash: '%s'\n\n", err)
		}
		sum := h.Sum(nil)

		// clean up the signature
		signature := strings.TrimSpace(signingStrings[len(signingStrings)-1])

		// get the public key
		pubkey := signingStrings[len(signingStrings)-2]

		// trim off the signature part from the signature
		signingStrings = signingStrings[:len(signingStrings)-1]

		// add the hash of the file in place
		signingStrings = append(signingStrings, hex.EncodeToString(sum))

		// generate the signing material
		message := strings.Join(signingStrings, "_")

		// hash the signing material
		messageHash := sha256.Sum256([]byte(message))

		// decode the signature
		var sig *schnorr.Signature
		sig, err = nostr.DecodeSignature(signature)
		if err != nil {

			cfg.Fatal("ERROR: decoding signature '%s'\n", err)
		}

		// decode the public key
		var pk *secp.PublicKey
		pk, err = nostr.DecodePublicKey(pubkey)

		// verify the hash and the signature match the public key
		if sig.Verify(messageHash[:], pk) {

			fmt.Println("VALID")

		} else {

			fmt.Println("INVALID")
		}
	},
}

func init() {
	verifyCmd.PersistentFlags().StringVar(&PubKey, "pubkey", "",
		"public key to check with if custom protocol omits it from the output")
	verifyCmd.PersistentFlags().StringVarP(&Custom, "custom", "k", "",
		"custom additional namespace")
	rootCmd.AddCommand(verifyCmd)
}
