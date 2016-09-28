package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

var configHolder Config

// Config contains general configuration details
type Config struct {
	Kms             bool
	Token           string
	DefaultLanguage string
	Blacklist       []string
	Whitelist       []string
	Languages       map[string]languageConfig
	LanguageDir     string
}

type languageConfig struct {
	Plugins  map[string]LanguagePluginDetails
	Language map[string]string
}

// LanguagePluginDetails holds the details for a plugin in a language
type LanguagePluginDetails struct {
	Description string
	Commands    map[string]LanguagePluginCommandDetails
}

// LanguagePluginCommandDetails holds the details for a plugin command in a language
type LanguagePluginCommandDetails struct {
	Description string
	Command     string
	Texts       map[string]string
}

var configFile []byte
var jsonConfig = true
var fallbackLanguage = "english.yml"

// GeneralConfig reads the configuration file and parses its general information
func GeneralConfig() (Config, error) {
	if configHolder.Token == "" {
		config := Config{}
		err := ParseConfig(&config)
		if err != nil {
			return config, err
		}
		config.Token, err = decryptValue(config.Kms, config.Token)
		if err != nil {
			return config, err
		}
		if config.LanguageDir == "" {
			config.LanguageDir = "language"
		}

		languageFiles, err := ioutil.ReadDir(config.LanguageDir)
		if err != nil {
			return config, err
		}
		languages := make(map[string]languageConfig)
		for _, file := range languageFiles {
			lConfig := languageConfig{}
			err = parseLanguageFile(config.LanguageDir+"/"+file.Name(), &lConfig)
			if err != nil {
				return config, err
			}
			languages[file.Name()] = lConfig
		}
		config.Languages = languages
		if config.DefaultLanguage == "" {
			config.DefaultLanguage = fallbackLanguage
		} else {
			config.DefaultLanguage = strings.Replace(config.DefaultLanguage, ".yml", "", -1) + ".yml"
		}
		configHolder = config
	}
	return configHolder, nil
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
		if err := yaml.Unmarshal(configFile, values); err != nil {
			return err
		}
	}
	return nil
}

// parseLanguageFile returns the requested language file
func parseLanguageFile(languagefile string, values interface{}) error {
	filename, _ := filepath.Abs(languagefile)
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(configFile, values); err != nil {
		return err
	}
	return nil
}
