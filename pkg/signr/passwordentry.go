package signr

import (
	"golang.org/x/crypto/ssh/terminal"
)

const (
	PasswordEntryViaTTY = iota
)

func PasswordEntry(prompt string, entryType int) (pass []byte, err error) {

	switch entryType {
	case PasswordEntryViaTTY:

		PrintErr(prompt)
		pass, err = terminal.ReadPassword(1)
		PrintErr("\n")
	}

	return
}
