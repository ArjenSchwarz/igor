package config

import (
	// "fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Service string `yaml:"service"`
}

type PluginConfig struct {
	Name string
}

// TODO ensure this can be used for plugin config
func ReadConfig() *Config {
	filename, _ := filepath.Abs("./config.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	config := new(Config)

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	return config
}
