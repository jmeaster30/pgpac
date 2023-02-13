package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func prompt(message string, defaultValue string) string {
	fmt.Printf("%s Default [%s]: ", message, defaultValue)
	var value string
	fmt.Scanln(&value)
	if value == "" {
		return defaultValue
	}
	return value
}

func mkdir(filename string) {
	err := os.MkdirAll(filename, os.ModePerm)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}

func InitializeProject(configFilename string) {
	fmt.Println("Initializing Project")
	fmt.Printf("Supplied file: %s\n", configFilename)

	config := BlankPacConfig()

	if _, err := os.Stat(configFilename); err == nil {
		fmt.Printf("File exists\n")
		err = config.LoadConfig(configFilename)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", config)
	} else {
		fmt.Println("Config file does not exist yet!")

		selectedConfigFilename := prompt("What to name your config file?", configFilename)
		shouldSetupProject := strings.ToLower(prompt("Would you like to set up a project now?", "Yes"))
		if shouldSetupProject != "yes" && shouldSetupProject != "y" {
			fmt.Println("Not setting up projects yet. Projects can be setup by modifying the yaml file and running the init command again.")
			err := config.SaveConfig(selectedConfigFilename)
			if err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}
			return
		}

		projectName := prompt("Name of project.", "MyProject")
		useDefaultSubDirectoryNames := strings.ToLower(prompt("Would you like to use the default subdirectory names (i.e. <dir>/schema, <dir>/predeploy, <dir>/postdeploy, <dir>/seeds)?", "Yes"))
		if useDefaultSubDirectoryNames == "yes" || useDefaultSubDirectoryNames == "y" {
			projectFolder := prompt("Project Directory Name.", "myproject")
			mkdir(filepath.Join(projectFolder, "schema"))
			mkdir(filepath.Join(projectFolder, "predeploy"))
			mkdir(filepath.Join(projectFolder, "postdeploy"))
			mkdir(filepath.Join(projectFolder, "seeds"))
			config.Projects[projectName] = ProjectConfig{ProjectDirectory: projectFolder}
		} else {
			schemaDirectory := prompt("Schema Directory Name.", "project/schema")
			preDeployDirectory := prompt("Pre Deploy Directory Name.", "project/predeploy")
			postDeployDirectory := prompt("Post Deploy Directory Name.", "project/postdeploy")
			seedDirectory := prompt("Seed Directory Name.", "project/seeds")
			mkdir(schemaDirectory)
			mkdir(preDeployDirectory)
			mkdir(postDeployDirectory)
			mkdir(seedDirectory)
			config.Projects[projectName] = ProjectConfig{
				SchemaDirectory:     schemaDirectory,
				PreDeployDirectory:  preDeployDirectory,
				PostDeployDirectory: postDeployDirectory,
				SeedDirectory:       seedDirectory,
			}
		}
		err := config.SaveConfig(selectedConfigFilename)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
	}
}
