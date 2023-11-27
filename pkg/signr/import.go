package signr

import (
	"encoding/hex"
	"fmt"
	"strings"

	"mleku.online/git/ec/schnorr"
	secp "mleku.online/git/ec/secp"
	"mleku.online/git/signr/pkg/nostr"
)

func (s *Signr) Import(secKey, name string) (err error) {

	var sec *secp.SecretKey
	if strings.HasPrefix(secKey, nostr.SecHRP) {
		if sec, err = nostr.NsecToSecretKey(secKey); err != nil {
			err = fmt.Errorf("ERROR: while decoding key: '%v'", err)
			return
		}
	} else {
		var secBytes []byte
		if secBytes, err = hex.DecodeString(secKey); err != nil {
			err = fmt.Errorf("key is mangled, '%s', cannot decode: '%v'",
				secKey, err)
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
	npub, _ := nostr.PublicKeyToNpub(pub)
	s.Log("hex:\n"+
		"\tsecret: %s\n"+
		"\tpublic: %s\n",
		hex.EncodeToString(secBytes),
		hex.EncodeToString(schnorr.SerializePubKey(pub)),
	)
	if s.Verbose {
		nsec, _ := nostr.SecretKeyToNsec(sec)
		s.Err("nostr:\n"+
			"\tsecret: %s\n"+
			"\tpublic: %s\n\n",
			nsec, npub)
	}
	if err = s.Save(name, secBytes, npub); err != nil {
		err = fmt.Errorf("error saving keys: %v", err)
	}
	return
}
