package signr

import (
	"encoding/hex"
	"github.com/mleku/ec/schnorr"
	"github.com/mleku/signr/pkg/nostr"
)

func (cfg *Config) Gen(name string) {

		sec, pub, err := GenKeyPair()
		if err != nil {

			Fatal("error generating key: '%s'", err)
		}

		secBytes := sec.Serialize()

		npub, _ := nostr.PublicKeyToString(pub)

		if cfg.Verbose {

			pubBytes := schnorr.SerializePubKey(pub)

			PrintErr(
				"generated key pair:\n"+
					"\nhex:\n"+
					"\tsecret: %s\n"+
					"\tpublic: %s\n\n",
				hex.EncodeToString(secBytes),
				hex.EncodeToString(pubBytes),
			)

			nsec, _ := nostr.SecretKeyToString(sec)
			PrintErr(
				"nostr:\n"+
					"\tsecret: %s\n"+
					"\tpublic: %s\n\n", nsec,
				npub)
		}

		if err = Save(cfg, name, secBytes, npub); err != nil {

			Fatal("error saving keys: %v", err)
		}

	return
}