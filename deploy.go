package main

import (
	"os"
	"path"
	"path/filepath"

	pg_query "github.com/pganalyze/pg_query_go/v2"
)

func Deploy(configFilename string, projectName string) {
	config := BlankPacConfig()
	configFolder := path.Dir(configFilename)
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

	project, prs := config.Projects[projectName]
	if !prs {
		LogError("error", "Project '%s' was not found in the config file '%s'.", projectName, configFilename)
	}
	if project.ProjectDirectory != "" {
		// we will be using the default folder names
		//predeployList := BuildFileList(filepath.Join(configFolder, project.ProjectDirectory, "predeploy"))
		//postdeployList := BuildFileList(filepath.Join(configFolder, project.ProjectDirectory, "postdeploy"))
		schemaList := BuildFileList(config, filepath.Join(configFolder, project.ProjectDirectory, "schema"))
		//seedList := BuildFileList(filepath.Join(configFolder, project.ProjectDirectory, "seeds"))
		LogDebug(config.Options.LogLevel, "%s", schemaList)
		_ = BuildSchema(config, schemaList)

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

func BuildSchema(config *PacConfig, files []string) []SchemaObject {
	results := []SchemaObject{}

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
			if create_stmt != nil {
				table := Table{
					TableName: create_stmt.Relation.Relname,
				}

				for _, val := range create_stmt.TableElts {
					columnDefinition := val.GetColumnDef()
					if columnDefinition != nil {
						column := Column{
							Name:      columnDefinition.Colname,
							TableName: table.TableName,
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
							BuildColumnConstraint(config, &column.Constraint, constraintNode)
						}

						table.Columns = append(table.Columns, column)
						results = append(results, column)
						continue
					}

					tableConstraint := val.GetConstraint()
					if tableConstraint != nil {
						//LogDebug(config.Options.LogLevel, "Found table constraint: %+v", tableConstraint)
						BuildTableConstraint(config, &table.TableConstraint, tableConstraint)
						continue
					}

					LogWarning(config.Options.LogLevel, "IGNORING UNSUPPORTED TABLE ELEMENT: %+v", val)
					continue
				}

				LogDebug(config.Options.LogLevel, "Table Result\n%s", table.String())
				results = append(results, table)
				continue
			}

			index_stmt := val.Stmt.GetIndexStmt()
			if index_stmt != nil {
				index := Index{
					Name:        index_stmt.Idxname,
					IndexType:   index_stmt.AccessMethod,
					TableName:   index_stmt.Relation.Relname,
					IndexParams: []IndexParam{},
				}

				// TODO add asserts for the stuff we haven't implemented

				for _, param := range index_stmt.IndexParams {
					index_elem := param.GetIndexElem()
					index.IndexParams = append(index.IndexParams, IndexParam{
						ColumnName: index_elem.Name,
					})
				}

				LogDebug(config.Options.LogLevel, "%+v", index)
				results = append(results, index)
				continue
			}

			enum_stmt := val.Stmt.GetCreateEnumStmt()
			if enum_stmt != nil {
				enum := Enum{
					Name:   enum_stmt.TypeName[0].GetString_().Str,
					Values: []string{},
				}

				for _, val := range enum_stmt.Vals {
					enum.Values = append(enum.Values, val.GetString_().Str)
				}

				LogDebug(config.Options.LogLevel, "%v", enum)
				results = append(results, enum)
				continue
			}

			extension_stmt := val.Stmt.GetCreateExtensionStmt()
			if extension_stmt != nil {
				ext := Extension{
					ExtensionName: extension_stmt.Extname,
				}

				LogDebug(config.Options.LogLevel, "%v", ext)
				results = append(results, ext)
				continue
			}

			view_stmt := val.Stmt.GetViewStmt()
			if view_stmt != nil {
				LogWarning(config.Options.LogLevel, "FOUND VIEW - Unimplemented")
				continue
			}

			LogWarning(config.Options.LogLevel, "Unexpected non-create statement in schema file. Ignoring...")
			LogWarning(config.Options.LogLevel, "%s", val.Stmt)
		}

		LogInfo(config.Options.LogLevel, "Done Processing '%s' :)", val)
	}

	return results
}

func BuildColumnConstraint(config *PacConfig, columnConstraint *ColumnConstraint, constraintNode *pg_query.Constraint) {
	if constraintNode.Contype.String() == "CONSTR_FOREIGN" {
		LogAssert(config.Options.LogLevel, len(constraintNode.PkAttrs) != 1, "Unhandled amount of foreign key references %d. %+v", len(constraintNode.FkAttrs), constraintNode)
		columnConstraint.ForeignKey = Some(ForeignKeyConstraint{
			ReferencingTableName:  constraintNode.Pktable.Relname,
			ReferencingColumnName: constraintNode.PkAttrs[0].GetString_().Str, // FIXME can reference multiple columns
			MatchType:             constraintNode.FkMatchtype,
			OnDeleteAction:        constraintNode.FkDelAction,
			OnUpdateAction:        constraintNode.FkUpdAction,
		})
	} else if constraintNode.Contype.String() == "CONSTR_PRIMARY" {
		columnConstraint.PrimaryKey = Some(PrimaryKeyConstraint{})
	} else if constraintNode.Contype.String() == "CONSTR_UNIQUE" {
		columnConstraint.Unique = Some(UniqueConstraint{})
	} else if constraintNode.Contype.String() == "CONSTR_NOTNULL" {
		columnConstraint.NotNull = true
	} else if constraintNode.Contype.String() == "CONSTR_NULL" {
		columnConstraint.NotNull = false
	} else if constraintNode.Contype.String() == "CONSTR_IDENTITY" {
		columnConstraint.NotNull = true
		columnConstraint.Identity = Some(IdentityConstraint{
			GeneratedAlways: constraintNode.GeneratedWhen == "a",
		})
	} else if constraintNode.Contype.String() == "CONSTR_GENERATED" {
		columnConstraint.GeneratedValue = Some(GeneratedValueConstraint{
			GeneratedAlways: constraintNode.GeneratedWhen == "a",
			Expression:      constraintNode.RawExpr.String(),
		})
	} else if constraintNode.Contype.String() == "CONSTR_DEFAULT" {
		columnConstraint.DefaultValue = Some(DefaultConstraint{
			Expression: constraintNode.RawExpr.String(),
		})
	} else {
		LogWarning(config.Options.LogLevel, "Unimplemented constraint node type '%s'", constraintNode.Contype.String())
		LogWarning(config.Options.LogLevel, "%+v", constraintNode)
	}
}

func BuildTableConstraint(config *PacConfig, tableConstraint *TableConstraint, constraintNode *pg_query.Constraint) {
	if constraintNode.Contype.String() == "CONSTR_FOREIGN" {
		LogAssert(config.Options.LogLevel, len(constraintNode.FkAttrs) != 1, "Unhandled amount of foreign key sources %d. %+v", len(constraintNode.FkAttrs), constraintNode)
		LogAssert(config.Options.LogLevel, len(constraintNode.PkAttrs) != 1, "Unhandled amount of foreign key references %d. %+v", len(constraintNode.FkAttrs), constraintNode)
		tableConstraint.ForeignKey = Some(ForeignKeyConstraint{
			ColumnName:            constraintNode.FkAttrs[0].GetString_().Str, // FIXME can reference multiple columns
			ReferencingTableName:  constraintNode.Pktable.Relname,
			ReferencingColumnName: constraintNode.PkAttrs[0].GetString_().Str, // FIXME can reference multiple columns
			MatchType:             constraintNode.FkMatchtype,
			OnDeleteAction:        constraintNode.FkDelAction,
			OnUpdateAction:        constraintNode.FkUpdAction,
		})
	} else if constraintNode.Contype.String() == "CONSTR_PRIMARY" {
		columnNames := []string{}
		for _, s := range constraintNode.Keys {
			columnNames = append(columnNames, s.GetString_().Str)
		}
		tableConstraint.PrimaryKey = Some(PrimaryKeyConstraint{
			ColumnNames: columnNames,
		})
	} else if constraintNode.Contype.String() == "CONSTR_UNIQUE" {
		keys := []string{}
		for _, key := range constraintNode.Keys {
			keys = append(keys, key.GetString_().Str)
		}
		tableConstraint.Unique = Some(UniqueConstraint{
			ColumnNames: keys,
		})
	} else {
		LogWarning(config.Options.LogLevel, "Unimplemented constraint node type '%s'", constraintNode.Contype.String())
		LogWarning(config.Options.LogLevel, "%+v", constraintNode)
	}
}
