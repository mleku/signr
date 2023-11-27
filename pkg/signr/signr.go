package signr

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/gookit/color"
	"mleku.online/git/ec/schnorr"
	"mleku.online/git/signr/pkg/nostr"
)

type SigID int

const (
	// Signature type identifiers
	SchnorrType SigID = iota
	BTCType
)

const (
	// strings used in signature prefix and other places
	AppName                  = "signr"
	ProtocolVersion          = "0"
	DefaultHashFunction      = "SHA256"
	SchnorrSignatures        = "SCHNORR"
	BitcoinCompactSignatures = "ECDSA"
)

const (
	// filename extensions found in app data directory
	ConfigExt  = "yaml"
	DeletedExt = "del"
	ConfigName = "config"
	PubExt     = "pub"
)

const (
	// expected filesystem permissions/masks
	DataDirPerm    os.FileMode = 0700
	ConfigFilePerm os.FileMode = 0600
	KeyFilePerm    os.FileMode = 0400
	DataFileMask   os.FileMode = 0077
)

// SigTypes are the available signature algorithms
var SigTypes = []string{SchnorrSignatures, BitcoinCompactSignatures}

func (s *Signr) GetCfgFilename() string {
	return filepath.Join(s.DataDir, ConfigName+"."+ConfigExt)
}

// GetDefaultSigningStrings returns a slice of strings that forms the prefix of
// a signature/signing material block.
//
// sigType is used to optionally switch to ECDSA bitcoin transaction signatures
// to enable the use case of signing PBSTs for on-chain transactions such as
// anchoring hashes for a chain-bound protocol.
func GetDefaultSigningStrings(sigType ...SigID) (signingStrings []string) {

	var sig SigID
	if len(sigType) > 0 {
		sig = sigType[0]
	}
	signingStrings = []string{
		AppName, ProtocolVersion, DefaultHashFunction, SigTypes[sig],
	}
	return
}

// AddCustom string to the signature string for namespacing purposes.
func (s *Signr) AddCustom(ss []string,
	Custom string) (signingStrings []string) {

	signingStrings = ss
	// Add the custom protocol string to the base if provided:
	if Custom != "" {
		var err error
		Custom, err = s.Sanitize(Custom)
		if err != nil {
			s.Log("error sanitizing custom string: %s\n", err)
			return ss
		}
		// no matter the variation of non-printable characters in the string so
		// long as the printable characters and the positions of their
		// interstitial spaces will be canonical.
		signingStrings = append(signingStrings, Custom)
	}
	return
}

// Sanitize replaces all nonprintable characters with spaces, eliminates spaces
// more than 1 character in a row, removes leading and following spaces and
// finally replaces all remaining interstitial spaces with hyphens.
func (s *Signr) Sanitize(in string) (out string, err error) {

	out = in
	// eliminate all non-printable characters first
	out = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return ' '
	}, out)
	// all multiple non-printables then should be collapsed to single.
	out = strings.Replace(out, "  ", " ", -1)
	// leading and following space characters are removed
	out = strings.TrimSpace(out)
	// spaces are not permitted in custom string, but they could be
	// added, so they will be replaced with hyphens, as are underscores.
	out = strings.ReplaceAll(out, " ", "-")
	if len(out) < 1 {
		err = fmt.Errorf("empty string after sanitizing")
	}
	return
}

// FormatSig takes a slice of signing strings and stitches them together with
// underscores, and snips off the hash and replaces it with the provided
// signature.
func FormatSig(signingStrings []string, sig *schnorr.Signature) (str string,
	err error) {

	prefix := signingStrings[:len(signingStrings)-1]
	var sigStr string
	sigStr, err = nostr.EncodeSignature(sig)
	if err != nil {
		err = fmt.Errorf("error while formatting signature: %s", err)
		return
	}
	return strings.Join(append(prefix, sigStr), "_"), err
}

// Log prints if verbose is enabled, and adds some color if it is enabled.
func (s *Signr) Log(format string, a ...interface{}) {

	if !s.Verbose {
		return
	}
	format = "> " + format
	if s.Color {
		format = color.C256(214).Sprint(format)
	}
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
}

// Err prints an error message, adds some color if enabled.
func (s *Signr) Err(format string, a ...interface{}) {

	if s.Color {
		format = color.Red.Sprint(format)
	}
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
}

// Info prints a message to stderr that won't be picked up by a standard simple
// pipe/redirection.
func (s *Signr) Info(format string, a ...interface{}) {

	if s.Color {
		format = color.Blue.Sprint(format)
	}
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
}

func Newline() {
	_, _ = fmt.Fprintf(os.Stderr, "\n")
}

// Fatal prints an error and then terminates the program.
func (s *Signr) Fatal(format string, a ...interface{}) {

	if s.Color {
		format = color.Red.Sprint("FATAL: ") + format
	}
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

// GetNonceHex returns a random 16 charater hexadecimal string derived from a
// 64 bit random value acquired through the system's strong entropy source.
func (s *Signr) GetNonceHex() (nonceHex string, err error) {

	// add the signature nonce
	nonce := make([]byte, 8)
	_, err = rand.Read(nonce)
	if err != nil {
		err = fmt.Errorf("error getting entropy: %s", err)
		return
	}
	nonceHex = hex.EncodeToString(nonce)
	return
}
