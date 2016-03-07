package config

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Token string
}

func ReadConfig() Config {
	configFile := GetConfigFile()
	config := Config{}

	err := yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}
	return config
}

func GetConfigFile() []byte {
	filename, _ := filepath.Abs("./config.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}
	return yamlFile
}
