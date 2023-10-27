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
	SigHRP = "nsig"
)

// ConvertForBech32 performs the bit expansion required for encoding into
// Bech32.
func ConvertForBech32(b8 []byte) (b5 []byte, err error) {
	return bech32.ConvertBits(b8, 8, 5, true)
}

// ConvertFromBech32 collapses together the bit expanded 5 bit numbers encoded
// in bech32.
func ConvertFromBech32(b5 []byte) (b8 []byte, err error) {
	return bech32.ConvertBits(b5, 5, 8, true)
}

// SecretKeyToNsec encodes an secp256k1 secret key as a Bech32 string (nsec).
func SecretKeyToNsec(sk *secp.SecretKey) (encoded string, err error) {

	var b5 []byte
	if b5, err = ConvertForBech32(sk.Serialize()); err != nil {
		return
	}
	return bech32.Encode(SecHRP, b5)
}

// PublicKeyToNpub encodes a public kxey as a bech32 string (npub).
func PublicKeyToNpub(pk *secp.PublicKey) (encoded string, err error) {

	var bits5 []byte
	if bits5, err = ConvertForBech32(schnorr.SerializePubKey(pk)); err != nil {
		return
	}
	return bech32.Encode(PubHRP, bits5)
}

// NsecToSecretKey decodes a nostr secret key (nsec) and returns the secp256k1
// secret key.
func NsecToSecretKey(encoded string) (sk *secp.SecretKey, err error) {

	var b5, b8 []byte
	var hrp string
	hrp, b5, err = bech32.Decode(encoded)
	if err != nil {
		return
	}
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

// NpubToPublicKey decodes an nostr public key (npub) and returns an secp256k1
// public key.
func NpubToPublicKey(encoded string) (pk *secp.PublicKey, err error) {

	var b5, b8 []byte
	var hrp string
	hrp, b5, err = bech32.Decode(encoded)
	if err != nil {
		err = fmt.Errorf("ERROR: '%s'", err)
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

// EncodeSignature encodes a schnorr signature as Bech32 with the HRP "nsig" to
// be consistent with the key encodings 4 characters starting with 'n'.
func EncodeSignature(sig *schnorr.Signature) (str string, err error) {

	var b5 []byte
	b5, err = ConvertForBech32(sig.Serialize())
	if err != nil {
		err = fmt.Errorf("ERROR: '%s'", err)
		return
	}
	str, err = bech32.Encode(SigHRP, b5)
	return
}

// DecodeSignature decodes a Bech32 encoded nsig nostr (schnorr) signature into
// its runtime binary form.
func DecodeSignature(encoded string) (sig *schnorr.Signature, err error) {

	var b5, b8 []byte
	var hrp string
	hrp, b5, err = bech32.DecodeNoLimit(encoded)
	if err != nil {
		err = fmt.Errorf("ERROR: '%s'", err)
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
