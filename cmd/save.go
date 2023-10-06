package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	PubExtension = ".pub"
)

func Save(name string) {

	secPath := filepath.Join(dataDir, name)
	pubPath := secPath + PubExtension
	fmt.Fprintf(os.Stderr,
		"saving secret key in '%s', public key in '%s'\n",
		secPath, pubPath)

}
