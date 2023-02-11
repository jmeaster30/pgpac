package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	pg_query "github.com/pganalyze/pg_query_go/v2"
)

func Deploy(configFilename string, projectName string) {
	config := BlankPacConfig()
	if err := config.LoadConfig(configFilename); err != nil {
		log.Fatalf("Failed to load config of name '%s'\n%s\n", configFilename, err)
	}

	if len(config.Projects) == 0 {
		log.Fatalf("No projects specified in loaded config file.\n")
	}

	if projectName == "" {
		log.Fatalf("No project supplied. Please specify a project to deploy.\nTODO: make this select the one if there is only 1 project")
	}

	project := config.Projects[projectName]
	if project.SchemaDirectory != "" {
		// we will be using the default folder names
		fmt.Println("here")
		results := BuildFileList(filepath.Join(project.SchemaDirectory, "schema"))
		fmt.Printf("%s\n", results)
		BuildSchema(results)
	}
}

func BuildFileList(foldername string) []string {
	list := []string{}
	err := filepath.Walk(foldername,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				list = append(list, path)
			}
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}
	return list
}

func readFile(filename string) string {
	dat, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed reading file '%s'\n", filename)
	}
	return string(dat)
}

func BuildSchema(files []string) {
	for _, val := range files {
		log.Printf("Processing '%s'...\n", val)
		content := readFile(val)
		tree, err := pg_query.Parse(content)
		if err != nil {
			log.Fatalf("Failed processing file '%s'...\n", val)
		}
		log.Printf("%s\n", tree.String())
		log.Printf("Done Processing '%s' :)\n", val)
	}
}
