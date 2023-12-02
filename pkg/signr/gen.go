package signr

import (
	"encoding/hex"
	"fmt"

	"mleku.online/git/ec/schnorr"
	"mleku.online/git/signr/pkg/nostr"
)

func (s *Signr) Generate(name string) (err error) {

	sec, pub, err := s.GenKeyPair()
	if err != nil {
		err = fmt.Errorf("error generating key: %s", err)
		return
	}
	secBytes := sec.Serialize()
	var npub string
	npub, err = nostr.PublicKeyToNpub(pub)
	if err != nil {
		err = fmt.Errorf("error generating npub: %s", err)
		return
	}
	if s.Verbose.Load() {
		s.Log(
			"generated key pair:\n"+
				"\nhex:\n"+
				"\tsecret: %s\n"+
				"\tpublic: %s\n\n",
			hex.EncodeToString(secBytes),
			hex.EncodeToString(schnorr.SerializePubKey(pub)),
		)
		nsec, _ := nostr.SecretKeyToNsec(sec)
		s.Log("nostr:\n"+
			"\tsecret: %s\n"+
			"\tpublic: %s\n\n",
			nsec, npub)
	}
	if err = s.Save(name, secBytes, npub); err != nil {
		err = fmt.Errorf("error saving keys: %v", err)
		return
	}
	return
}
