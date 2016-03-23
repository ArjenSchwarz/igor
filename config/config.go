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
var jsonConfig bool

// GeneralConfig reads the configuration file and parses its general information
func GeneralConfig() Config {
	configFile := getRawConfig()
	config := Config{}
	if jsonConfig {
		err := json.Unmarshal(configFile, &config)
		if err != nil {
			panic(err)
		}
	} else {
		err := yaml.Unmarshal(configFile, &config)
		if err != nil {
			panic(err)
		}
	}
	return config
}

// SetRawConfig allows you to set the raw, unparsed, configuration data.
func SetRawConfig(data []byte) {
	configFile = data
}

// getRawConfig allows you to retrieve the raw, unparsed, configuration data.
// If no configuration is present, it will pull this from the configuration
// file.
func getRawConfig() []byte {
	if len(configFile) == 0 {
		configFile = getConfigFile()
	}
	return configFile
}

// GetConfigFile retrieves the contents of the config file as a byte array
func getConfigFile() []byte {
	envConf := os.Getenv("IGOR_CONFIG")
	if envConf != "" {
		jsonConfig = true
		return []byte(envConf)
	}
	filename, _ := filepath.Abs("./config.json")
	if _, err := os.Stat(filename); err == nil {
		jsonConfig = true
	} else {
		jsonConfig = false
		filename, _ = filepath.Abs("./config.yml")
	}
	configFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}
	return configFile
}

// ParsePluginConfig parses the plugin file and unmarshals it into the
// provided interface
func ParsePluginConfig(values interface{}) error {
	configFile := getRawConfig()

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
