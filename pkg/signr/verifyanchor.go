package signr

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/minio/sha256-simd"
	"mleku.online/git/ec/schnorr"
	secp "mleku.online/git/ec/secp"
	"mleku.online/git/signr/pkg/nostr"
)

// VerifyAnchor takes in the 3 elements found in an anchor transaction
// inscription, the NPUB, MERKLE and NSIG, encoded in hex as a single 256
// character long hex string and splits it into its parts, and validates
// according to the supplied custom protocol string and signr signing material
// protocol.
func (s *Signr) VerifyAnchor(input, custom string) (valid bool, err error) {

	if len(input) != 256 {
		err = fmt.Errorf("error: inscription must be 256 hexadecimal digits")
		return
	}

	// slice up the inscription into its respective segments
	NPUB, merkle, nsig := input[:64], input[64:128], input[128:]

	s.Log("\nnpub   %s\nmerkle %s\nnsig   %s\n", NPUB, merkle, nsig)

	// decode hexadecimal segments to raw bytes
	var pubBytes, sigBytes []byte
	pubBytes, err = hex.DecodeString(NPUB)
	if err != nil {
		err = fmt.Errorf("error: unable to decode public key segment '%s': %s",
			NPUB, err)
		return
	}

	// this only needs to be checked that it is valid hex as it is used as is.
	_, err = hex.DecodeString(merkle)
	if err != nil {
		err = fmt.Errorf("error: unable to decode merkle root segment '%s': %s",
			merkle, err)
		return
	}
	sigBytes, err = hex.DecodeString(nsig)
	if err != nil {
		err = fmt.Errorf("error: unable to decode signature segment '%s': %s",
			nsig, err)
		return
	}

	var pubKey *secp.PublicKey
	pubKey, err = schnorr.ParsePubKey(pubBytes)
	if err != nil {
		err = fmt.Errorf("error: unable to parse public key: %s", err)
		return
	}

	var sig *schnorr.Signature
	sig, err = schnorr.ParseSignature(sigBytes)
	if err != nil {
		err = fmt.Errorf("error: unable to parse signature: %s", err)
		return
	}

	// construct the proper signing material to get the real hash
	signingMaterialStrings := GetDefaultSigningStrings()
	if custom != "" {
		signingMaterialStrings = append(signingMaterialStrings, custom)
	}

	var npub string
	npub, err = nostr.PublicKeyToNpub(pubKey)
	signingMaterialStrings = append(signingMaterialStrings, npub)

	signingMaterialStrings = append(signingMaterialStrings, merkle)
	signingMaterialString := strings.Join(signingMaterialStrings, "_")
	s.Log("signing material reconstructed: %s\n", signingMaterialString)

	signingMaterial := sha256.Sum256(
		[]byte(signingMaterialString))

	valid = sig.Verify(signingMaterial[:], pubKey)
	return
}
