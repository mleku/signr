package cmd

import (
	"fmt"
	"github.com/mleku/appdata"
	"github.com/mleku/signr/pkg/signr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var cfg signr.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "signr",
	Short: "A CLI generator, importer, signer, verifier and keychain for Nostr keys",
	Long: `A command line interface for generating, importing, signing, verifying and managing keys used with the Nostr protocol.

Designed to function in a similar way to ssh-keygen in that it keeps the keychain in a user directory with named key pairs and a configuration file.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	cfg.DataDir = appdata.GetDataDir(rootCmd.Use, false)

	if fi, err := os.Stat(cfg.DataDir); err != nil {

		if os.IsNotExist(err) {

			signr.PrintErr(
				"creating signr data directory at '%s'\n", cfg.DataDir)
			if err = os.MkdirAll(cfg.DataDir, 0700); err != nil {
				signr.Fatal("unable to create data dir, cannot proceed\n")
			}

		} else {

			signr.PrintErr("%s\n", err)
			os.Exit(1)
		}

	} else {

		// check the permissions
		if fi.Mode().Perm()&0077 != 0 {

			err = fmt.Errorf(
				"data directory '%s' has insecure permissions %s"+
					" recommended to restore it to -rwx------ (0700), "+
					"and investigate how it got changed",
				cfg.DataDir, fi.Mode().Perm())

			signr.Fatal("ERROR: '%s'\n", err)
		}
	}

	cfg.CfgFile = filepath.Join(cfg.DataDir, rootCmd.Use+"."+signr.ConfigExt)

	if _, err := os.Stat(cfg.DataDir); err != nil {

		if os.IsNotExist(err) {

			signr.PrintErr(
				"creating signr data directory at '%s'\n", cfg.DataDir)

		} else {

			signr.Fatal("%s\n", err)
		}
	}

	rootCmd.PersistentFlags().
		BoolVarP(&cfg.Verbose, "Verbose", "v", false, "prints more things")

	cobra.OnInitialize(initConfig)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.SetConfigName(signr.ConfigName)
	viper.SetConfigType(signr.ConfigExt)
	viper.AddConfigPath(cfg.DataDir)

	// read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && cfg.Verbose {

		signr.PrintErr("Using config file: %s\n", viper.ConfigFileUsed())
	}

	cfg.DefaultKey = viper.GetString("default")
}
