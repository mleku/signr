package signr

import (
	"encoding/hex"
	"fmt"
	"github.com/mleku/bech32"
	secp "github.com/mleku/ec/secp"
	"strings"
	"time"

	"github.com/mleku/ec/schnorr"
	"github.com/mleku/signr/pkg/nostr"
)

type Position int

const prefix = nostr.PubHRP + "1"

const (
	PositionBeginning = iota
	PositionContains
	PositionEnding
)

func (s *Signr) Vanity(str, name string, where Position) (err error) {

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

	found := false
	var sec *secp.SecretKey
	var pub *secp.PublicKey
	var counter int
	var npub string
	started := time.Now()
	for !found {
		counter++
		sec, pub, err = s.GenKeyPair()
		if err != nil {
			return fmt.Errorf("error generating key: %s", err)
		}
		npub, err = nostr.PublicKeyToNpub(pub)
		if err != nil {
			s.Fatal("fatal error generating npub: %s\n", err)
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
	if err = s.Save(name, secBytes, npub); err != nil {
		err = fmt.Errorf("error saving keys: %v", err)
		return
	}
	return
}
