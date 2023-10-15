package signr

import (
	"fmt"
	"github.com/minio/sha256-simd"
	"io"
	"os"
)

// HashFile accepts a filename, interpreting "-" to mean to read from stdin, and
// computes the sha256 hash of the file.
func HashFile(filename string) (sum []byte, err error) {

	var f io.ReadCloser

	switch {
	case filename == "-":

		// read from stdin
		f = os.Stdin

	default:

		// read from the file
		f, err = os.Open(filename)
		if err != nil {
			err = fmt.Errorf("ERROR: unable to open file: %s\n\n", err)
			return
		}

		defer func(f io.ReadCloser) {

			err := f.Close()
			if err != nil {
				err = fmt.Errorf("error while closing file: %s\n", err)
				return
			}
		}(f)

	}

	h := sha256.New()

	// run the file data through the hash function
	if _, err = io.Copy(h, f); err != nil {
		err = fmt.Errorf(
			"ERROR: unable to read file to generate hash: %s\n\n",
			err)
		return
	}

	sum = h.Sum(nil)
	return
}
