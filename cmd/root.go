package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "ipv",
	Short: "IPVanish CLI utility",
	Long: `  IPVanish is a VPN provider
  This command lists the servers and connects to the
  selected server in a particular country.
  Complete documentation is available at http://ipvanish.com/.`,
	Version: "0.1",
	Run:     func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Config dir
	configDir := filepath.Join(home, ".ipv")

	// Tell viper to look for `.config.toml` in configuration folder
	viper.AddConfigPath(configDir)
	viper.SetConfigType("toml")
	viper.SetConfigName(".config")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Cannot find configuration file")
		}
		if _, ok := err.(viper.UnsupportedConfigError); ok {
			fmt.Println("Unsupported config file type, expected toml")
		}
	}
}
