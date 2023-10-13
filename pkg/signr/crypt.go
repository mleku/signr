package signr

import (
	"encoding/hex"
	"github.com/mleku/ec/schnorr"
	"github.com/mleku/ec/secp"
	"github.com/mleku/signr/pkg/nostr"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/argon2"
	"strings"
)

// XOR two same length slices of bytes.
func XOR(dest, src []byte) []byte {
	if len(src) != len(dest) {
		PrintErr("key and secret must be the same length")
	}
	for i := range dest {
		dest[i] = dest[i] ^ src[i]
	}
	return dest
}

const UnlockPrompt = "type password to unlock encrypted secret key: "

func GetKey(dataDir, name, pass string,
	cmd *cobra.Command) (key *secp256k1.SecretKey,
	err error) {

	var keyBytes []byte
	keyBytes, err = ReadFile(dataDir, name)
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

	if pass != "" {

		actualKey := argon2.Key([]byte(pass),
			[]byte("signr"), 3, 1024*1024, 4, 32)

		secret = XOR(secret, actualKey)

		sec := secp256k1.SecKeyFromBytes(secret)
		pub := sec.PubKey()

		pubBytes := schnorr.SerializePubKey(pub)
		npub, _ := nostr.PublicKeyToString(pub)

		// check the decrypted secret yields the stored pubkey
		pubBytes, err = ReadFile(dataDir, name+"."+PubExt)
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

			pass, err := PasswordEntry(UnlockPrompt, 0)
			if err != nil {
				PrintErr(
					"error in password input: '%s'\n", err)
				continue
			}

			actualKey := argon2.Key(pass, []byte("signr"), 3, 1024*1024, 4, 32)

			secret = XOR(secret, actualKey)

			sec := secp256k1.SecKeyFromBytes(secret)
			pub := sec.PubKey()

			pubBytes := schnorr.SerializePubKey(pub)
			npub, _ := nostr.PublicKeyToString(pub)

			// check the decrypted secret generates the stored pubkey
			pubBytes, err = ReadFile(dataDir, name+"."+PubExt)
			npubReal := strings.TrimSpace(string(pubBytes))

			if npub != npubReal {
				PrintErr("password failed to unlock key, try again\n",
					err)
				tryCount++
				continue

			} else {

				key = sec
				break
			}
		}

	} else {

		key = secp256k1.SecKeyFromBytes(originalSecret)
	}
	return
}

// GenKeyPair creates a fresh new key pair using the entropy source used by
// crypto/rand (ie, /dev/random on posix systems).
func GenKeyPair() (sec *secp256k1.SecretKey, pub *secp256k1.PublicKey,
	err error) {

	sec, err = secp256k1.GenerateSecretKey()
	if err != nil {

		PrintErr("error generating key: '%s'", err)
		return
	}

	pub = sec.PubKey()

	return
}
