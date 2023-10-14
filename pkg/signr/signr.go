package signr

import (
	"fmt"
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

func (cfg *Config) GetCfgFilename() string {
	return filepath.Join(cfg.DataDir, ConfigName+"."+ConfigExt)
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

func PrintErr(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
}

func Fatal(format string, a ...interface{}) {
	PrintErr(format, a...)
	os.Exit(1)
}
