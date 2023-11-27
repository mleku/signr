package signr

import (
	"encoding/hex"
	"fmt"

	"mleku.online/git/ec/schnorr"
	secp "mleku.online/git/ec/secp"
	"mleku.online/git/signr/pkg/btc"
)

func (s *Signr) GetAnchor(args []string, pass, custom string) (WIF,
	NPUB, MERKLE, NSIG string, err error) {

	// args[0] should be a 32 byte hash, presumably a merkle root
	MERKLE = args[0]
	if len(MERKLE) != 64 {
		err = fmt.Errorf("merkle/hash not expected length of 32 byes/64 hex characters, got %d characters",
			len(MERKLE))
		return
	}
	var merkleBytes []byte
	if merkleBytes, err = hex.DecodeString(MERKLE); err != nil {
		err = fmt.Errorf("error decoding merkle/hash: %s", err)
		return
	}
	// we encode it back because then it's for sure the same case as the rest
	// of the outputs we need to generate.
	MERKLE = hex.EncodeToString(merkleBytes)

	var keyName string
	if len(args) > 1 {
		keyName = args[1]
	} else {
		keyName = s.DefaultKey
	}
	var key *secp.SecretKey
	NSIG, key, err = s.Sign([]string{MERKLE, keyName}, pass, custom, true,
		false)
	if err != nil {
		err = fmt.Errorf("error signing merkle/hash: %s", err)
		return
	}
	WIF = btc.GetWIF(key)
	NPUB = hex.EncodeToString(schnorr.SerializePubKey(key.PubKey()))
	return
}
