package signr

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gookit/color"
	"github.com/mleku/ec/schnorr"
	"github.com/mleku/signr/pkg/nostr"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

const (
	AppName                    = "signr"
	ConfigExt                  = "yaml"
	DeletedExt                 = "del"
	ConfigName                 = "config"
	PubExt                     = "pub"
	DataDirPerm    os.FileMode = 0700
	ConfigFilePerm os.FileMode = 0600
	KeyFilePerm    os.FileMode = 0400
	DataFileMask   os.FileMode = 0077
)

func (s *Signr) GetCfgFilename() string {
	return filepath.Join(s.DataDir, ConfigName+"."+ConfigExt)
}

func GetDefaultSigningStrings() (signingStrings []string) {

	// for now the first 4 are always the same
	signingStrings = []string{
		"signr", "0", "SHA256", "SCHNORR",
	}
	return
}

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

	// eliminate all non-printable characters first
	in = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return ' '
	}, in)

	// all multiple non-printables then should be collapsed to single.
	in = strings.Replace(in, "  ", " ", -1)

	// leading and following space characters are removed
	in = strings.TrimSpace(in)

	// spaces are not permitted in custom string, but they could be
	// added, so they will be replaced with hyphens, as are underscores.
	in = strings.ReplaceAll(in, " ", "-")

	if len(in) < 1 {
		err = fmt.Errorf("empty string after sanitizing")
	}

	return
}

func FormatSig(signingStrings []string, sig *schnorr.Signature) (str string,
	err error) {

	prefix := signingStrings[:len(signingStrings)-1]

	var sigStr string
	sigStr, err = nostr.EncodeSignature(sig)
	if err != nil {

		err = fmt.Errorf("ERROR: while formatting signature: '%s'\n", err)
		return
	}

	return strings.Join(append(prefix, sigStr), "_"), err
}

func (s *Signr) Log(format string, a ...interface{}) {
	if !s.Verbose {
		return
	}
	if s.Color {
		format = color.C256(214).Sprint("> ") + format
	}
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
}

func (s *Signr) Err(format string, a ...interface{}) {
	if !s.Verbose {
		return
	}
	if s.Color {
		format = color.Red.Sprint("! ") + format
	}
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
}

func (s *Signr) Fatal(format string, a ...interface{}) {
	if s.Color {
		format = color.Red.Sprint("FATAL: ") + format
	}
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func (s *Signr) PrintErr(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
}

func (s *Signr) GetNonceHex() (nonceHex string, err error) {

	// add the signature nonce
	nonce := make([]byte, 8)

	_, err = rand.Read(nonce)
	if err != nil {

		err = fmt.Errorf("error getting entropy: %s\n", err)
		return
	}

	nonceHex = hex.EncodeToString(nonce)
	return
}
