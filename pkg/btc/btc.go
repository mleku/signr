package btc

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/minio/sha256-simd"
	"github.com/mleku/base58"
	btcec "github.com/mleku/ec"
	"github.com/mleku/ec/chainhash"
	secp "github.com/mleku/ec/secp"
)

const (
	SecretKeyWIFPrefix = 0x80
)

// GetWIF converts the signr secret key to WIF format for importing into a
// bitcoin wallet.
func GetWIF(sk *secp.SecretKey) (secKey string) {

	// get the raw bytes
	keyBytes := sk.Serialize()
	// hash the key twice
	check := sha256.Sum256(keyBytes)
	check = sha256.Sum256(check[:])
	checkBytes := check[:4]
	// key bytes for base58 check = 0x80 | keyBytes |
	keyBytes = append(append(
		[]byte{SecretKeyWIFPrefix}, keyBytes...), checkBytes...)
	secKey = base58.Encode(keyBytes)

	return
}

// GetTaprootKeys takes a signr secret key and generates the taproot secret and
// public keys for a generic spend, as can be used to embed a signature in the
// witness.
//
// This reimplements the functionality in
// github.com/btcsuite/btcd/txscript/taproot.go except without using btcd and
// co.
func GetTaprootAddress(sk *secp.SecretKey,
	scriptRoot []byte) (taprootAddress string, err error) {

	// if the secret key's pubkey has an odd y coord, negate the seckey,
	// as per BIP 341
	skScalar := sk.Key
	pkBytes := sk.PubKey().SerializeCompressed()
	if pkBytes[0] == secp.PubKeyFormatCompressedOdd {
		skScalar.Negate()
	}

	// compute the tweak hash that commits to the internal key and the merkle
	// script root.
	schnorrKeyBytes := pkBytes[1:]
	tapHash := chainhash.TaggedHash(chainhash.TagTapTweak, schnorrKeyBytes,
		scriptRoot)

	// create a ModNScalar from the secret key
	var tweakScalar secp.ModNScalar
	tweakScalar.SetBytes((*[32]byte)(tapHash))

	tapSecKey := btcec.SecKeyFromScalar(skScalar.Add(&tweakScalar))
	tapPubKey := tapSecKey.PubKey()
	var tapAddr *btcutil.AddressTaproot
	tapAddr, err = btcutil.NewAddressTaproot(tapPubKey.SerializeCompressed(),
		&chaincfg.MainNetParams)
	if err != nil {
		err = fmt.Errorf("error creating taproot address: %s", err)
		return
	}

	// then serialize to the appropriate bech32 strings
	tapPubKey.SerializeCompressed()

	taprootAddress = tapAddr.EncodeAddress()

	return
}
