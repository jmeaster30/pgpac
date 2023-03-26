package main

import (
	"fmt"
	"strings"
)

type SchemaObject interface {
	getDependencies(objects []SchemaObject) ([]SchemaObject, error)
	isNameMatch(objectType string, name string) bool
}

func findMatch(objectList []SchemaObject, objType string, name string) []SchemaObject {
	results := []SchemaObject{}
	for _, o := range objectList {
		if o.isNameMatch(objType, name) {
			results = append(results, o)
		}
	}
	return results
}

// create table
type Table struct {
	TableName       string
	TableConstraint TableConstraint
	Columns         []Column
}

func (t Table) getDependencies(objects []SchemaObject) ([]SchemaObject, error) {
	return []SchemaObject{}, nil
}

func (t Table) isNameMatch(objectType string, name string) bool {
	if objectType == "table" {
		return t.TableName == name
	}
	return false
}

func (t *Table) String() string {
	value := fmt.Sprintf("Table: \033[34m\033[4m%s\033[0m", t.TableName)
	value += fmt.Sprintf("\n\tConstraints: %s", t.TableConstraint.String())
	for _, col := range t.Columns {
		value += "\n\t" + col.String()
	}
	return value
}

type Column struct {
	Name       string
	TableName  string
	ColumnType ColumnType
	Constraint ColumnConstraint
}

func (c Column) getDependencies(objects []SchemaObject) ([]SchemaObject, error) {
	tableMatches := findMatch(objects, "table", c.TableName)
	if len(tableMatches) == 0 {
		return []SchemaObject{}, fmt.Errorf("expected to find a match for table name '%s'", c.TableName)
	}
	if len(tableMatches) > 1 {
		return []SchemaObject{}, fmt.Errorf("too many matches for table name '%s'", c.TableName)
	}

	results := tableMatches

	// get last name? in typepath
	// not sure if we should check all of them
	baseType := c.ColumnType.TypePath[len(c.ColumnType.TypePath)-1]
	if !isPgType(baseType) {
		typeMatches := findMatch(objects, "type", baseType)
		if len(tableMatches) == 0 {
			return []SchemaObject{}, fmt.Errorf("expected to find a match for type '%s'", baseType)
		}
		if len(tableMatches) > 1 {
			return []SchemaObject{}, fmt.Errorf("too many matches for type '%s'", baseType)
		}
		results = append(results, typeMatches[0])
	}

	return results, nil
}

func (c Column) isNameMatch(objectType string, name string) bool {
	if objectType == "column" {
		return c.Name == name
	}
	return false
}

func (c *Column) String() string {
	return fmt.Sprintf("Column: \033[34m\033[4m%s\033[0m \033[1;37m[%s]\033[0m%s", c.Name, c.ColumnType.String(), c.Constraint.String())
}

type ColumnType struct {
	TypePath []string
	TypeMods []int // FIXME I believe these don't have to be just integers
}

func (c *ColumnType) String() string {
	s := ""
	for i, val := range c.TypeMods {
		if i != 0 {
			s += ", "
		}
		s += fmt.Sprint(val)
	}
	return fmt.Sprintf("%s(%s)", strings.Join(c.TypePath, "."), s)
}

func isPgType(typename string) bool {
	// pulled from https://www.postgresql.org/docs/current/datatype.html
	return typename == "bigint" ||
		typename == "int8" ||
		typename == "bigserial" ||
		typename == "serial8" ||
		typename == "bit" ||
		typename == "varbit" ||
		typename == "boolean" ||
		typename == "bool" ||
		typename == "box" ||
		typename == "bytea" ||
		typename == "character" ||
		typename == "char" ||
		typename == "varchar" ||
		typename == "cidr" ||
		typename == "circle" ||
		typename == "date" ||
		typename == "float8" ||
		typename == "inet" ||
		typename == "integer" ||
		typename == "int" ||
		typename == "int4" ||
		typename == "interval" ||
		typename == "json" ||
		typename == "jsonb" ||
		typename == "line" ||
		typename == "lseg" ||
		typename == "macaddr" ||
		typename == "macaddr8" ||
		typename == "money" ||
		typename == "numeric" ||
		typename == "decimal" ||
		typename == "path" ||
		typename == "pg_lsn" ||
		typename == "pg_snapshot" ||
		typename == "point" ||
		typename == "polygon" ||
		typename == "real" ||
		typename == "float4" ||
		typename == "smallint" ||
		typename == "int2" ||
		typename == "smallserial" ||
		typename == "serial2" ||
		typename == "serial" ||
		typename == "serial4" ||
		typename == "text" ||
		typename == "time" ||
		typename == "timetz" ||
		typename == "timestamp" ||
		typename == "timestamptz" ||
		typename == "tsquery" ||
		typename == "tsvector" ||
		typename == "txid_snapshot" ||
		typename == "uuid" ||
		typename == "xml"
}

type TableConstraint struct {
	Name       Optional[string]
	Check      Optional[CheckConstraint]
	Unique     Optional[UniqueConstraint]
	PrimaryKey Optional[PrimaryKeyConstraint]
	Exclusion  Optional[ExclusionConstraint]
	ForeignKey Optional[ForeignKeyConstraint]
}

func (t *TableConstraint) String() string {
	value := ""
	if t.ForeignKey.HasValue() {
		value += fmt.Sprintf("\n\t\t\033[1;36m%s\033[0m", t.ForeignKey.Value().String())
	}
	if t.PrimaryKey.HasValue() {
		value += fmt.Sprintf("\n\t\t\033[1;32m%s\033[0m", t.PrimaryKey.Value().String())
	}
	return value
}

type ColumnConstraint struct {
	Name           Optional[string]
	NotNull        bool
	Check          Optional[CheckConstraint]
	DefaultValue   Optional[DefaultConstraint]
	GeneratedValue Optional[GeneratedValueConstraint]
	Identity       Optional[IdentityConstraint]
	ForeignKey     Optional[ForeignKeyConstraint]
	Unique         Optional[UniqueConstraint]
	PrimaryKey     Optional[PrimaryKeyConstraint]
	// collation
}

func (c *ColumnConstraint) String() string {
	value := ""
	if c.NotNull {
		value += "\n\t\t\033[1;31mNotNull\033[0m"
	}
	if c.ForeignKey.HasValue() {
		value += fmt.Sprintf("\n\t\t\033[1;36m%s\033[0m", c.ForeignKey.Value().String())
	}
	if c.Identity.HasValue() {
		value += fmt.Sprintf("\n\t\t\033[1;32m%s\033[0m", c.Identity.Value().String())
	}
	if c.DefaultValue.HasValue() {
		value += fmt.Sprintf("\n\t\t\033[1;32m%s\033[0m", c.DefaultValue.Value().String())
	}
	if c.GeneratedValue.HasValue() {
		value += fmt.Sprintf("\n\t\t\033[1;32m%s\033[0m", c.GeneratedValue.Value().String())
	}
	if c.PrimaryKey.HasValue() {
		value += fmt.Sprintf("\n\t\t\033[1;32m%s\033[0m", c.PrimaryKey.Value().String())
	}
	return value
}

type CheckConstraint struct {
	Expression string
	NoInherit  bool
}

type DefaultConstraint struct {
	Expression string
}

func (d DefaultConstraint) String() string {
	return fmt.Sprintf("Default '%s'", d.Expression)
}

type GeneratedValueConstraint struct {
	GeneratedAlways bool
	Expression      string
}

func (g GeneratedValueConstraint) String() string {
	defaultString := "By Default"
	if g.GeneratedAlways {
		defaultString = "Always"
	}
	return fmt.Sprintf("Generated %s as '%s'", defaultString, g.Expression)
}

type IdentityConstraint struct {
	GeneratedAlways bool // false is By Default
	// TODO sequence options
}

func (i IdentityConstraint) String() string {
	defaultString := "By Default"
	if i.GeneratedAlways {
		defaultString = "Always"
	}
	return fmt.Sprintf("Identity Generated %s", defaultString)
}

type UniqueConstraint struct {
	ColumnNames []string // only needed for table constraints
	Nulls       bool
	NotDistinct bool
	// TODO index parameters
}

func (u UniqueConstraint) String() string {
	result := "Unique ("
	for i, col := range u.ColumnNames {
		result += col
		if i != len(u.ColumnNames)-1 {
			result += ", "
		}
	}
	result += ")"
	return result
}

type PrimaryKeyConstraint struct {
	ColumnNames     []string // only needed for table constraints
	IndexParameters string   // TODO
}

func (p PrimaryKeyConstraint) String() string {
	result := "Primary Key ("
	for i, col := range p.ColumnNames {
		result += col
		if i != len(p.ColumnNames)-1 {
			result += ", "
		}
	}
	result += ")"
	return result
}

type ForeignKeyConstraint struct {
	ColumnName            string // FIXME can reference multiple columns
	ReferencingTableName  string
	ReferencingColumnName string // FIXME can reference multiple columns
	MatchType             string // "f", "p", or "s"
	OnDeleteAction        string // TODO referential action
	OnUpdateAction        string // TODO referential action
}

func (f ForeignKeyConstraint) String() string {

	return fmt.Sprintf("References %s(%s) MatchType: '%s'", f.ReferencingTableName, f.ReferencingColumnName, f.MatchType)
}

type ExclusionConstraint struct {
	IndexType  string
	Exclusions []string // seems to be a list of columnName and operator pairs
}

// create index
type Index struct {
	Name                string
	Unique              bool // only btree indexes can be set as unique
	IndexType           string
	TableName           string
	IndexParams         []IndexParam
	IncludedColumnNames []string
}

func (t Index) getDependencies(objects []SchemaObject) ([]SchemaObject, error) {
	// todo
	return []SchemaObject{}, nil
}

func (i Index) isNameMatch(objectType string, name string) bool {
	if objectType == "index" {
		return i.Name == name
	}
	return false
}

type IndexParam struct {
	ColumnName string
	// TODO ordering
	// TODO nulls ordering
	// TODO indexing expressions
	// TODO partial indexes
	// TODO operator classes
	// TODO collation
}

// create extension
type Extension struct {
	ExtensionName string
}

func (t Extension) getDependencies(objects []SchemaObject) ([]SchemaObject, error) {
	return []SchemaObject{}, nil
}

func (ext Extension) isNameMatch(objectType string, name string) bool {
	if objectType == "extension" || objectType == "type" {
		return ext.ExtensionName == name
	}
	return false
}

// create enum type
type Enum struct {
	Name   string
	Values []string
}

func (t Enum) getDependencies(objects []SchemaObject) ([]SchemaObject, error) {
	return []SchemaObject{}, nil
}

func (enum Enum) isNameMatch(objectType string, name string) bool {
	if objectType == "enum" || objectType == "type" {
		return enum.Name == name
	}
	return false
}
