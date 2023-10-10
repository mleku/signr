package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/minio/sha256-simd"
	"github.com/mleku/ec/schnorr"
	secp "github.com/mleku/ec/secp"
	"github.com/mleku/signr/pkg/nostr"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/argon2"
	"io"
	"os"
	"strings"
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign [flags] <file> [key name]",
	Short: "Generate a signature on a file",
	Long: `Generates a signature for the provided file, using either the default key, or if specified, another key.

Use the filename '-' to indicate that the data is being piped in with < or |.

Currently there isn't actually any flags, but there can be in the future.
`,
	Run: func(cmd *cobra.Command, args []string) {

		signingKey := defaultKey
		filename := args[0]

		var err error

		switch {
		case len(args) < 1:

			PrintErr("ERROR: at minimum a file to be signed needs to " +
				"be specified\n\n")
			os.Exit(1)

		case len(args) > 1:

			var keySlice []string
			keySlice, err = GetKeyPairNames()
			if err != nil {

				PrintErr("ERROR: '%s'\n", err)
				os.Exit(1)
			}

			var found bool
			for _, k := range keySlice {

				if k == args[1] {
					found, signingKey = true, k
				}
			}

			if !found {

				PrintErr("'%s' key not found\n", args[1])

				listkeysCmd.Run(cmd, nil)

				return
			}
		}

		// for now the first 4 are always the same
		signingStrings := []string{
			"signr", "0", "SHA256", "SCHNORR",
		}

		// add the signature nonce
		nonce := make([]byte, 8)

		_, err = rand.Read(nonce)
		if err != nil {
			PrintErr("ERROR: '%s'\n\n", err)
			os.Exit(1)
		}

		signingStrings = append(signingStrings, hex.EncodeToString(nonce))

		// add the public key
		var pkb []byte
		pkb, err = ReadFile(signingKey + "." + pubExt)
		if err != nil {
			PrintErr("ERROR: '%s'\n", err)
			os.Exit(1)
		}

		signingStrings = append(signingStrings, strings.TrimSpace(string(pkb)))

		var f io.ReadCloser

		switch {
		case filename == "-":

			f = os.Stdin
			// read from stdin
		default:

			f, err = os.Open(filename)
			if err != nil {
				PrintErr(
					"ERROR: unable to open file: '%s'\n\n", err)
				os.Exit(1)
			}

			defer func(f io.ReadCloser) {
				err := f.Close()
				if err != nil {
					PrintErr("ERROR: '%s'\n", err)
					os.Exit(1)
				}
			}(f)

		}

		h := sha256.New()

		if _, err := io.Copy(h, f); err != nil {
			PrintErr(
				"ERROR: unable to read file to generate hash: '%s'\n\n", err)
			os.Exit(1)
		}

		sum := h.Sum(nil)

		signingStrings = append(signingStrings, hex.EncodeToString(sum))

		message := strings.Join(signingStrings, "_")

		messageHash := sha256.Sum256([]byte(message))

		var key *secp.SecretKey
		key, err = GetKey(signingKey, cmd)
		if err != nil {
			PrintErr("ERROR: '%s'\n", err)
			os.Exit(1)
		}

		var sig *schnorr.Signature
		sig, err = schnorr.Sign(key, messageHash[:])

		var str string
		str, err = FormatSig(signingStrings, sig)
		if err != nil {
			PrintErr("ERROR: '%s'\n", err)
			os.Exit(1)
		}

		fmt.Println(str)
	},
}

func FormatSig(signingStrings []string, sig *schnorr.Signature) (str string,
	err error) {

	prefix := signingStrings[:len(signingStrings)-1]

	var sigStr string
	sigStr, err = nostr.EncodeSignature(sig)
	return strings.Join(
		append(prefix, sigStr), "_"), err
}

const unlockPrompt = "type password to unlock encrypted secret key: "

func GetKey(name string, cmd *cobra.Command) (key *secp.SecretKey, err error) {

	var keyBytes []byte
	keyBytes, err = ReadFile(name)
	if err != nil {
		err = errors.Wrap(err, "error getting key bytes:")
	}

	var encrypted bool

	for i, sb := range keyBytes {

		if sb == ' ' {
			if len(keyBytes) >= 64 {

				if keyBytes[i+1] == '*' {

					keyBytes = keyBytes[:64]
					encrypted = true
					break
				}
			}
		} else if sb == '\n' {

			keyBytes = keyBytes[:64]
			break
		}
	}
	_, err = hex.Decode(keyBytes, keyBytes)
	if err != nil {
		PrintErr("ERROR: '%v", err)
		return
	}

	originalSecret := keyBytes[:32]

	secret := make([]byte, 32)
	copy(secret, originalSecret)

	passFlag := cmd.PersistentFlags().Lookup("pass")

	if passFlag.Value.String() != "" {

		pass := passFlag.Value.String()
		actualKey := argon2.Key([]byte(pass),
			[]byte("signr"), 3, 1024*1024, 4, 32)

		secret = xor(secret, actualKey)

		sec := secp.SecKeyFromBytes(secret)
		pub := sec.PubKey()

		pubBytes := schnorr.SerializePubKey(pub)
		npub, _ := nostr.PublicKeyToString(pub)

		// check the decrypted secret yields the stored pubkey
		pubBytes, err = ReadFile(name + "." + pubExt)
		npubReal := strings.TrimSpace(string(pubBytes))

		if npub != npubReal {

			PrintErr("password failed to unlock key\n", err)
			return

		} else {

			key = sec
			return
		}
	}

	if encrypted {

		var tryCount int
		for tryCount < 3 {

			pass, err := PasswordEntry(unlockPrompt)
			if err != nil {
				PrintErr(
					"error in password input: '%s'\n", err)
				continue
			}

			actualKey := argon2.
				Key(pass, []byte("signr"), 3, 1024*1024, 4, 32)

			secret = xor(secret, actualKey)

			sec := secp.SecKeyFromBytes(secret)
			pub := sec.PubKey()

			pubBytes := schnorr.SerializePubKey(pub)
			npub, _ := nostr.PublicKeyToString(pub)

			// check the decrypted secret generates the stored pubkey
			pubBytes, err = ReadFile(name + "." + pubExt)
			npubReal := strings.TrimSpace(string(pubBytes))

			if npub != npubReal {
				PrintErr("password failed to unlock key, try again\n", err)
				tryCount++
				continue

			} else {

				key = sec
				break
			}
		}

	} else {

		key = secp.SecKeyFromBytes(originalSecret)
	}
	return
}

func init() {
	rootCmd.AddCommand(signCmd)
	signCmd.PersistentFlags().String("pass", "", "password to unlock the key")
}
