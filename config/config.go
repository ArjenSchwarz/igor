package config

import (
	"io/ioutil"
	"path/filepath"

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

func SetRawConfig(data []byte) {
	configFile = data
}

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
