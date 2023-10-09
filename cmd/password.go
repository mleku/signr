package cmd

import (
	"golang.org/x/term"
	"syscall"
)

func PasswordEntry(prompt string) (pass []byte, err error) {
	printErr(prompt)
	pass, err = term.ReadPassword(int(syscall.Stdin))
	printErr("\n")
	if err != nil {
		return
	}
	return
}
