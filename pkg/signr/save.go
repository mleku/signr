package signr

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const PassPrompt = "type password to use for secret key (press enter for none): "

func (s *Signr) Equal(first, second []byte) (same bool) {

	if len(first) != len(second) {
		s.Log("length of inputs differs: first: %d; second %d",
			len(first) != len(second))
		return
	}
	for i := range first {
		if first[i] != second[i] {
			return
		}
	}
	same = true
	return
}

func Zero(bytes []byte) {
	for i := range bytes {
		bytes[i] = 0
	}
}

func (s *Signr) Save(name string, secret []byte,
	npub string) (err error) {

	// check that the name isn't already taken
	newPath := filepath.Join(s.DataDir, name)
	var exists bool
	_, exists, err = CheckFileExists(newPath)
	if err != nil {
		err = fmt.Errorf("error checking for existence of file %s: %s", newPath,
			err)
		return
	}
	if exists {
		err = fmt.Errorf("'%s' already exists, please use a different name "+
			"or delete the other", newPath)
		return
	}
	var pass1, pass2 []byte
	var tryCount int
	const maxTries = 3
	for tryCount < maxTries {
		pass1, err = s.PasswordEntry(PassPrompt, s.PassEntryType)
		if err != nil {
			Zero(pass1)
			err = fmt.Errorf("error in password input: %s", err)
			return
		}
		// s.Log("'%s'\n", pass1)
		if len(pass1) == 0 {
			pass2, err = s.PasswordEntry(
				"again (press enter \nagain to confirm no encryption): ",
				s.PassEntryType)
		} else {
			pass2, err = s.PasswordEntry("again: ", s.PassEntryType)
		}
		if err != nil {
			Zero(pass2)
			err = fmt.Errorf("error in password input: %s", err)
			return
		}
		// s.Log("'%s'\n", pass2)
		if len(pass1) == 0 && len(pass2) == 0 {
			s.Log("secret key will not be encrypted\n")
			break
		}
		tryCount++
		if matched := s.Equal(pass1, pass2); !matched {
			s.Err("passwords didn't match, try again (try %d of %d)\n",
				tryCount+1)
			// sanitation
			Zero(pass1)
			Zero(pass2)
		} else {
			break
		}
	}
	passwordProtected := ""
	if len(pass1) > 0 {
		passwordProtected = " *"
		// sanitation
		Zero(pass2)
		actualKey := ArgonKey(pass1)
		secret = s.XOR(secret, actualKey)
		// sanitation
		Zero(pass1)
		Zero(actualKey)
	}
	secPath := filepath.Join(s.DataDir, name)
	pubPath := secPath + "." + PubExt
	s.Log("saving secret in '%s', public in '%s'\n", secPath, pubPath)
	secretString := fmt.Sprintf("%x%s", secret, passwordProtected)
	if s.DefaultKey == "" {
		s.DefaultKey = name
		viper.Set("default", s.DefaultKey)
		if err = viper.WriteConfig(); err != nil {
			s.Err("error writing config: '%v'\n", err)
		}
	}
	// key files are created with mode 0400, so that when they are deleted, the
	// `rm` command requires a confirmation, in addition to not being readable
	// by any other than the user themselves.
	err = os.WriteFile(secPath, []byte(secretString+"\n"), KeyFilePerm)
	if err != nil {
		err = fmt.Errorf("unable to write secret key file '%s': %v",
			secPath, err)
		return
	}
	err = os.WriteFile(pubPath, []byte(npub+"\n"), KeyFilePerm)
	if err != nil {
		err = fmt.Errorf("unable to write public key file '%s': %v",
			pubPath, err)
		return
	}
	return
}
