package main

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Rest struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Rest string `yaml:"restPath"`
	} `yaml:"rest"`
}

func readConf(cfg *Config) {
	f, err := os.Open("config.yaml")
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func baseUrl(cfg Config) string {
	return strings.Join([]string{
		"http://",
		cfg.Rest.Host,
		":",
		cfg.Rest.Port,
		cfg.Rest.Rest,
	}, "")
}
