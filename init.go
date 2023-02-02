package main

import (
	"fmt"
	"log"
	"os"
)

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
	} else {
		fmt.Printf("File does not exist\n")
	}

	// check if config file provided exists already
	// if yes pull the config and verify it.
	// 		confirm to the user if they want to ues the existing projects
	// 		if no existing projects in found config then walk through setting up the project
	// 		confirm to the user if they want to use the existing server configs
	// 		if no servers exist ask if the user wants to add server configurations
	// 		if yes walk through server configuration
	// if no confirm the name of the config file that the user wants to use
	//		walk through setup of projects
	//    walk through setup of servers

	// by the end here the user will have a valid configuration and have the file structure set up
}
