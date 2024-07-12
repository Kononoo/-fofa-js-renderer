package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	ApiKey  string `yaml:"api_key"`
	MaxUrls int    `yaml:"max_urls"`
}

func loadConfig() (*Config, error) {
	configFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
