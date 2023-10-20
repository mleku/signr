package signr

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (s *Signr) ReadFile(name string) (data []byte, err error) {

	path := filepath.Join(s.DataDir, name)
	// check the permissions are secure first
	var fi os.FileInfo
	if fi, err = os.Stat(path); err != nil {
		err = fmt.Errorf("error getting file info for %s: %v", name, err)
		return
	}
	// secret key files that are readable by other than the owner may not be
	// used
	if fi.Mode().Perm()&0077 != 0 &&
		!strings.HasSuffix(name, "."+PubExt) {
		err = fmt.Errorf("secret key '%s' has insecure permissions %s", name,
			fi.Mode().Perm())
		return
	}
	data, err = os.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("error reading file '%s': %v", name, err)
		return
	}
	return
}
