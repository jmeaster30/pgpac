package main

import (
	"github.com/jackc/pgx"
)

type Plan struct {
	idToObjectMap      map[int64]PlanObject
	parallelObjectList [][]int64
}

type Operation int

const (
	UNSET Operation = iota
	REPLACE
	CREATE
	DELETE
	SKIP
)

type PlanObject struct {
	PlanObjectId int64
	Operation    Operation
	SchemaObject SchemaObject
}

func BuildPlan(config *PacConfig, objs []SchemaObject, host string, port uint16, database string, user string, pass string) *Plan {
	LogInfo(config.Options.LogLevel, "Starting to build plan")
	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     host,
		Port:     port,
		Database: database,
		User:     user,
		Password: pass,
	})
	if err != nil {
		LogError(config.Options.LogLevel, "Unable to connect to database :(")
		panic(err)
	}
	defer conn.Close()

	rows, err := conn.Query("SELECT table_name FROM information_schema.tables;")
	if err != nil {
		LogError(config.Options.LogLevel, "Unable to query database :(")
		panic(err)
	}

	vals, err := rows.Values()
	if err != nil {
		LogError(config.Options.LogLevel, "Some error happened :(")
		panic(err)
	}

	LogInfo(config.Options.LogLevel, "%+v", vals)

	return &Plan{}
}
