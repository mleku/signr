package cmd

import (
	"fmt"
	"github.com/mleku/appdata"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var dataDir, cfgFile string
var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "signr",
	Short: "A CLI generator, importer, signer, verifier and keychain for Nostr keys",
	Long: `A command line interface for generating, importing, signing, verifying and managing keys used with the Nostr protocol.

Designed to function in a similar way to ssh-keygen in that it keeps the keychain in a user directory with named key pairs and a configuration file.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// },
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
	cobra.OnInitialize(initConfig)
	dataDir = appdata.GetDataDir(rootCmd.Use, false)
	cfgFile = filepath.Join(dataDir, rootCmd.Use+".yml")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"prints more things")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {

		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)

	} else {

		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cmd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cmd")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
