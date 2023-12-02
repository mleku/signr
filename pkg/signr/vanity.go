package signr

import (
	"encoding/hex"
	"fmt"
	"mleku.online/git/bech32"
	secp "mleku.online/git/ec/secp"
	"strings"
	"time"

	"mleku.online/git/ec/schnorr"
	"mleku.online/git/qu"
	"mleku.online/git/signr/pkg/nostr"
)

type Position int

const prefix = nostr.PubHRP + "1"

const (
	PositionBeginning = iota
	PositionContains
	PositionEnding
)

func (s *Signr) Vanity(str, name string, where Position) (e error) {

	// check the string has valid bech32 ciphers
	for i := range str {
		wrong := true
		for j := range bech32.Charset {
			if str[i] == bech32.Charset[j] {
				wrong = false
				break
			}
		}
		if wrong {
			return fmt.Errorf("found invalid character '%c' only ones from '%s' allowed\n",
				str[i], bech32.Charset)
		}
	}

	started := time.Now()

	quit := qu.T()
	var sec *secp.SecretKey
	var counter int
	var npub string
	var pub *secp.PublicKey
	sec, pub, npub, counter, e = s.mine(str, where, quit)

	s.Info("generated in %d attempts, taking %v\n", counter,
		started.Sub(time.Now()))
	secBytes := sec.Serialize()
	if s.Verbose {
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
	if e = s.Save(name, secBytes, npub); e != nil {
		e = fmt.Errorf("error saving keys: %v", e)
		return
	}
	return
}

func (s *Signr) mine(str string, where Position,
	quit qu.C) (sec *secp.SecretKey, pub *secp.PublicKey,
	npub string, counter int, e error) {

	found := false
	for !found {
		counter++
		sec, pub, e = s.GenKeyPair()
		if e != nil {
			e = fmt.Errorf("error generating key: %s", e)
			return
		}
		npub, e = nostr.PublicKeyToNpub(pub)
		if e != nil {
			s.Fatal("fatal error generating npub: %s\n", e)
		}
		// s.Log("%s\n", npub)
		switch where {
		case PositionBeginning:
			if strings.HasPrefix(npub, prefix+str) {
				found = true
			}
		case PositionEnding:
			if strings.HasSuffix(npub, str) {
				found = true
			}
		case PositionContains:
			if strings.Contains(npub, str) {
				found = true
			}
		}
		if counter%1000000 == 0 {
			s.Log("attempt %d\n", counter)
		}
	}
	return
}
