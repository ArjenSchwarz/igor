package config

import (
	"io/ioutil"
	"path/filepath"
// 	"errors"

	"gopkg.in/yaml.v2"
)

// Config contains general configuration details
type Config struct {
	Token     string
	Blacklist []string
	Whitelist []string
}

var configFile []byte

// GeneralConfig reads the configuration file and parses its general information
func GeneralConfig() Config {
	configFile := GetRawConfig()
	config := Config{}

	err := yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}
	return config
}

// SetRawConfig allows you to set the raw, unparsed, configuration data.
func SetRawConfig(data []byte) {
	configFile = data
}

// GetRawConfig allows you to retrieve the raw, unparsed, configuration data.
// If no configuration is present, it will pull this from the configuration
// file.
func GetRawConfig() []byte {
	if len(configFile) == 0 {
		configFile = getConfigFile()
	}
	return configFile
}

// GetConfigFile retrieves the contents of the config file as a byte array
func getConfigFile() []byte {
	filename, _ := filepath.Abs("./config.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}
	return yamlFile
}

func ParsePluginConfig(values interface{}) error {
    configFile := GetRawConfig()

	err := yaml.Unmarshal(configFile, values)
	return err
}