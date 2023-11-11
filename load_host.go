package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Host []string `yaml:"host"`
}

func LoadHost() ([]string, error) {

	fileName, _ := filepath.Abs("mikrotik.yaml")
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config.Host, nil
}
