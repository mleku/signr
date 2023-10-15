package signr

import (
	"fmt"
	"github.com/mleku/appdata"
	"os"
)

// Signr stores the configuration for signr.
type Signr struct {
	DataDir    string
	CfgFile    string
	DefaultKey string
	Verbose    bool
	Color      bool
}

// Init sets up the data directory if it doesn't exist, checks the permissions
// of the directory and configuration file.
func Init() (s *Signr) {

	s = &Signr{}

	s.DataDir = appdata.GetDataDir(AppName, false)

	fi, exists, err := CheckFileExists(s.DataDir)

	s.CfgFile = s.GetCfgFilename()

	if !exists {
		s.PrintErr("First run: Creating signr data directory at %s\n\n",
			s.DataDir)

		if err = os.MkdirAll(s.DataDir, DataDirPerm); err != nil {
			s.Fatal("unable to create data dir, cannot proceed\n")
		}

		// Touch the config file so it is ready to write to.
		os.WriteFile(s.CfgFile, []byte{}, ConfigFilePerm)

		// check the permissions
	} else if fi.Mode().Perm()&DataFileMask != 0 {

		err = fmt.Errorf(
			"data directory %s has insecure permissions %s"+
				" recommended to restore it to %s (0%o), "+
				"and investigate how it got changed",
			s.DataDir, fi.Mode().Perm(),
			DataDirPerm, DataDirPerm)

		s.Fatal("ERROR: '%s'\n", err)
	}

	if fi, err = os.Stat(s.CfgFile); err != nil {
		s.Fatal("Unexpected error probing config file %s: '%s'\n",
			s.CfgFile, err)
	}

	// check configuration permissions
	if fi.Mode().Perm()&DataFileMask != 0 {

		err = fmt.Errorf(
			"configuration file %s has insecure permissions %s (0%o)"+
				" recommended to restore it to %s (0%o), "+
				"and investigate how it got changed",
			s.CfgFile, fi.Mode().Perm(), fi.Mode().Perm(),
			ConfigFilePerm, ConfigFilePerm)

		s.Fatal("ERROR: %s\n", err)
	}

	return s
}
