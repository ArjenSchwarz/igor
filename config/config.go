package config

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Config contains general configuration details
type Config struct {
	Token     string
	Blacklist []string
	Whitelist []string
}

var configFile []byte
var jsonConfig = true

// GeneralConfig reads the configuration file and parses its general information
func GeneralConfig() (Config, error) {
	config := Config{}
	err := ParseConfig(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}

// GetConfigFile retrieves the contents of the config file as a byte array
func getConfigFile() ([]byte, error) {
	if len(configFile) != 0 {
		return configFile, nil
	}
	envConf := os.Getenv("IGOR_CONFIG")
	if envConf != "" {
		return []byte(envConf), nil
	}
	filename, _ := filepath.Abs("./config.json")
	if _, err := os.Stat(filename); err != nil {
		jsonConfig = false
		filename, _ = filepath.Abs("./config.yml")
	}
	configFile, err := ioutil.ReadFile(filename)

	if err != nil {
		return configFile, err
	}
	return configFile, nil
}

// ParseConfig parses the config file and unmarshals it into the
// provided interface
func ParseConfig(values interface{}) error {
	configFile, err := getConfigFile()
	if err != nil {
		return err
	}

	if jsonConfig {
		if err := json.Unmarshal(configFile, &values); err != nil {
			return err
		}
	} else {
		if err := yaml.Unmarshal(configFile, &values); err != nil {
			return err
		}
	}
	return nil
}
