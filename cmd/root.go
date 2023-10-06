package cmd

import (
	"fmt"
	"github.com/mleku/appdata"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var dataDir, cfgFile, defaultKey string
var verbose bool

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

	dataDir = appdata.GetDataDir(rootCmd.Use, false)
	if _, err := os.Stat(dataDir); err != nil {
		if os.IsNotExist(err) {
			_, _ = fmt.Fprintf(os.Stderr,
				"creating signr data directory at '%s'\n", dataDir)
			if err = os.MkdirAll(dataDir, 0700); err != nil {
				_, _ = fmt.Fprintf(os.Stderr,
					"unable to create data dir, cannot proceed\n")
				os.Exit(1)
			}
		} else {
			panic(err)
		}
	}
	cfgFile = filepath.Join(dataDir, rootCmd.Use+".yaml")
	if _, err := os.Stat(dataDir); err != nil {
		if os.IsNotExist(err) {
			_, _ = fmt.Fprintf(os.Stderr,
				"creating signr data directory at '%s'\n", dataDir)
		} else {
			panic(err)
		}
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"prints more things")
	rootCmd.PersistentFlags().StringVarP(&defaultKey, "usekey", "u", "",
		"set secret key for signing by name")
	cobra.OnInitialize(initConfig)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dataDir)

	// read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:",
			viper.ConfigFileUsed())
	}
	defaultKey = viper.GetString("default")
}
