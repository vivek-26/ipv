// Package config is used to generate config file for the user
package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/manifoldco/promptui"

	"github.com/vivek-26/ipv/reporter"
)

const configFileName = "/.config.toml"

// Prompt template
var templates = &promptui.PromptTemplates{
	Prompt:  "{{ . }} ",
	Valid:   "{{ . | green }} ",
	Invalid: "{{ . | red }} ",
	Success: "{{ . | bold }} ",
}

// Supported protocols
var protocols = []string{"udp", "tcp"}

// cfg has fields for user configuration
type cfg struct {
	Username    string `toml:"username"`
	CountryCode string `toml:"countryCode"` // 2 letter country code
	Protocol    string `toml:"protocol"`
}

// Implement stringer interface
func (c *cfg) String() string {
	return fmt.Sprintf("Username: %v, Country: %v, Protocol: %v",
		c.Username, c.CountryCode, c.Protocol)
}

// Generate creates a new config file
func Generate(configDirPath string) {
	c := &cfg{
		Username:    getUsername(),
		CountryCode: getCountryCode(),
		Protocol:    getProtocol(),
	}

	// Encode `cfg` struct as toml
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(c); err != nil {
		reporter.Error(err)
		os.Exit(1)
	}

	// Check if config directory exists; if it doesn't exist, create it
	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
		if err := os.Mkdir(configDirPath, os.ModePerm); err != nil {
			reporter.Error(err)
			os.Exit(1)
		}
	}

	// Write config file
	configFile := configDirPath + configFileName
	if err := ioutil.WriteFile(configFile, buf.Bytes(), 0400); err != nil {
		reporter.Error(err)
		os.Exit(1)
	}

	reporter.Info("Saved configuration to " + configFile)
}

func getUsername() string {
	prompt := promptui.Prompt{
		Label:     "Username:",
		Templates: templates,
	}

	uname, err := prompt.Run()
	if err != nil {
		reporter.Error(err)
		os.Exit(1)
	}

	return uname
}

func getCountryCode() string {
	prompt := promptui.Prompt{
		Label:     "Country Code:",
		Templates: templates,
	}

	countryCode, err := prompt.Run()
	if err != nil {
		reporter.Error(err)
		os.Exit(1)
	}

	return countryCode
}

func getProtocol() string {
	// Select template
	selectTemplates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "> {{ . | green }}",
		Inactive: "  {{ . | cyan }}",
		Selected: "Protocol: {{ . }}",
	}

	prompt := promptui.Select{
		Label:     "Protocol",
		Items:     protocols,
		Templates: selectTemplates,
		Size:      2,
	}

	_, protocol, err := prompt.Run()
	if err != nil {
		reporter.Error(err)
		os.Exit(1)
	}

	return protocol
}
