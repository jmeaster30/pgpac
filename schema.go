package main

import (
	"fmt"
	"strings"
)

type IndexType int

const (
	IndexType_BTREE IndexType = iota
	IndexType_HASH
	IndexType_GIST
	IndexType_SPGIST
	IndexType_GIN
	IndexType_BRIN
)

// create table
type Table struct {
	TableName string
	//tableConstraints []TableConstraint
	Columns []Column
}

func (t *Table) String() string {
	value := fmt.Sprintf("Table: \033[34m\033[4m%s\033[0m", t.TableName)
	for _, col := range t.Columns {
		value += "\n\t" + col.String()
	}
	return value
}

type Column struct {
	Name       string
	ColumnType ColumnType
	Constraint ColumnConstraint
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

type TableConstraint struct {
	Name       Optional[string]
	Check      Optional[CheckConstraint]
	Unique     Optional[UniqueConstraint]
	PrimaryKey Optional[PrimaryKeyConstraint]
	Exclusion  Optional[ExclusionConstraint]
	ForeignKey Optional[ForeignKeyConstraint]
}

type ColumnConstraint struct {
	Name              Optional[string]
	NotNull           bool
	Check             Optional[CheckConstraint]
	DefaultValue      Optional[DefaultConstraint]
	GeneratedValue    Optional[GeneratedValueConstraint]
	GeneratedIdentity Optional[GeneratedIdentityConstraint]
	ForeignKey        Optional[ForeignKeyConstraint]
	// collation
}

func (c *ColumnConstraint) String() string {
	value := ""
	if c.NotNull {
		value += "\t\t\033[1;31mNotNull\033[0m"
	}

	if c.ForeignKey.HasValue() {
		value += fmt.Sprintf("\t\t\033[1;36m%s\033[0m", c.ForeignKey.Value().String())
	}

	if value == "" {
		return ""
	}
	return "\n" + value
}

type CheckConstraint struct {
	Expression string
	NoInherit  bool
}

type DefaultConstraint struct {
	Expression string
}

type GeneratedValueConstraint struct {
	Expression string
}

type GeneratedIdentityConstraint struct {
	Always bool // false is By Default
	// TODO sequence options
}

type UniqueConstraint struct {
	ColumnNames []string // only needed for table constraints
	Nulls       bool
	NotDistinct bool
	// TODO index parameters
}

type PrimaryKeyConstraint struct {
	ColumnNames     []string // only needed for table constraints
	IndexParameters string   // TODO
}

type ForeignKeyConstraint struct {
	ColumnNames           []string // only needed for table constraints
	ReferencingTableName  string
	ReferencingColumnName string // FIXME can reference multiple columns
	MatchFull             bool
	MatchPartial          bool
	MatchSimple           bool
	OnDeleteAction        string // TODO referential action
	OnUpdateAction        string // TODO referential action
}

func (f ForeignKeyConstraint) String() string {
	return fmt.Sprintf("REFERENCES %s(%s)", f.ReferencingTableName, f.ReferencingColumnName)
}

type ExclusionConstraint struct {
	IndexType  IndexType
	Exclusions []string // seems to be a list of columnName and operator pairs
}

// create index
type Index struct {
	Name                string
	Unique              bool // only btree indexes can be set as unique
	IndexType           IndexType
	TableName           string
	ColumnNames         []string
	IncludedColumnNames []string
	// TODO ordering
	// TODO indexing expressions
	// TODO partial indexes
	// TODO operator classes
	// TODO collation
}

type Type interface {
}

// create extension
type Extension struct {
	ExtensionName string
}

// create enum type
type Enum struct {
	Name   string
	Values []string
}
