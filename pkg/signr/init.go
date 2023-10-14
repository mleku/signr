package signr

import (
	"fmt"
	"github.com/mleku/appdata"
	"github.com/spf13/viper"
	"os"
)

// Config stores the configuration for signr.
type Config struct {
	DataDir    string
	CfgFile    string
	DefaultKey string
	Verbose    bool
}

// Init sets up the data directory if it doesn't exist, checks the permissions
// of the directory and configuration file.
func Init() (cfg *Config) {

	cfg = &Config{}

	cfg.DataDir = appdata.GetDataDir(AppName, false)

	fi, exists, err := CheckFileExists(cfg.DataDir)

	cfg.CfgFile = cfg.GetCfgFilename()

	if !exists {
		PrintErr("First run: Creating signr data directory at %s\n\n",
			cfg.DataDir)

		if err = os.MkdirAll(cfg.DataDir, DataDirPerm); err != nil {
			Fatal("unable to create data dir, cannot proceed\n")
		}

		// Touch the config file so it is ready to write to.
		os.WriteFile(cfg.CfgFile, []byte{}, ConfigFilePerm)

		// check the permissions
	} else if fi.Mode().Perm()&DataFileMask != 0 {

		err = fmt.Errorf(
			"data directory %s has insecure permissions %s"+
				" recommended to restore it to %s (0%o), "+
				"and investigate how it got changed",
			cfg.DataDir, fi.Mode().Perm(),
			DataDirPerm, DataDirPerm)

		Fatal("ERROR: '%s'\n", err)
	}

	if fi, err = os.Stat(cfg.CfgFile); err != nil {
		Fatal("Unexpected error probing config file %s: '%s'\n",
			cfg.CfgFile, err)
	}

	// check configuration permissions
	if fi.Mode().Perm()&DataFileMask != 0 {

		err = fmt.Errorf(
			"configuration file %s has insecure permissions %s (0%o)"+
				" recommended to restore it to %s (0%o), "+
				"and investigate how it got changed",
			cfg.CfgFile, fi.Mode().Perm(), fi.Mode().Perm(),
			ConfigFilePerm, ConfigFilePerm)

		Fatal("ERROR: %s\n", err)
	}

	cfg.DefaultKey = viper.GetString("default")

	return cfg
}

