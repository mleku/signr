package signr

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"mleku.online/git/signr/pkg/signr"
)

var s *signr.Signr

const customInfo = `The custom namespace string is sanitized, removing all non-printable/whitespace and removing any number of these characters found sequentially into a single hyphen "-" with all leading and following whitespace removed. The only limitation on the rest of the content of this parameter is that the remainder of the characters are designated as "printable" in the UTF-8 standard, this includes most punctuation and foreign language characters.
`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   signr.AppName,
	Short: "signr - A CLI key generator, importer, signer, verifier and keychain for Nostr keys",
	Long: `signr

A command line interface for generating, importing, signing, verifying and managing keys used with the Nostr protocol.

Designed to function in a similar way to ssh-keygen in that it keeps the keychain in a user directory with named key as pairs of files and a configuration file.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		s.Fatal("%s\n", err)
	}
}

var Verbose, Color bool

func init() {

	var err error
	// because this is a CLI app we know the user can enter passwords this way.
	// other types of apps using this can load the environment variables.
	s, err = signr.Init(signr.PasswordEntryViaTTY)
	if err != nil {
		s.Fatal("fatal error: %s\n", err)
	}
	rootCmd.PersistentFlags().BoolVarP(&Verbose,
		"verbose", "v", false, "prints more things")
	rootCmd.PersistentFlags().BoolVarP(&Color,
		"color", "c", false, "prints color things")
	cobra.OnInitialize(initConfig(s))
}

// initConfig reads in config file and ENV variables if set.
func initConfig(cfg *signr.Signr) func() {
	return func() {
		viper.SetConfigName(signr.ConfigName)
		viper.SetConfigType(signr.ConfigExt)
		viper.AddConfigPath(cfg.DataDir)
		// read in environment variables that match
		viper.SetEnvPrefix(signr.AppName)
		viper.AutomaticEnv()
		s.Verbose.Store(Verbose)
		s.Color.Store(Color)
		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil && cfg.Verbose.Load() {
			cfg.Log("Using config file: %s\n", viper.ConfigFileUsed())
		}

		// if pass is given on CLI it overrides environment, but if it is empty and environment has a value, load it
		if Pass == "" {
			if p := viper.GetString("pass"); p != "" {
				Pass = p
			}
		}
		cfg.DefaultKey = viper.GetString("default")
	}
}
