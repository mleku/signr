package nostr

import (
	"fmt"
	"github.com/mleku/bech32"
	"github.com/mleku/ec/schnorr"
	secp "github.com/mleku/ec/secp"
)

const (
	secHRP = "nsec"
	pubHRP = "npub"
)

func ConvertForBech32(b8 []byte) (b5 []byte, err error) {
	return bech32.ConvertBits(b8, 8, 5, true)
}

func SecretKeyToString(sk *secp.SecretKey) (encoded string, err error) {

	var bits5 []byte
	if bits5, err = ConvertForBech32(sk.Serialize()); err != nil {
		return
	}

	return bech32.Encode(secHRP, bits5)
}

func PublicKeyToString(pk *secp.PublicKey) (encoded string, err error) {

	var bits5 []byte
	if bits5, err = ConvertForBech32(schnorr.SerializePubKey(pk)); err != nil {
		return
	}

	return bech32.Encode(pubHRP, bits5)
}

func DecodeSecretKey(encoded string) (sk *secp.SecretKey, err error) {
	var b []byte
	var hrp string
	hrp, b, err = bech32.Decode(encoded)
	if hrp != secHRP {
		err = fmt.Errorf("wrong human readable part, got '%s' want '%s'",
			hrp, secHRP)
		return
	}
	sk = secp.SecKeyFromBytes(b)
	return
}

func DecodePublicKey(encoded string) (sk *secp.PublicKey, err error) {
	var b []byte
	var hrp string
	hrp, b, err = bech32.Decode(encoded)
	if hrp != pubHRP {
		err = fmt.Errorf("wrong human readable part, got '%s' want '%s'",
			hrp, pubHRP)
		return
	}
	return schnorr.ParsePubKey(b)
}
