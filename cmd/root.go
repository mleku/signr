package cmd

import (
	"github.com/mleku/signr/pkg/signr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg *signr.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   signr.AppName,
	Short: "A CLI generator, importer, signer, verifier and keychain for Nostr keys",
	Long: `A command line interface for generating, importing, signing, verifying and managing keys used with the Nostr protocol.

Designed to function in a similar way to ssh-keygen in that it keeps the keychain in a user directory with named key as pairs of files and a configuration file.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		signr.Fatal("%s\n", err)
	}
}

func init() {

	cfg = signr.Init()
	rootCmd.PersistentFlags().BoolVarP(&cfg.Verbose,
		"verbose", "v", false, "prints more things")
	cobra.OnInitialize(initConfig(cfg))
}

// initConfig reads in config file and ENV variables if set.
func initConfig(cfg *signr.Config) func() {
	return func() {
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
}
