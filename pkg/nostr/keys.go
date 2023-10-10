package nostr

import (
	"fmt"
	"github.com/mleku/bech32"
	"github.com/mleku/ec/schnorr"
	secp "github.com/mleku/ec/secp"
)

const (
	SecHRP = "nsec"
	PubHRP = "npub"
	SigHRP = "sig"
)

func ConvertForBech32(b8 []byte) (b5 []byte, err error) {
	return bech32.ConvertBits(b8, 8, 5, true)
}

func ConvertFromBech32(b5 []byte) (b8 []byte, err error) {
	return bech32.ConvertBits(b5, 5, 8, true)
}

func SecretKeyToString(sk *secp.SecretKey) (encoded string, err error) {

	var b5 []byte
	if b5, err = ConvertForBech32(sk.Serialize()); err != nil {
		return
	}

	return bech32.Encode(SecHRP, b5)
}

func PublicKeyToString(pk *secp.PublicKey) (encoded string, err error) {

	var bits5 []byte
	if bits5, err = ConvertForBech32(schnorr.SerializePubKey(pk)); err != nil {
		return
	}

	return bech32.Encode(PubHRP, bits5)
}

func DecodeSecretKey(encoded string) (sk *secp.SecretKey, err error) {
	var b5, b8 []byte
	var hrp string
	hrp, b5, err = bech32.Decode(encoded)
	if hrp != SecHRP {
		err = fmt.Errorf("wrong human readable part, got '%s' want '%s'",
			hrp, SecHRP)
		return
	}
	b8, err = ConvertFromBech32(b5)
	if err != nil {
		return
	}
	sk = secp.SecKeyFromBytes(b8)
	return
}

func DecodePublicKey(encoded string) (pk *secp.PublicKey, err error) {
	var b5, b8 []byte
	var hrp string
	hrp, b5, err = bech32.Decode(encoded)
	if err != nil {
		err = fmt.Errorf("ERROR: '%s'\n", err)
		return
	}
	if hrp != PubHRP {
		err = fmt.Errorf("wrong human readable part, got '%s' want '%s'",
			hrp, PubHRP)
		return
	}
	b8, err = ConvertFromBech32(b5)
	if err != nil {
		return
	}
	return schnorr.ParsePubKey(b8[:32])
}

func EncodeSignature(sig *schnorr.Signature) (str string, err error) {
	var b5 []byte
	b5, err = ConvertForBech32(sig.Serialize())
	if err != nil {
		err = fmt.Errorf("ERROR: '%s'\n", err)
		return
	}
	str, err = bech32.Encode(SigHRP, b5)
	return
}

func DecodeSignature(encoded string) (sig *schnorr.Signature, err error) {

	var b5, b8 []byte
	var hrp string
	hrp, b5, err = bech32.DecodeNoLimit(encoded)
	if err != nil {
		fmt.Printf("ERROR: '%s'\n", err)
		return
	}
	if hrp != SigHRP {
		err = fmt.Errorf("wrong human readable part, got '%s' want '%s'",
			hrp, SigHRP)
		return
	}
	b8, err = ConvertFromBech32(b5)
	if err != nil {
		return
	}
	return schnorr.ParseSignature(b8[:64])
}
