package signr

import (
	"encoding/hex"
	"fmt"
	"github.com/mleku/ec/schnorr"
	secp "github.com/mleku/ec/secp"
	"github.com/mleku/signr/pkg/nostr"
	"strings"
)

func (cfg *Config) Import(secKey, name string) (err error) {

	var sec *secp.SecretKey
	if strings.HasPrefix(secKey, nostr.SecHRP) {

		if sec, err = nostr.DecodeSecretKey(secKey); err != nil {

			err = fmt.Errorf("ERROR: while decoding key: '%v'", err)
		}

	} else {

		var secBytes []byte
		if secBytes, err = hex.DecodeString(secKey); err != nil {

			err = fmt.Errorf(
				"key is mangled, '%s', cannot decode: '%v'", secKey, err)
			return
		}

		sec = secp.SecKeyFromBytes(secBytes)
	}

	if sec == nil {

		err = fmt.Errorf("input did not match any known formats")
		return
	}

	pub := sec.PubKey()
	secBytes := sec.Serialize()

	npub, _ := nostr.PublicKeyToString(pub)

	if cfg.Verbose {

		pubBytes := schnorr.SerializePubKey(pub)

		PrintErr("hex:\n"+
			"\tsecret: %s\n"+
			"\tpublic: %s\n",
			hex.EncodeToString(secBytes),
			hex.EncodeToString(pubBytes),
		)

		nsec, _ := nostr.SecretKeyToString(sec)

		PrintErr("nostr:\n"+
			"\tsecret: %s\n"+
			"\tpublic: %s\n\n",
			nsec, npub)
	}

	if err = cfg.Save(name, secBytes, npub); err != nil {

		err = fmt.Errorf("error saving keys: %v", err)
	}
	return
}
