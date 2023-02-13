package main

type Table struct {
	TableName string
}

type Column struct {
	name    string
	coltype ColumnType
}

type ColumnType struct {
	baseType string
}
