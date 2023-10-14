package signr

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/minio/sha256-simd"
	"github.com/mleku/ec/schnorr"
	secp "github.com/mleku/ec/secp"
	"github.com/mleku/signr/pkg/nostr"
	"io"
	"os"
	"strings"
)

func (cfg *Config) Sign(args []string, pass, custom string,
	asHex, sigOnly bool) (sigStr string, err error) {

	signingKey := cfg.DefaultKey
	filename := args[0]

	switch {
	case len(args) < 1:

		err = fmt.Errorf("ERROR: at minimum a file to be signed needs to " +
			"be specified\n\n")
		return

	case len(args) > 1:

		var keySlice []string
		keySlice, err = cfg.GetKeyPairNames()
		if err != nil {

			err = fmt.Errorf("ERROR: '%s'\n", err)
			return
		}

		var found bool
		for _, k := range keySlice {

			if k == args[1] {
				found, signingKey = true, k
			}
		}

		if !found {

			err = fmt.Errorf("'%s' key not found\n", args[1])
			return
		}
	}

	signingStrings := GetDefaultSigningStrings()

	signingStrings = AddCustom(signingStrings, custom)

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
			err = fmt.Errorf("ERROR: '%s'\n\n", err)
			return
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
				err = fmt.Errorf("ERROR: unable to open file: '%s'\n\n", err)
				return
			}

			defer func(f io.ReadCloser) {
				err := f.Close()
				if err != nil {
					err = fmt.Errorf("ERROR: '%s'\n", err)
					return
				}
			}(f)

		}
		h := sha256.New()

		if _, err = io.Copy(h, f); err != nil {
			err = fmt.Errorf(
				"ERROR: unable to read file to generate hash: '%s'\n\n",
				err)
			return
		}

		sum = h.Sum(nil)

	}

	// add the public key. This must always be present as it isolates
	// the namespace of even intra-protocol signing.
	var pkb []byte
	pkb, err = cfg.ReadFile(signingKey + "." + PubExt)
	if err != nil {

		err = fmt.Errorf("ERROR: '%s'\n", err)
		return
	}

	// the keychain stores secrets as hex but the pubkeys in nostr npub.
	signingStrings = append(signingStrings, strings.TrimSpace(string(pkb)))

	// append the checksum.
	signingStrings = append(signingStrings, hex.EncodeToString(sum))

	// construct the signing material.
	message := strings.Join(signingStrings, "_")

	if cfg.Verbose {
		PrintErr("signing on message: %s\n", message)
	}

	messageHash := sha256.Sum256([]byte(message))

	var key *secp.SecretKey
	key, err = cfg.GetKey(signingKey, pass)
	if err != nil {
		Fatal("ERROR: '%s'\n", err)
	}

	var sig *schnorr.Signature
	sig, err = schnorr.Sign(key, messageHash[:])

	if skipRandomness {

		if asHex {

			sigStr = hex.EncodeToString(sig.Serialize())

		} else {

			sigStr, err = nostr.EncodeSignature(sig)
			if err != nil {

				err = fmt.Errorf("ERROR: while formatting signature: '%s'\n",
					err)
				return
			}
		}
		return

	} else if sigOnly {

		if asHex {

			sigStr = hex.EncodeToString(sig.Serialize())

		} else {

			sigStr, err = nostr.EncodeSignature(sig)
			if err != nil {

				err = fmt.Errorf("ERROR: while formatting signature: '%s'\n",
					err)
				return
			}
		}

		return
	}
	sigStr, err = FormatSig(signingStrings, sig)
	if err != nil {

		err = fmt.Errorf("ERROR: '%s'\n", err)
		return
	}

	return
}
