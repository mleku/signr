package signr

import (
	"fmt"
	"os"
	"path/filepath"
)

func (s *Signr) Delete(name string) (err error) {

	var grid [][]string
	grid, _, err = s.GetList(nil)
	for i := range grid {

		if grid[i][0] == name {

			// generate random nonce to ensure no collisions are possible (ie,
			// extremely unlikely)
			var nonce string
			nonce, err = s.GetNonceHex()

			// construct deleted names for keys
			sk := name + "." + nonce + "." + DeletedExt
			pk := name + "." + nonce + "." + PubExt + "." + DeletedExt

			// add path prefixes to old and new filenames
			origSk := filepath.Join(s.DataDir, name)
			origPk := filepath.Join(s.DataDir, name+"."+PubExt)
			newSk := filepath.Join(s.DataDir, sk)
			newPk := filepath.Join(s.DataDir, pk)

			// rename all the things to deleted status
			if err = os.Rename(origSk, newSk); err != nil {
				return
			}
			if err = os.Rename(origPk, newPk); err != nil {
				return
			}

			// success!
			return
		}
	}

	// key name was not found
	err = fmt.Errorf("key named '%s' not found in keychain; "+
		"use signr listkeys to see list of keys in keychain", name)

	return
}
