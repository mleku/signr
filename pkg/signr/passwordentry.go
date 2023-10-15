package signr

import (
	"golang.org/x/crypto/ssh/terminal"
)

const (
	PasswordEntryViaTTY = iota
	// other entry types (eg, GUIs) can be added here.
)

func (s *Signr) PasswordEntry(prompt string, entryType int) (pass []byte, err error) {

	switch entryType {
	case PasswordEntryViaTTY:

		s.PrintErr(prompt)
		pass, err = terminal.ReadPassword(1)
		s.PrintErr("\n")

	default:
		s.PrintErr("password entry type %d not implemented\n", entryType)
	}

	return
}
