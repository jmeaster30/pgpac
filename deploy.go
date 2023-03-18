package main

import (
	"os"
	"path/filepath"

	pg_query "github.com/pganalyze/pg_query_go/v2"
)

func Deploy(configFilename string, projectName string) {
	config := BlankPacConfig()
	if err := config.LoadConfig(configFilename); err != nil {
		LogError("error", "Failed to load config of name '%s'", configFilename)
		LogError("error", "%s", err)
		os.Exit(1)
	}

	if len(config.Projects) == 0 {
		LogError(config.Options.LogLevel, "No projects specified in loaded config file.")
		os.Exit(1)
	}

	if projectName == "" {
		LogError(config.Options.LogLevel, "No project supplied. Please specify a project to deploy.\n\tTODO: make this select the one if there is only 1 project")
		os.Exit(1)
	}

	project := config.Projects[projectName]
	if project.ProjectDirectory != "" {
		// we will be using the default folder names
		//predeployList := BuildFileList(filepath.Join(project.ProjectDirectory, "predeploy"))
		//postdeployList := BuildFileList(filepath.Join(project.ProjectDirectory, "postdeploy"))
		schemaList := BuildFileList(config, filepath.Join(project.ProjectDirectory, "schema"))
		//seedList := BuildFileList(filepath.Join(project.ProjectDirectory, "seeds"))
		LogDebug(config.Options.LogLevel, "%s", schemaList)
		BuildSchema(config, schemaList)
	}
}

func BuildFileList(config *PacConfig, foldername string) []string {
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
		LogError(config.Options.LogLevel, "%s", err)
		os.Exit(1)
	}
	return list
}

func readFile(config *PacConfig, filename string) string {
	dat, err := os.ReadFile(filename)
	if err != nil {
		LogError(config.Options.LogLevel, "Failed reading file '%s'", filename)
		os.Exit(1)
	}
	return string(dat)
}

func BuildSchema(config *PacConfig, files []string) {
	for _, val := range files {
		LogInfo(config.Options.LogLevel, "Processing '%s'...", val)
		content := readFile(config, val)
		tree, err := pg_query.Parse(content)
		if err != nil {
			LogError(config.Options.LogLevel, "Failed processing file '%s'...", val)
			LogError(config.Options.LogLevel, "%s", err)
			os.Exit(1)
		}

		LogInfo(config.Options.LogLevel, "Found %d statements", len(tree.Stmts))
		for _, val := range tree.Stmts {
			create_stmt := val.Stmt.GetCreateStmt()
			enum_stmt := val.Stmt.GetCreateEnumStmt()
			extension_stmt := val.Stmt.GetCreateExtensionStmt()

			if create_stmt == nil && enum_stmt == nil && extension_stmt == nil {
				LogWarning(config.Options.LogLevel, "Unexpected non-create statement in schema file. Ignoring...")
				LogWarning(config.Options.LogLevel, "%s", val.Stmt)
				continue
			}

			if create_stmt != nil {
				LogDebug(config.Options.LogLevel, "CREATE TABLE")
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

					for _, typeMod := range columnDefinition.TypeName.Typmods {
						column.ColumnType.TypeMods = append(column.ColumnType.TypeMods, int(typeMod.GetAConst().Val.GetInteger().Ival))
					}

					column.Constraint = ColumnConstraint{
						NotNull: false,
					}
					for _, constraint := range columnDefinition.Constraints {
						constraintNode := constraint.GetConstraint()
						BuildConstraint(config, &column.Constraint, constraintNode)
					}

					LogDebug(config.Options.LogLevel, "%v", columnDefinition)
					table.Columns = append(table.Columns, column)
				}

				//LogDebug("%v", create_stmt)
				LogDebug(config.Options.LogLevel, "Table Result\n%s", table.String())
			}

			if enum_stmt != nil {
				LogDebug(config.Options.LogLevel, "CREATE ENUM")
				enum := Enum{
					Name:   enum_stmt.TypeName[0].GetString_().Str,
					Values: []string{},
				}

				for _, val := range enum_stmt.Vals {
					enum.Values = append(enum.Values, val.GetString_().Str)
				}

				LogDebug(config.Options.LogLevel, "%v", enum_stmt)
				LogDebug(config.Options.LogLevel, "%v", enum)
			}

			if extension_stmt != nil {
				LogDebug(config.Options.LogLevel, "CREATE EXTENSION")
				ext := Extension{
					ExtensionName: extension_stmt.Extname,
				}

				LogDebug(config.Options.LogLevel, "%v", extension_stmt)
				LogDebug(config.Options.LogLevel, "%v", ext)
			}
		}

		//log.Printf("%s\n", tree.String())
		LogInfo(config.Options.LogLevel, "Done Processing '%s' :)", val)
	}
}

func BuildConstraint(config *PacConfig, columnConstraint *ColumnConstraint, constraintNode *pg_query.Constraint) {
	if constraintNode.Contype.String() == "CONSTR_FOREIGN" {
		columnConstraint.ForeignKey = Some(ForeignKeyConstraint{
			ReferencingTableName:  constraintNode.Pktable.Relname,
			ReferencingColumnName: constraintNode.PkAttrs[0].GetString_().Str, // FIXME can reference multiple columns
			MatchFull:             constraintNode.FkMatchtype == "f",
			MatchPartial:          constraintNode.FkMatchtype == "p",
			MatchSimple:           constraintNode.FkMatchtype == "s",
			OnDeleteAction:        constraintNode.FkDelAction,
			OnUpdateAction:        constraintNode.FkUpdAction,
		})
	} else if constraintNode.Contype.String() == "CONSTR_NOTNULL" {
		columnConstraint.NotNull = true
	} else {
		LogWarning(config.Options.LogLevel, "Unimplemented constraint node type '%s'", constraintNode.Contype.String())
	}
}
