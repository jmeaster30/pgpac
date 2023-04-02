package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// Functions for generating, loading, and validating configuration

type DeployConfig struct {
	OnFailure             string  `yaml:"onFailure,omitempty"`
	ColumnRenameTolerance float32 `yaml:"columnRenameTolerance,omitempty"`
	TableRenameTolerance  float32 `yaml:"tableRenameTolerance,omitempty"`
	RemoveUnusedTables    bool    `yaml:"removeUnusedTables,omitempty"`
	SchemaName            string  `yaml:"schemaName,omitempty"`
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

type OptionsConfig struct {
	LogLevel string `yaml:"logLevel,omitempty" validOptions:"debug,warn,info,error"`
}

type PacConfig struct {
	Projects map[string]ProjectConfig `yaml:"projects,omitempty"`
	Servers  map[string]ServerConfig  `yaml:"servers,omitempty"`
	Options  OptionsConfig            `yaml:"options,omitempty"`
}

func BlankPacConfig() *PacConfig {
	return &PacConfig{
		Projects: make(map[string]ProjectConfig),
		Servers:  make(map[string]ServerConfig),
		Options: OptionsConfig{
			LogLevel: "debug",
		},
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

	err = p.validateConfig()
	if err != nil {
		return err
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (p *PacConfig) validateConfig() error {
	errorMessages := []string{}

	// FIXME: make it so this goes through the whole pac config object
	optionsType := reflect.TypeOf(p.Options)
	for i := 0; i < optionsType.NumField(); i++ {
		field := optionsType.Field(i)
		value := reflect.ValueOf(p.Options)
		validOptions := strings.Split(field.Tag.Get("validOptions"), ",")
		if len(validOptions) != 0 && !contains(validOptions, strings.ToLower(value.FieldByName(field.Name).String())) {
			errorMessages = append(errorMessages, fmt.Sprintf("Expected %s to be one of %s but found '%s'.", field.Name, strings.Join(validOptions, ", "), p.Options.LogLevel))
		}
	}

	if len(errorMessages) == 0 {
		return nil
	}
	//lint:ignore ST1005 I want this to be capitalized
	return fmt.Errorf("There were %d error(s) in the config :(\n%s", len(errorMessages), strings.Join(errorMessages, "\n"))
}

func (p *PacConfig) SaveConfig(filename string) error {
	bytes, err := yaml.Marshal(p)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bytes, 0600)
}
