package signr

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/minio/sha256-simd"
)

// Hash accepts a filename, interpreting "-" to mean to read from stdin, and
// computes the sha256 hash of the file, or if the text is a 64 character long string that parses to hex, just decodes it and returns the decoded bytes.
func Hash(filename string) (sum []byte, wasHash bool, err error) {

	var exists bool
	var fi os.FileInfo
	fi, exists, err = CheckFileExists(filename)
	if err != nil {
		err = fmt.Errorf("error checking for file '%s': %s", filename, err)
		return
	}
	if exists && fi.IsDir() {
		err = fmt.Errorf("'%s' is a directory", filename)
		return
	}
	var f io.ReadCloser
	switch {
	case len(filename) == 64 && !exists:
		if sum, err = hex.DecodeString(filename); err == nil {
			wasHash = true
			return
		}
	case !exists:
		err = fmt.Errorf("file '%s' does not exist: %s", filename, err)
		return
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
