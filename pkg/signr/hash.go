package signr

import (
	"fmt"
	"io"
	"os"

	"github.com/minio/sha256-simd"
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
		if f, err = os.Open(filename); err != nil {
			err = fmt.Errorf("error: unable to open file: %s", err)
			return
		}
		defer func(f io.ReadCloser) {
			if err = f.Close(); err != nil {
				err = fmt.Errorf("error while closing file: %s", err)
				return
			}
		}(f)
	}
	h := sha256.New()
	// run the file data through the hash function
	if _, err = io.Copy(h, f); err != nil {
		err = fmt.Errorf("error: unable to read file to generate hash: %s",
			err)
		return
	}
	sum = h.Sum(nil)
	return
}
