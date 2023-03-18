package main

import (
	"os"
	"path/filepath"

	pg_query "github.com/pganalyze/pg_query_go/v2"
)

func Deploy(configFilename string, projectName string) {
	config := BlankPacConfig()
	if err := config.LoadConfig(configFilename); err != nil {
		LogError("Failed to load config of name '%s'", configFilename)
		LogError("%s", err)
		os.Exit(1)
	}

	if len(config.Projects) == 0 {
		LogError("No projects specified in loaded config file.")
		os.Exit(1)
	}

	if projectName == "" {
		LogError("No project supplied. Please specify a project to deploy.\n\tTODO: make this select the one if there is only 1 project")
		os.Exit(1)
	}

	project := config.Projects[projectName]
	if project.ProjectDirectory != "" {
		// we will be using the default folder names
		//predeployList := BuildFileList(filepath.Join(project.ProjectDirectory, "predeploy"))
		//postdeployList := BuildFileList(filepath.Join(project.ProjectDirectory, "postdeploy"))
		schemaList := BuildFileList(filepath.Join(project.ProjectDirectory, "schema"))
		//seedList := BuildFileList(filepath.Join(project.ProjectDirectory, "seeds"))
		LogDebug("%s", schemaList)
		BuildSchema(schemaList)
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
		LogError("%s", err)
		os.Exit(1)
	}
	return list
}

func readFile(filename string) string {
	dat, err := os.ReadFile(filename)
	if err != nil {
		LogError("Failed reading file '%s'", filename)
		os.Exit(1)
	}
	return string(dat)
}

func BuildSchema(files []string) {
	for _, val := range files {
		LogInfo("Processing '%s'...", val)
		content := readFile(val)
		tree, err := pg_query.Parse(content)
		if err != nil {
			LogError("Failed processing file '%s'...", val)
			LogError("%s", err)
			os.Exit(1)
		}

		LogInfo("Found %d statements", len(tree.Stmts))
		for _, val := range tree.Stmts {
			create_stmt := val.Stmt.GetCreateStmt()
			enum_stmt := val.Stmt.GetCreateEnumStmt()
			extension_stmt := val.Stmt.GetCreateExtensionStmt()

			if create_stmt == nil && enum_stmt == nil && extension_stmt == nil {
				LogWarning("Unexpected non-create statement in schema file. Ignoring...")
				LogWarning("%s", val.Stmt)
				continue
			}

			if create_stmt != nil {
				LogDebug("CREATE TABLE")
				table := Table{
					TableName: create_stmt.Relation.Relname,
				}

				for _, val := range create_stmt.TableElts {
					columnDefinition := val.GetColumnDef()

					column := Column{
						Name: columnDefinition.Colname,
					}

					for _, typeName := range columnDefinition.TypeName.Names {
						column.ColumnType.TypePath = append(column.ColumnType.TypePath, typeName.GetString_().Str)
					}

					LogDebug("%v", columnDefinition)
					table.Columns = append(table.Columns, column)
				}

				LogDebug("%v", create_stmt)
				LogDebug("%s", table)
			}

			if enum_stmt != nil {
				LogDebug("CREATE ENUM")
				enum := Enum{
					Name:   enum_stmt.TypeName[0].GetString_().Str,
					Values: []string{},
				}

				for _, val := range enum_stmt.Vals {
					enum.Values = append(enum.Values, val.GetString_().Str)
				}

				LogDebug("%v", enum_stmt)
				LogDebug("%v", enum)
			}

			if extension_stmt != nil {
				LogDebug("CREATE EXTENSION")
				ext := Extension{
					ExtensionName: extension_stmt.Extname,
				}

				LogDebug("%v", extension_stmt)
				LogDebug("%v", ext)
			}
		}

		//log.Printf("%s\n", tree.String())
		LogInfo("Done Processing '%s' :)", val)
	}
}
