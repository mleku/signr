package signr

import (
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"
	"os"
	"path/filepath"
)

const PassPrompt = "type password to use for secret key (press enter for none): "

func (cfg *Config) Save(name string, secret []byte,
	npub string) (err error) {

	// check that the name isn't already taken

	newPath := filepath.Join(cfg.DataDir, name)
	if _, err := os.Stat(cfg.DataDir); err != nil {
		if !os.IsNotExist(err) {
			PrintErr("'%s' already exists, please use a different name or delete the other.\n",
				newPath)
		}
	}

	var pass1, pass2 []byte
	var tryCount int
	for tryCount < 3 {

		pass1, err = PasswordEntry(PassPrompt, 0)
		if err != nil {

			PrintErr(
				"error in password input: '%s'\n", err)
			return
		}

		if len(pass1) == 0 {
			pass2, err = PasswordEntry("again (press enter again to confirm no encryption): ",
				0)

		} else {

			pass2, err = PasswordEntry("again: ", 0)
		}

		if err != nil {
			PrintErr("error in password input: '%s'\n", err)
			return
		}

		if len(pass1) == 0 && len(pass2) == 0 {
			PrintErr("secret key will not be encrypted\n")
			break
		}

		tryCount++

		if len(pass1) != len(pass2) {
			PrintErr(
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
				break
			}
		}

		if !matched {
			PrintErr(
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

		actualKey := argon2.Key(pass1,
			[]byte("signr"), 3, 1024*1024, 4, 32)
		secret = XOR(secret, actualKey)

		// sanitation
		for i := range pass1 {
			pass1[i] = 0
		}
		for i := range pass2 {
			pass2[i] = 0
		}
	}

	secPath := filepath.Join(cfg.DataDir, name)
	pubPath := secPath + "." + PubExt

	PrintErr("saving secret in '%s', public in '%s'\n", secPath, pubPath)

	passwordProtected := ""
	if len(pass1) > 0 {
		passwordProtected = " *"
	}

	secretString := fmt.Sprintf("%x%s", secret, passwordProtected)

	if cfg.DefaultKey == "" {

		cfg.DefaultKey = name

		viper.Set("default", cfg.DefaultKey)

		if err = viper.WriteConfig(); err != nil {
			PrintErr("error writing config: '%v'\n", err)
		}
	}

	// key files are created with mode 0400, so that when they are deleted, the
	// `rm` command requires a confirmation, in addition to not being readable
	// by any other than the user themselves.
	err = os.WriteFile(secPath, []byte(secretString+"\n"), 0400)
	if err != nil {
		PrintErr(
			"unable to write secret key file '%s': %v\n", secPath, err)
		os.Exit(1)
	}

	err = os.WriteFile(pubPath, []byte(npub+"\n"), 0400)
	if err != nil {
		PrintErr(
			"unable to write public key file '%s': %v\n", pubPath, err)
		os.Exit(1)
	}

	return
}
