package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"
	"os"
	"path/filepath"
)

const passPrompt = "type password to use for secret key (press enter for none): "

func Save(name string, secret []byte, npub string) (err error) {
	var pass1, pass2 []byte
	var tryCount int
	for tryCount < 3 {
		pass1, err = PasswordEntry(passPrompt)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr,
				"error in password input: '%s'\n", err)
			return
		}
		if len(pass1) == 0 {
			pass2, err = PasswordEntry(
				"again (press enter again to confirm no encryption): ")
		} else {
			pass2, err = PasswordEntry("again: ")
		}
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr,
				"error in password input: '%s'\n", err)
			return
		}
		if len(pass1) == 0 && len(pass2) == 0 {
			_, _ = fmt.Fprintf(os.Stderr,
				"secret key will not be encrypted\n")
			break
		}
		tryCount++
		if len(pass1) != len(pass2) {
			_, _ = fmt.Fprintf(os.Stderr,
				"passwords don't match, try again (try %d of 3)\n",
				tryCount+1)
			// sanitation
			for i := range pass1 {
				pass1[i] = 0
			}
			for i := range pass2 {
				pass2[i] = 0
			}
			continue
		}
		var matched = true
		for i := range pass1 {
			if pass1[i] != pass2[i] {
				matched = false
			}
		}
		if !matched {
			_, _ = fmt.Fprintf(os.Stderr,
				"passwords didn't match, try again (try %d of 3)\n",
				tryCount+1)
			// sanitation
			for i := range pass1 {
				pass1[i] = 0
			}
			for i := range pass2 {
				pass2[i] = 0
			}
		} else {
			break
		}
	}
	if len(pass1) > 0 {
		actualKey := argon2.Key([]byte(pass1), []byte("signr"),
			3, 32*1024, 4, 32)
		secret = xor(secret, actualKey)
		// sanitation
		for i := range pass1 {
			pass1[i] = 0
		}
		for i := range pass2 {
			pass2[i] = 0
		}
	}

	secPath := filepath.Join(dataDir, name)
	pubPath := secPath + "." + pubExt
	_, _ = fmt.Fprintf(os.Stderr,
		"saving secret key in '%s', public key in '%s'\n",
		secPath, pubPath)

	prefix := "c"
	if len(pass1) > 0 {
		prefix = "e"
	}
	secretString := fmt.Sprintf("%s%x", prefix, secret)
	_ = secretString
	if defaultKey == "" {
		defaultKey = name
		viper.Set("default", defaultKey)
		if err = viper.SafeWriteConfig(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr,
				"error: '%v'\n", err)
		}
	}
	err = os.WriteFile(secPath, []byte(secretString+"\n"), 0600)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr,
			"unable to write secret key file '%s': %v\n", secPath, err)
		os.Exit(1)
	}
	err = os.WriteFile(pubPath, []byte(npub+"\n"), 0600)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr,
			"unable to write public key file '%s': %v\n", pubPath, err)
		os.Exit(1)
	}

	return
}

func xor(dest, src []byte) []byte {
	if len(src) != len(dest) {
		panic("key and secret must be the same length")
	}
	for i := range dest {
		dest[i] = dest[i] ^ src[i]
	}
	return dest
}
