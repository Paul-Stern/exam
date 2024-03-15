package main

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Key  string `yaml:"key"`
		Cert string `yaml:"cert"`
	} `yaml:"server"`
	Rest struct {
		Host  string `yaml:"host"`
		Port  string `yaml:"port"`
		Rest  string `yaml:"restPath"`
		Nodes struct {
			GetQuestions    string `yaml:"getQuestions"`
			SaveTestResults string `yaml:"saveTestResult"`
			Register        string `yaml:"register"`
			User            string `yaml:"User"`
			Authenticate    string `yaml:"authenticate"`
		} `yaml:"nodes"`
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
