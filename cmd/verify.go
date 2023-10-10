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

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify <file> <signature file/signature>",
	Short: "check that a file matches a signature",
	Long: `checks that a signature is valid for a given file.

the public key is embedded in the signature should also match the key expected to be used on the signature.

use the filename '-' to indicate the file is being piped in via stdin.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			printErr("ERROR: at minimum a file and a keyfile name must be specified\n\n")
			os.Exit(1)
		}
		filename := args[0]
		sigOrSigFile := args[1]
		var signingStrings []string
		var err error
		if strings.HasPrefix(sigOrSigFile, "signr_0_SHA256_SCHNORR_") {
			signingStrings = strings.Split(sigOrSigFile, "_")
		} else {
			var data []byte
			data, err = os.ReadFile(sigOrSigFile)
			if err != nil {
				printErr(
					"ERROR: reading file '%s': %v\n", sigOrSigFile, err)
				return
			}
			signingStrings = strings.Split(string(data), "_")
		}
		var f io.ReadCloser
		switch {
		case filename == "-":
			f = os.Stdin
			// read from stdin
		default:
			f, err = os.Open(filename)
			if err != nil {
				printErr(
					"ERROR: unable to open file: '%s'\n\n", err)
				os.Exit(1)
			}
			defer func(f io.ReadCloser) {
				err := f.Close()
				if err != nil {
					printErr("ERROR: closing file '%s'\n", err)
					os.Exit(1)
				}
			}(f)

		}
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			printErr(
				"ERROR: unable to read file to generate hash: '%s'\n\n", err)
			os.Exit(1)
		}
		sum := h.Sum(nil)
		signature := strings.TrimSpace(signingStrings[len(signingStrings)-1])
		pubkey := signingStrings[len(signingStrings)-2]
		signingStrings = signingStrings[:len(signingStrings)-1]
		signingStrings = append(signingStrings, hex.EncodeToString(sum))
		message := strings.Join(signingStrings, "_")
		messageHash := sha256.Sum256([]byte(message))

		var sig *schnorr.Signature
		sig, err = nostr.DecodeSignature(signature)
		if err != nil {
			printErr("ERROR: decoding signature '%s'\n", err)
			os.Exit(1)
		}
		var pk *secp.PublicKey
		pk, err = nostr.DecodePublicKey(pubkey)
		if sig.Verify(messageHash[:], pk) {
			fmt.Println("VALID")
		} else {
			fmt.Println("INVALID")
		}
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// verifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
