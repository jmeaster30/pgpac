package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// Functions for generating, loading, and validating configuration

type DeployConfig struct {
	OnFailure             string  `yaml:"onFailure,omitempty"`
	ColumnRenameTolerance float32 `yaml:"columnRenameTolerance,omitempty"`
	TableRenameTolerance  float32 `yaml:"tableRenameTolerance,omitempty"`
	RemoveUnusedTables    bool    `yaml:"removeUnusedTables,omitempty"`
}

type ProjectConfig struct {
	Deploy              DeployConfig `yaml:"deploy,omitempty"`
	SchemaDirectory     string       `yaml:"schemaDirectory,omitempty"`
	PreDeployDirectory  string       `yaml:"preDeployDirectory,omitempty"`
	PostDeployDirectory string       `yaml:"postDeployDirectory,omitempty"`
	SeedDirectory       string       `yaml:"seedDirectory,omitempty"`
	ProjectDirectory    string       `yaml:"projectDirectory,omitempty"`
}

type ServerConfig struct {
	Project    string `yaml:"project,omitempty"`
	Connection string `yaml:"connection,omitempty"`
	Hostname   string `yaml:"hostname,omitempty"`
	Port       string `yaml:"port,omitempty"`
	Database   string `yaml:"database,omitempty"`
	Username   string `yaml:"username,omitempty"`
}

type PacConfig struct {
	Projects map[string]ProjectConfig `yaml:"projects,omitempty"`
	Servers  map[string]ServerConfig  `yaml:"servers,omitempty"`
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

	return nil
}

func (p *PacConfig) SaveConfig(filename string) error {
	bytes, err := yaml.Marshal(p)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bytes, 0600)
}
