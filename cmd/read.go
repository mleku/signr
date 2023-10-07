package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadFile(name string) (data []byte, err error) {

	path := filepath.Join(dataDir, name)

	// check the permissions are secure first
	var fi os.FileInfo
	fi, err = os.Stat(path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr,
			"error getting file info for %s: %v\n", name, err)
	}

	// secret key files that are readable by other than the owner may not be used
	if fi.Mode().Perm()&0077 != 0 && strings.HasSuffix(name, "."+pubExt) {
		err = fmt.Errorf("secret key has insecure permissions %s",
			fi.Mode().Perm())
		return
	}

	data, err = os.ReadFile(path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr,
			"error reading file '%s': %v\n", name, err)
		return
	}
	return
}
