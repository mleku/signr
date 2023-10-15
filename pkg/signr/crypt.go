package signr

import (
	"encoding/hex"
	"fmt"
	"github.com/mleku/ec/schnorr"
	"github.com/mleku/ec/secp"
	"github.com/mleku/signr/pkg/nostr"
	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
	"strings"
)

const UnlockPrompt = "type password to unlock encrypted secret key"

// GetKey scans the keychain for a named key, with optional password string to
// decrypt the key in the file in the keychain.
func (s *Signr) GetKey(name, passStr string) (key *secp256k1.SecretKey,
	err error) {

	var keyBytes []byte
	keyBytes, err = s.ReadFile(name)
	if err != nil {
		err = errors.Wrap(err, "error getting key bytes:")
		return
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
		s.PrintErr("ERROR: decoding key hex: '%v", err)
		return
	}

	originalSecret := keyBytes[:32]

	secret := make([]byte, 32)

	if passStr != "" {
		copy(secret, originalSecret)
		if key, err = s.DeriveAndCheckKey(name, secret,
			[]byte(passStr)); err != nil {

			s.PrintErr("password failed to unlock key: %s\n", err)

			return
		}

	}

	var pass []byte

	if encrypted {

		var tryCount int
		retryStr := ""
		for tryCount < 3 {

			copy(secret, originalSecret)
			if tryCount > 0 {
				retryStr = fmt.Sprintf(" (attempt %d of %d)", tryCount, 3)
			}
			unlockPrompt := fmt.Sprintf("%s%s:", UnlockPrompt, retryStr)
			pass, err = s.PasswordEntry(unlockPrompt, 0)
			if err != nil {
				s.PrintErr(
					"error in password input: '%s'\n", err)
				continue
			}

			if key, err = s.DeriveAndCheckKey(name, secret,
				pass); err != nil {
				tryCount++
				continue
			} else {
				break
			}

		}

	} else {

		key = secp256k1.SecKeyFromBytes(originalSecret)
	}
	return
}

func (s *Signr) DeriveAndCheckKey(name string,
	secret, pass []byte) (sec *secp256k1.SecretKey, err error) {

	actualKey := ArgonKey(pass)

	secret = s.XOR(secret, actualKey)

	sec = secp256k1.SecKeyFromBytes(secret)
	pub := sec.PubKey()

	pubBytes := schnorr.SerializePubKey(pub)
	npub, _ := nostr.PublicKeyToString(pub)

	// check the decrypted secret generates the stored pubkey
	pubBytes, err = s.ReadFile(name + "." + PubExt)
	if err != nil {
		s.PrintErr("error reading pubkey: %s\n", err)
		return
	}
	npubReal := strings.TrimSpace(string(pubBytes))

	s.Log("secret decrypted: %v; decrypted->pub: %s, stored pub; %s\n",
		npub == npubReal, npub, npubReal)

	if npub != npubReal {
		err = fmt.Errorf("ERROR: %s, password failed to unlock key, try again\n",
			err)
	}

	return
}

// XOR two same length slices of bytes.
func (s *Signr) XOR(dest, src []byte) []byte {
	if len(src) != len(dest) {
		s.Err("key and secret must be the same length")
	}
	for i := range dest {
		dest[i] = dest[i] ^ src[i]
	}
	return dest
}

// ArgonKey hash grinds the input password string to derive the actual
// encryption key used on the secret key.
func ArgonKey(pass []byte) []byte {
	return argon2.Key(pass, []byte("signr"), 3, 1024*1024, 4, 32)
}

// GenKeyPair creates a fresh new key pair using the entropy source used by
// crypto/rand (ie, /dev/random on posix systems).
func (s *Signr) GenKeyPair() (sec *secp256k1.SecretKey,
	pub *secp256k1.PublicKey,
	err error) {

	sec, err = secp256k1.GenerateSecretKey()
	if err != nil {

		s.PrintErr("error generating key: '%s'", err)
		return
	}

	pub = sec.PubKey()

	return
}
