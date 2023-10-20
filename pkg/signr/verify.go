package signr

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/minio/sha256-simd"
	"github.com/mleku/ec/schnorr"
	secp "github.com/mleku/ec/secp"
	"github.com/mleku/signr/pkg/nostr"
)

func (s *Signr) Verify(filename, sigOrSigFile, PubKey,
	Custom string) (valid bool, err error) {

	var pubKey, sumHex, nonce, pkInSig string
	var signingStrings []string
	var sig *schnorr.Signature
	var sum []byte
	pubKey = PubKey
	// get the hash of the file that has been signed
	if sum, err = HashFile(filename); err != nil {
		err = fmt.Errorf("error getting hash of file/input: %s", err)
		return
	}
	sumHex = hex.EncodeToString(sum)
	if sig, pkInSig, nonce = s.RecogniseSig(sigOrSigFile); pkInSig != "" {
		s.Log("pubkey in signature: %s\n", pkInSig)
		pubKey = pkInSig
	}
	s.Log("loading pubkey: %s\n", pubKey)
	// if the signature was not in the parameter, we try to get one from a
	// file named in the parameter.
	if sig == nil {
		var data []byte
		data, err = s.ReadFile(sigOrSigFile)
		if err != nil {
			err = fmt.Errorf("error reading file: %s", err)
			return
		}
		sig, pkInSig, nonce = s.RecogniseSig(string(data))
		if pkInSig != "" {
			s.Log("pubkey in signature: %s\n", pkInSig)
			pubKey = pkInSig
		}
	}
	// decode the public key
	var pk *secp.PublicKey
	if pk, err = nostr.NpubToPublicKey(pubKey); err != nil {
		err = fmt.Errorf("error decoding pubkey: %s", err)
		return
	}
	// start with the essentials
	signingStrings = GetDefaultSigningStrings()
	if Custom != "" {
		s.Log("adding custom namespace %s\n", Custom)
		// if a custom protocol field is specified, it goes before the pubkey:
		signingStrings = s.AddCustom(signingStrings, Custom)
	}
	// if a nonce was present in the signature, it must be added to the signing
	// material
	if nonce != "" {
		s.Log("adding nonce %s\n", Custom)
		signingStrings = append(signingStrings, nonce)
	}
	// add the public key to the strings
	signingStrings = append(signingStrings, pubKey)
	// append the checksum of the file/input to the end of the signing
	// string slice
	signingStrings = append(signingStrings, sumHex)
	// generate the signing material
	message := strings.Join(signingStrings, "_")
	s.Log("message: %s\n", message)
	// hash the signing material
	messageHash := sha256.Sum256([]byte(message))
	// verify the hash and the signature match the public key
	valid = sig.Verify(messageHash[:], pk)
	return
}

var FullSigPrefix = strings.Join(GetDefaultSigningStrings(), "_")

func (s *Signr) RecogniseSig(possibleSig string) (sig *schnorr.Signature,
	pubKey, nonce string) {

	var err error
	switch {
	case strings.HasPrefix(possibleSig, FullSigPrefix):
		signingStrings := strings.Split(possibleSig, "_")
		// if it is a full signature, the last part will match to the next:
		possibleSig = signingStrings[len(signingStrings)-1]
		// scan for the possible pubkey segment and add to returns it if found.
		for i, pk := range signingStrings {
			if strings.HasPrefix(pk, nostr.PubHRP) {
				pubKey = pk
				// before the pubkey can be a nonce also, return it if it
				// decodes as hex
				if i > 0 {
					if _, err = hex.DecodeString(signingStrings[i-1]); err == nil {
						nonce = signingStrings[i-1]
						s.Log("nonce found %s\n", nonce)
					}
				}
			}
		}
		// if it's a sig, the next case will decode it.
		fallthrough
	case strings.HasPrefix(possibleSig, nostr.SigHRP):
		// decode the signature
		if sig, err = nostr.DecodeSignature(possibleSig); err != nil {
			s.Log("not possible sig: %s\n")
		}
		// a hex signature can only be an exact number of characters long
	case len(possibleSig) == 128:
		var sigBytes []byte
		if sigBytes, err = hex.DecodeString(possibleSig); err != nil {
			s.Log("not a possible hex signature: %s\n", err)
		}
		if sig, err = schnorr.ParseSignature(sigBytes); err != nil {
			s.Log("not possible nsig: %s\n")
		}
	}
	return
}
