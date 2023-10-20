package signr

import (
	"encoding/hex"
	"fmt"

	"github.com/mleku/ec/schnorr"
	"github.com/mleku/signr/pkg/nostr"
)

func (s *Signr) Gen(name string) (err error) {

	sec, pub, err := s.GenKeyPair()
	if err != nil {
		err = fmt.Errorf("error generating key: %s", err)
		return
	}
	secBytes := sec.Serialize()
	npub, _ := nostr.PublicKeyToNpub(pub)
	s.Log(
		"generated key pair:\n"+
			"\nhex:\n"+
			"\tsecret: %s\n"+
			"\tpublic: %s\n\n",
		hex.EncodeToString(secBytes),
		hex.EncodeToString(schnorr.SerializePubKey(pub)),
	)
	if s.Verbose {
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
