package signr

import "os"

func CheckFileExists(name string) (fi os.FileInfo, exists bool, err error) {
	exists = true
	if fi, err = os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			exists = false
			err = nil
		}
	}
	return
}
