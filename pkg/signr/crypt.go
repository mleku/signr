package signr

import (
	"encoding/hex"
	"fmt"
	"strings"

	secp "github.com/mleku/ec/secp"
	"github.com/mleku/signr/pkg/nostr"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"
)

const UnlockPrompt = "type password to unlock encrypted secret key"

// GetKey scans the keychain for a named key, with optional password string to
// decrypt the key in the file in the keychain.
func (s *Signr) GetKey(name, passStr string) (key *secp.SecretKey,
	err error) {

	var keyBytes []byte
	if keyBytes, err = s.ReadFile(name); err != nil {
		err = errors.Wrap(err, "error getting key bytes:")
		return
	}
	var encrypted bool
	for i, sb := range keyBytes {
		if sb == ' ' && i == 64 && len(keyBytes) >= 66 {
			if keyBytes[i+1] == '*' {
				keyBytes = keyBytes[:64]
				encrypted = true
				break
			}
			// key must be before any linefeed, anything after is ignored.
		} else if sb == '\n' {
			keyBytes = keyBytes[:64]
			break
		}
	}
	if _, err = hex.Decode(keyBytes, keyBytes); err != nil {
		s.Err("ERROR: decoding key hex: '%v", err)
		return
	}
	originalSecret := keyBytes[:32]
	// sanitise memory whenever possible
	Zero(keyBytes[32:])
	secret := make([]byte, 32)
	if p := viper.GetString("pass"); p != "" {
		passStr = p
	}
	if passStr != "" {
		copy(secret, originalSecret)
		if key, err = s.
			DeriveAndCheckKey(name, secret, []byte(passStr)); err != nil {
			s.Err("password failed to unlock key: %s\n", err)
			return
		}
		return
	}
	var pass []byte
	if encrypted {
		var tryCount int
		retryStr := ""
		for tryCount < 3 {
			copy(secret, originalSecret)
			if tryCount > 0 {
				retryStr = fmt.Sprintf(" (attempt %d of %d)", tryCount+1, 3)
			}
			unlockPrompt := fmt.Sprintf("%s%s:", UnlockPrompt, retryStr)
			if pass, err = s.PasswordEntry(unlockPrompt, 0); err != nil {
				s.Err("error in password input: '%s'\n", err)
				continue
			}
			if key, err = s.
				DeriveAndCheckKey(name, secret, pass); err != nil {
				tryCount++
				continue
			} else {
				break
			}
		}
	} else {
		key = secp.SecKeyFromBytes(originalSecret)
	}
	return
}

func (s *Signr) DeriveAndCheckKey(name string,
	secret, pass []byte) (sec *secp.SecretKey, err error) {

	actualKey := ArgonKey(pass)
	secret = s.XOR(secret, actualKey)
	sec = secp.SecKeyFromBytes(secret)
	pub := sec.PubKey()
	npub, _ := nostr.PublicKeyToNpub(pub)
	// check the decrypted secret generates the stored pubkey
	var pubBytes []byte
	pubBytes, err = s.ReadFile(name + "." + PubExt)
	if err != nil {
		s.Err("error reading pubkey: %s\n", err)
		return
	}
	npubReal := strings.TrimSpace(string(pubBytes))
	s.Log("secret decrypted: %v; decrypted->pub: %s, stored pub; %s\n",
		npub == npubReal, npub, npubReal)
	if npub != npubReal {
		err = fmt.Errorf("ERROR: %s, password failed to unlock key, try again",
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
func (s *Signr) GenKeyPair() (sec *secp.SecretKey,
	pub *secp.PublicKey, err error) {

	sec, err = secp.GenerateSecretKey()
	if err != nil {
		err = fmt.Errorf("error generating key: %s", err)
		return
	}
	pub = sec.PubKey()
	return
}
