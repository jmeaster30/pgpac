package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// Functions for generating, loading, and validating configuration
type ProjectConfig struct {
	schema_directory   string
	one_shot_directory string
	seed_directory     string
	project_directory  string
}

type ServerConfig struct {
	project    string
	connection string
	hostname   string
	port       string
	database   string
	username   string
	password   string //should we even support storing this in the config? It may make sense but I want to leave more towards forcing the user to use this in a safe way
}

type PacConfig struct {
	projects map[string]ProjectConfig
	servers  map[string]ServerConfig
}

func BlankPacConfig() *PacConfig {
	return &PacConfig{
		projects: make(map[string]ProjectConfig),
		servers:  make(map[string]ServerConfig),
	}
}

func (p *PacConfig) LoadConfig(filename string) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		return err
	}

	for k, v := range data {
		fmt.Printf("%s -> %s\n", k, v)
	}

	return nil
}
