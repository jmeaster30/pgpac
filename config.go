package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// Functions for generating, loading, and validating configuration
type ProjectConfig struct {
	SchemaDirectory  string `yaml:"schemaDirectory"`
	OneShotDirectory string `yaml:"oneShotDirectory"`
	SeedDirectory    string `yaml:"seedDirectory"`
	ProjectDirectory string `yaml:"projectDirectory"`
}

type ServerConfig struct {
	Project    string
	Connection string
	Hostname   string
	Port       string
	Database   string
	Username   string
	Password   string //should we even support storing this in the config? It may make sense but I want to leave more towards forcing the user to use this in a safe way
}

type PacConfig struct {
	Projects map[string]ProjectConfig
	Servers  map[string]ServerConfig
}

func BlankPacConfig() *PacConfig {
	return &PacConfig{
		Projects: make(map[string]ProjectConfig),
		Servers:  make(map[string]ServerConfig),
	}
}

func (p *PacConfig) LoadConfig(filename string) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, p)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", p)

	return nil
}
