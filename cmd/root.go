package cmd

import (
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vivek-26/ipv/config"
	"github.com/vivek-26/ipv/reporter"
)

// configDirName is the directory name where user config will be stored
const configDirName = ".ipv"

// rootCmd represents the base command when called without any subcommands.
func rootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ipv",
		Short: "IPVanish CLI utility",
		Long: `  IPVanish is a VPN provider
		This command lists the servers and connects to the
		selected server in a particular country.
		Complete documentation is available at http://ipvanish.com/.`,
		Version: "0.1",
		Run:     func(cmd *cobra.Command, args []string) {},
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// Read config file or create if it does not exist
	cobra.OnInitialize(initConfig)

	// Add child commands
	rootCommand := rootCmd()
	connectCommand := connectCmd()
	rootCommand.AddCommand(connectCommand)

	if err := rootCommand.Execute(); err != nil {
		reporter.Error(err)
	}
}

// initConfig reads in config file.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		reporter.Error(err)
	}

	// Config directory path
	configDirPath := filepath.Join(home, configDirName)

	// Tell viper to look for `.config.toml` in configuration folder
	viper.AddConfigPath(configDirPath)
	viper.SetConfigType("toml")
	viper.SetConfigName(".config")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			reporter.Warn("Cannot find configuration file, generating new one...")
			// Generate new config file
			config.Generate(configDirPath)
		}
		if _, ok := err.(viper.UnsupportedConfigError); ok {
			reporter.Error("Unsupported config file type, expected toml")
		}
	}
}
