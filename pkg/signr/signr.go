package signr

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/mleku/ec/schnorr"
	"github.com/mleku/signr/pkg/nostr"
	"os"
	"path/filepath"
	"strings"
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

func AddCustom(ss []string, Custom string) (signingStrings []string) {
	signingStrings = ss
	// Add the custom protocol string to the base if provided:
	if Custom != "" {

		// leading and following space characters are removed
		Custom := strings.TrimSpace(Custom)

		// spaces are not permitted in custom string, but they could be
		// added, so they will be replaced with hyphens, as are underscores.
		Custom = strings.ReplaceAll(Custom, " ", "-")
		Custom = strings.ReplaceAll(Custom, "\n", "-")
		Custom = strings.ReplaceAll(Custom, "_", "-")

		signingStrings = append(signingStrings, Custom)
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

