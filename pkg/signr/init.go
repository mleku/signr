package signr

import (
	"fmt"
	"os"
	"runtime"

	"mleku.online/git/appdata"
)

// Signr stores the configuration for signr.
type Signr struct {
	DataDir       string
	CfgFile       string
	DefaultKey    string
	Verbose       bool
	Color         bool
	PassEntryType int
}

// Init sets up the data directory if it doesn't exist, checks the permissions
// of the directory and configuration file.
//
// Applications consuming this library can use alternative password input
// methods or use their own and pass the value via environment variables when
// calling the CLI.
func Init(passEntryType int) (s *Signr, err error) {

	s = &Signr{PassEntryType: passEntryType}
	s.DataDir = appdata.GetDataDir(runtime.GOOS, AppName, false)
	fi, exists, err := CheckFileExists(s.DataDir)
	if err != nil {
		err = fmt.Errorf("error checking if datadir exists: %s", err)
		return
	}
	s.CfgFile = s.GetCfgFilename()
	if !exists {
		s.Err("First run: Creating signr data directory at %s\n\n",
			s.DataDir)
		if err = os.MkdirAll(s.DataDir, DataDirPerm); err != nil {
			err = fmt.Errorf("unable to create data dir, cannot proceed: %s",
				err)
			return
		}
		// Touch the config file so it is ready to write to.
		if err = os.WriteFile(s.CfgFile, []byte{}, ConfigFilePerm); err != nil {
			err = fmt.Errorf("error writing config file '%s': %s", s.CfgFile,
				err)
			return
		}
		// check the permissions
	} else if fi.Mode().Perm()&DataFileMask != 0 {
		err = fmt.Errorf(
			"data directory %s has insecure permissions %s"+
				" recommended to restore it to %s (0%o), "+
				"and investigate how it got changed",
			s.DataDir, fi.Mode().Perm(),
			DataDirPerm, DataDirPerm)
		return
	}
	if fi, err = os.Stat(s.CfgFile); err != nil {
		err = fmt.Errorf("unexpected error probing config file %s: %s",
			s.CfgFile, err)
		return
	}
	// check configuration permissions
	if fi.Mode().Perm()&DataFileMask != 0 {
		err = fmt.Errorf(
			"configuration file %s has insecure permissions %s (0%o)"+
				" recommended to restore it to %s (0%o), "+
				"and investigate how it got changed",
			s.CfgFile, fi.Mode().Perm(), fi.Mode().Perm(),
			ConfigFilePerm, ConfigFilePerm)
		return
	}
	return
}
