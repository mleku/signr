package cmd

import (
	"golang.org/x/term"
	"syscall"
)

func PasswordEntry(prompt string) (pass []byte, err error) {

	PrintErr(prompt)

	pass, err = term.ReadPassword(int(syscall.Stdin))

	PrintErr("\n")

	return
}
