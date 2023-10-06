package cmd

import (
	"fmt"
	"golang.org/x/term"
	"syscall"
)

func PasswordEntry(prompt string) (pass []byte, err error) {
	fmt.Print(prompt)
	pass, err = term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return
	}
	return
}
